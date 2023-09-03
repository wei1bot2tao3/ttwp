package handler

import (
	"bysj0.3/cmd/app1/db"
	"bysj0.3/cmd/app1/meta"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
)

// 通过用户名返回表名
func UserTable(username string) string {
	meta, _ := db.GetUserMeta(username)
	return meta.Table
}

// 文件模块
// GeT 接收用户姓名 打开用户个人文件首页
func Home(c *gin.Context) {
	html, err := os.ReadFile("/root/bysj0.3/web/index.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}
	io.WriteString(c.Writer, string(html))

}

// 获取根目录
func GetHome(c *gin.Context) {
	username := c.Query("username")
	tablename := UserTable(fmt.Sprint(username))
	token, _ := c.Get("NewTokenAdd")
	if token == nil {
		token = c.Query("token")
	}
	form, err := db.Showtabel(tablename, 1)
	if err != nil {
		fmt.Printf("数据库查询失败：%v", err)
	}
	forms := meta.FormJson{
		Code: http.StatusOK,
		Data: form,
	}
	c.JSON(200, gin.H{"code": forms.Code, "data": forms.Data, "token": token, "tree": "home", "parentid": 1})
	//c.Writer.Write(jsons)
}

// 获取选择的目录文件
func SelectUser(c *gin.Context) {
	//通过查询当前用户表格名
	//处理请求
	username := c.Query("username")
	treeed := c.Query("tree")
	parent := c.Query("tolocation")
	parentid := ParentId(parent)
	parentname := c.Query("parentname")
	//收到了 需要访问的id tree 拼接一下
	//返回目录
	tree := treeed + "/" + parentname
	//应该是用redis返回的
	tablename := UserTable(fmt.Sprint(username))
	token, _ := c.Get("NewTokenAdd")
	if token == nil {
		token = c.Query("token")
	}
	//查询当前路径的所有文件

	form, err := db.Showtabel(tablename, parentid)
	if err != nil {
		fmt.Printf("数据库查询失败：%v", err)
	}
	//返回成jison
	forms := meta.FormJson{
		Code: http.StatusOK,
		Data: form,
	}
	//返回当前目录和过去目录
	fmt.Println("开始返回")
	c.JSON(200, gin.H{"code": forms.Code, "data": forms.Data, "token": token, "tree": tree, "parentid": parentid})
}

// 搜索
func SearchFilemeta(c *gin.Context) {
	//获取前端信息 返回出用户表名
	username := c.Query("username")
	key := c.Query("key")
	tablename := UserTable(fmt.Sprint(username))
	//通过数据库获取想要的文件
	data := db.SearchFile(tablename, key)
	//返回道前端
	c.JSON(200, gin.H{"code": http.StatusOK, "data": data})

}

// 分类
func SearchFileType(c *gin.Context) {
	//获取信息返回表名
	username := c.Query("username")
	searchtype := c.Query("type")
	tablename := UserTable(fmt.Sprint(username))
	//查询数据库
	data := db.Searchtype(tablename, searchtype)
	//返回道前端
	c.JSON(200, gin.H{"code": http.StatusOK, "data": data})

}

// 打开上传页面
func Upfile(c *gin.Context) {
	html, err := os.ReadFile("/root/bysj0.3/web/upfile.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}
	io.WriteString(c.Writer, string(html))

}

func ParentId(s string) int {

	var parentid int
	var err error
	if s == " " {
		parentid = 1
	} else {
		parentid, err = strconv.Atoi(s)
		if err != nil {
			fmt.Printf("转换失败")
		}
	}
	return parentid
}

// 七天回收
func ChecckDay7(c *gin.Context) {
	username := c.Query("username")
	tablename := UserTable(fmt.Sprint(username))

	data, _ := db.ShowDay7(tablename)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": data})
}

// 用户模块
// 登录
// 加载登录页面
func Login(c *gin.Context) {
	html, err := os.ReadFile("/root/bysj0.3/web/login/test.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}

	io.WriteString(c.Writer, string(html))
}
