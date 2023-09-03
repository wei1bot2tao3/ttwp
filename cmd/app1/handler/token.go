package handler

import (
	cache "bysj0.3/cmd/app1/cache/redis"
	"bysj0.3/cmd/app1/db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

// 生成结构体

type UserToken struct {
	Username           string `json:"username"`
	jwt.StandardClaims        //标准格式 主要是过期时间

}

// 定义密钥
var SecretKey = []byte("my_secret_key")

// 生成Token
func GenerateToken(username string, exit int64) (string, error) {
	// 获取当前时间然后设置过期时间
	expirationTime := exit
	//赋值给结构体
	usertoken := UserToken{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}
	// 创建一个JWT实例
	//根据签名生成token，NewWithClaims(加密方式,claims) ==》 头部，载荷，签证
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, usertoken)
	// 签名token并返回字符串形式的结果
	return token.SignedString(SecretKey)
}

// 验证token
func CheckToken(token string) (*UserToken, error) {
	//  ParseToken 解析token
	// 调用ParseWithClaims进行解析，使用签名解密出数据
	usertoken, err := jwt.ParseWithClaims(token, &UserToken{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if usertoken != nil {
		if c, ok := usertoken.Claims.(*UserToken); ok && usertoken.Valid {
			fmt.Println(c)
			return c, nil
		}
	}
	return nil, err
}

// 流程
func CheckPermission(c *gin.Context) {
	fmt.Printf("开始验证token")
	//获取token
	token := c.Query("token")
	if token == "" {
		token = c.PostForm("token")
	}
	fmt.Print("接收到到tokne：%v", token)
	//加一个判单token是否存在的

	//解析token
	usertoken, err := CheckToken(token)
	if err != nil {
		fmt.Println(err)
		fmt.Println("")
		c.Abort()
		return

	}
	//数据库判定
	_, check := db.GetUserMeta(usertoken.Username)
	if check == false {
		fmt.Println("用户不存在")
		c.Redirect(http.StatusMovedPermanently, "http://www.wbt1104.xyz/")
		c.Abort()
		return
	}

	main := time.Now().Add(30 * time.Minute).Unix()
	//超时怎么办 中止后面的函数
	if usertoken.ExpiresAt < time.Now().Unix() {
		fmt.Println("对不起，请您重新登录")
		c.JSON(http.StatusOK, gin.H{"text": "对不起您超时重新登录"})
		c.Redirect(http.StatusMovedPermanently, "http://www.wbt1104.xyz/")
		c.Abort()
		return
	} else if usertoken.ExpiresAt <= main {
		fmt.Println("token加时间")
		//要给token小于30分钟就加一小时
		usertoken.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		//再次生成一个新的token 传到下面到函数里然后返回
		NewTokenAdd, _ := GenerateToken(usertoken.Username, usertoken.ExpiresAt)
		c.Set("NewTokenAdd", NewTokenAdd)
		fmt.Println(NewTokenAdd)
		c.Next()
	} else {
		c.Next()
	}

}

type CheckRequest struct {
	Username string `json:"username"`
}

// 检查第一token
func FirstCheckToken(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token := strings.Replace(tokenString, "Bearer ", "", 1)
	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := req.Username
	fmt.Println(username)
	ctoken, err := CheckToken(token)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusFound, gin.H{"data": "token过期"})
		return
	}
	//读取redis里的token
	Rtoken, err := GetRedisToken(username)
	if Rtoken != token {
		//返回状态码
		fmt.Println("检查失败")
		c.JSON(http.StatusFound, gin.H{"data": "token匹配失败"})
		return
	}
	//检查时间
	main := time.Now().Add(30 * time.Minute).Unix()
	if ctoken.ExpiresAt <= main {
		//制作新的token
		exit := time.Now().Add(30 * time.Minute).Unix()
		token, _ = GenerateToken(token, exit)
		//更新redis
		UpdateToken(username, token)
		c.JSON(http.StatusOK, nil)
		return
	} else {
		fmt.Println(Rtoken)
		fmt.Println("redistoken检查成功")
		c.JSON(http.StatusOK, nil)
	}

	//超时重新生成一个redis存入

}

// 写入redis
// 存储用户名和token到Redis
func StoreToken(username string, token string) error {
	// 连接Redis数据库
	client, err := cache.Redis()
	if err != nil {
		return err
	}
	defer client.Close()
	// 设置键值对
	err = client.HSet(context.Background(), "tokens", username, token).Err()
	if err != nil {
		return err
	}
	// 设置过期时间为1小时
	err = client.Expire(context.Background(), username, time.Hour).Err()
	if err != nil {
		return err
	}
	fmt.Println("写入成功")
	return nil
}

// 获取redistoken
func GetRedisToken(username string) (string, error) {
	// 连接Redis数据库
	client, err := cache.Redis()
	if err != nil {
		return "", err
	}
	defer client.Close()
	// 获取键值对
	fmt.Println("开始读取redistoken")
	fmt.Println(username)
	token, err := client.HGet(context.Background(), "tokens", username).Result()
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	fmt.Println(token)
	return token, nil
}

// 更新redistoken
func UpdateToken(username string, newToken string) error {
	// 连接Redis数据库
	client, err := cache.Redis()
	if err != nil {
		return err
	}
	defer client.Close()
	// 更新键值对
	err = client.HSet(context.Background(), "tokens", username, newToken).Err()
	if err != nil {
		return err
	}
	//延长时间
	// 设置过期时间为1小时
	err = client.Expire(context.Background(), username, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

// 删除redis
func DeleteToken(username string) (int, error) {
	// 连接Redis数据库
	client, err := cache.Redis()
	if err != nil {
		return 0, err
	}
	defer client.Close()
	// 删除键值对
	fmt.Println("开始删除redistoken")
	fmt.Println(username)
	deletedCount, err := client.HDel(context.Background(), "tokens", username).Result()
	fmt.Println(err)
	if err != nil {
		return 0, err
	}
	fmt.Println(deletedCount)
	return 1, nil
}

func CheckLoginJs(c *gin.Context) {
	token := c.PostForm("token")
	usernmae := c.PostForm("username")
	fmt.Println("获取token" + token)
	fmt.Println(token)
	rtoken, _ := GetRedisToken(usernmae)
	if rtoken != token {
		c.JSON(http.StatusOK, gin.H{"code": -1})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1})
	}

}
