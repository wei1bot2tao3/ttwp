package handler

import (
	"bysj0.3/cmd/app1/db"
	"bysj0.3/cmd/app1/meta"
	ossapi "bysj0.3/cmd/app1/oss"
	"fmt"
	"github.com/gin-gonic/gin"

	"io"
	"net/http"
	"os"
	"time"
)

// CheckLogin 登录密码判定
func CheckLogin(c *gin.Context) {
	//接收参数
	username := c.PostForm("username")
	usaerpassword := c.PostForm("password")

	//数据库判定
	user, check := db.GetUserMeta(username)
	//返回值
	if check == false {
		//用户名不存在
		fmt.Printf("用户名不存在")
		c.JSON(http.StatusOK, gin.H{"coed": 0, "web": "用户名不存在"})

	}
	if check == true && user.Password != usaerpassword {
		fmt.Printf("密码错误")
		c.JSON(http.StatusOK, gin.H{"coed": 0, "web": "密码错误"})
	}
	if user.Status == "0" {
		fmt.Printf("对不起您账号被封禁")
		c.JSON(http.StatusOK, gin.H{"coed": 0, "web": "对不起您账号被封禁"})

	}
	if check == true && user.Password == usaerpassword && user.Status == "1" {
		fmt.Printf("登录成功")
		//token
		//生成token
		//获取当前时间生成过期时间
		extime := time.Now().Add(1 * time.Hour).Unix()
		token, _ := GenerateToken(username, extime)
		//把token写如数据库
		//fmt.Println("请求写入数据库")
		//ok := db.CreateToken(username, token)
		//if ok {
		//	fmt.Printf("写入数据库成功")
		//}
		//写入redis
		err := StoreToken(username, token)
		if err != nil {
			fmt.Printf("写入失败")
		}
		//打开file页面 返回token给全端
		//返回根目录
		c.JSON(http.StatusOK, gin.H{"result": "success", "username": username, "token": token, "userid": user.Userid})

	}

}

// 注册
// 用户文件夹 打开注册
func NewUser(c *gin.Context) {
	html, err := os.ReadFile("/root/bysj0.3/web/login/register.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}

	io.WriteString(c.Writer, string(html))
}

// 注册用户
func CreateUser(c *gin.Context) {
	//接收参数
	newuser := meta.User{
		Username:  c.PostForm("username"),
		Password:  c.PostForm("password"),
		Email:     c.PostForm("email"),
		Status:    "1",
		Lsb:       "20",
		Privilege: c.PostForm("privilege"),
		Phone:     c.PostForm("phone"),
	}
	//先写入 在返回userid 然后在组成表明
	umeta, _ := db.GetUserMeta(newuser.Username)

	////但是用户名修改的话需要判断是否重名
	if newuser.Username != umeta.Username {
		//数据库写入 users表 并且返回
		ok := db.WriteUser(newuser)
		if ok == 0 {
			c.JSON(http.StatusOK, gin.H{"ok": 0})
		} else {
			//本地创建用户文件夹
			created, _ := db.GetUserMeta(newuser.Username)
			pwd := fmt.Sprintf("/root/file/%d/home", created.Userid)
			err := os.MkdirAll(pwd, 0750)
			if err != nil && !os.IsExist(err) {
				fmt.Printf("本地文件夹创建失败：%v", err)
			}
			fmt.Println("开始创建用户表")
			createdTable := fmt.Sprintf("%dtable", created.Userid)
			fmt.Println(createdTable + "准备写入")

			//写入user表
			err = db.UpdateUserMeta(created.Userid, "table", createdTable)
			if err != nil {
				fmt.Println(err)
			}
			//创建用户专属文件表
			db.CreateUserTable(createdTable)
			//oss 创建oss用户目录
			ossuser1 := fmt.Sprintf("file/%d/home", created.Userid)
			ossuser := fmt.Sprintf("file/%d", created.Userid)
			fmt.Println(ossuser)
			fmt.Println(ossuser1)
			ossapi.Mkdir(ossuser)
			if !ossapi.Mkdir(ossuser1) {
				fmt.Println("创建oss用户目录失败")
			}
			//把根目录写进去，然后设置为0
			m := meta.CretaeFileMeta{
				Userid:       created.Userid,
				Filename:     "root",
				Filelocation: pwd,
				Filestatus:   1,
				Filetype:     "wjj",
				Parent:       0,
				Osslocation:  ossuser1,
			}
			db.WriteCretaeFile(createdTable, m)

		}
		//返回给页面
		c.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		c.JSON(http.StatusOK, gin.H{"ok": 0})
	}
}

// 打开个人主页
func GetUserInfo(c *gin.Context) {

	html, err := os.ReadFile("/root/bysj0.3/web/user/uesrindex.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}
	io.WriteString(c.Writer, string(html))
}

// 返回个人主页的信息
func GetUserMeta(c *gin.Context) {
	//获取参数 用户名 就ok
	Username := c.PostForm("username")
	//通过和数据库获取信息
	usermeat, _ := db.GetUserMeta(Username)
	//返回用户名，手机号邮件
	c.JSON(http.StatusOK, gin.H{"username": usermeat.Username, "phone": usermeat.Phone, "email": usermeat.Email})
}

// 修改函数用户信息
func PsotUserUMeta(c *gin.Context) {
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	username := c.PostForm("username")
	newname := c.PostForm("newusername")
	////修改
	umeta, _ := db.GetUserMeta(username)

	////但是用户名修改的话需要判断是否重名
	if newname != username {
		//检查用户是否重名
		ok := db.CheckUser(newname)
		if ok {
			db.UpdateUserMeta(umeta.Userid, "username", newname, "email", email, "phone", phone)
			c.JSON(http.StatusOK, gin.H{"ok": 1, "username": newname})
		} else {
			c.JSON(http.StatusOK, gin.H{"ok": 0, "username": username})
		}
	} else {
		db.UpdateUserMeta(umeta.Userid, "email", email, "phone", phone)
		c.JSON(http.StatusOK, gin.H{"ok": 1, "username": newname})
	}

}

// 打开修改密码页面
func GetUserPassword(c *gin.Context) {
	html, err := os.ReadFile("/root/bysj0.3/web/user/password.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}
	io.WriteString(c.Writer, string(html))

}

// 修改密码
func ChangePassword(c *gin.Context) {
	//收取参数
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	//confirmPassword := c.PostForm("confirm_password")
	username := c.PostForm("username")

	// 处理密码修改逻辑
	user, _ := db.GetUserMeta(username)
	if user.Password == oldPassword {
		fmt.Printf("原始密码ok开始进行修改密码")
		db.NewPassword(username, newPassword)
		c.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		c.JSON(http.StatusFound, gin.H{"ok:": 0})
	}

}

// 查询操作记录

func GetHistory(c *gin.Context) {
	//打开
	html, err := os.ReadFile("/root/bysj0.3/web/user/userhistory.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}

	io.WriteString(c.Writer, string(html))
}
func UserHistory(c *gin.Context) {
	//获取参数用户id就OK
	uderid := c.Query("userid")
	//查询id
	data := db.UserHsitory(uderid)
	//返回
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})

}

func Exit(c *gin.Context) {
	//获取用户，用户名
	username := c.PostForm("username")
	// 连接Redis数据库
	fmt.Println(username)
	DeleteToken(username)
	c.JSON(http.StatusOK, gin.H{"code": 1})
}
