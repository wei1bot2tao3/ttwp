package handler

import (
	"bysj0.3/cmd/app1/db"
	"bysj0.3/cmd/app1/k8s"
	"bysj0.3/cmd/app1/meta"
	ossapi "bysj0.3/cmd/app1/oss"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
)

func AdminGet(c *gin.Context) {
	html, err := os.ReadFile("/root/bysj0.3/web/login/admin.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}

	io.WriteString(c.Writer, string(html))
}
func AdminIndex(c *gin.Context) {
	html, err := os.ReadFile("/root/bysj0.3/web/user/admin.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}

	io.WriteString(c.Writer, string(html))
}
func GetUserList(c *gin.Context) {
	privilege := c.Query("privilege")
	//直接查询数据流所有非管理员的用户
	if privilege == "xhw" {
		data := db.AllXhw(0)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": data})
	} else {
		data := db.AllUser(privilege)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": data})
	}

	//

}
func HAdminUsermeta(c *gin.Context) {
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	username := c.PostForm("username")
	userid := c.PostForm("userid")
	fmt.Println(userid)
	userId, _ := strconv.Atoi(userid)
	ok := db.CheckUser(username)
	fmt.Println(ok)
	if ok {
		fmt.Println("开始修改")
		fmt.Println(phone, email, username, userId)
		err := db.UpdateUserMeta(userId, "username", username, "email", email, "phone", phone)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{"code": 200})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 0})
	}
}
func HAdminUserStatus(c *gin.Context) {
	userid := c.PostForm("userid")
	userId, _ := strconv.Atoi(userid)
	st := c.PostForm("status")
	status, _ := strconv.Atoi(st)
	//放到数据修改
	db.AdminSatus(userId, status)
	c.JSON(http.StatusOK, gin.H{"code": 200})
}

func HAdminPassword(c *gin.Context) {
	userid := c.PostForm("userid")
	userId, _ := strconv.Atoi(userid)
	db.AdminPassword(userId)
	c.JSON(http.StatusOK, gin.H{"code": 200})

}
func HAdminSerch(c *gin.Context) {
	key := c.Query("value")
	fmt.Println("要查询的名是" + key)

	data := db.AdminSearchUser(key)
	fmt.Println(data)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": data})
}

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// 创建用户
func HAdminCreateUser(c *gin.Context) {
	//获取姓名，
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Username := user.Username
	privilege := user.Role
	Ameta := meta.User{
		Username:  Username,
		Privilege: privilege,
		Status:    "1",
		Lsb:       "20",
		Password:  "123456",
	}
	//如果是admin 直接写入数据库
	//先判断是否重名
	if !db.CheckUser(Username) {
		//右重名
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "对不起重名了"})
	}
	if privilege == "admin" {
		//直接写入数据库

		ok := db.WriteUser(Ameta)
		if ok == 0 {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "对不起重名了"})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建成功"})
	} else {
		//创建用户
		//数据库写入 users表 并且返回
		ok := db.WriteUser(Ameta)
		if ok == 0 {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "对不起重名了"})
		} else {

			//本地创建用户文件夹
			created, _ := db.GetUserMeta(Username)
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
			//把根目录写进去，然后设置为0
			m := meta.CretaeFileMeta{
				Userid:       created.Userid,
				Filename:     "root",
				Filelocation: pwd,
				Filestatus:   1,
				Filetype:     "wjj",
				Parent:       0,
			}
			db.WriteCretaeFile(createdTable, m)
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建成功"})
		}

	}

}

// 打开个人主页
func HAdminUserInfo(c *gin.Context) {

	html, err := os.ReadFile("/root/bysj0.3/web/user/adminindex.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}
	io.WriteString(c.Writer, string(html))
}
func HAdminPsaaword(c *gin.Context) {

	html, err := os.ReadFile("/root/bysj0.3/web/user/adminpassored.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}
	io.WriteString(c.Writer, string(html))
}

func HAdminUserToHsitory(c *gin.Context) {
	//获取id 然后把时间和操作记录写进去

}

// 获取操作记录根据id
func HAdminUserSelect(c *gin.Context) {

}

func HAdminHistroy(c *gin.Context) {

	html, err := os.ReadFile("/root/bysj0.3/web/user/adminhistory.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}
	io.WriteString(c.Writer, string(html))
}

// history
func HAdminAllHsitory(c *gin.Context) {
	data := db.AdminGetAllHistory
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data()})

}

// k8s
func HAdminK8s(c *gin.Context) {

	html, err := os.ReadFile("/root/bysj0.3/web/user/admink8s.html")
	if err != nil {
		fmt.Printf("打开文件失败：err：%v", err)
		return
	}
	io.WriteString(c.Writer, string(html))
}

func HAdminK8sTop(c *gin.Context) {
	data, _ := k8s.GetNodeStatus()
	fmt.Println(data)
	oss := ossapi.GetOss()
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": data, "oss": oss})

}

func HAdminDelete(c *gin.Context) {
	user := c.PostForm("userid")
	userid := ParentId(user)

	if db.AdminDeleteUser(userid) {
		c.JSON(http.StatusOK, gin.H{"code": 200})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "删除失败"})
	}

}
