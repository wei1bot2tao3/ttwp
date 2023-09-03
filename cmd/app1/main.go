package main

import (
	"bysj0.3/cmd/app1/handler"
	"github.com/gin-gonic/gin"
)

func main() {

	//创建一个默认的路由引擎
	route := gin.Default()

	//注册全局中间件
	// 启动静态文件服务
	route.Static("/layui", "/root/bysj0.3/web/layui")
	//TODO
	//暂时不用
	//route.Static("/img", "/root/img/")
	//文件模块
	//定义一个路由组 添加检查jwt函数
	FileApi := route.Group("/file", handler.CheckPermission)
	{
		//获取当前用户文件信息
		FileApi.GET("/home", handler.GetHome)
		//搜索文件
		FileApi.GET("/search", handler.SearchFilemeta)
		FileApi.GET("/type", handler.SearchFileType)
		//上传文件页面
		//上传接口
		FileApi.GET("/upfile", handler.Upfile)
		//先判断是不是有存在的文件
		FileApi.POST("/upfile", handler.CheckOnlyFile, handler.UpfilePost)
		//TODO
		//大文件
		//小文件
		//下载接口
		FileApi.POST("/down", handler.DownFile)
		//删除接口
		FileApi.POST("/delete", handler.DeleteFile)
		//修改接口
		FileApi.POST("/edit", handler.InstallFileMeta)
		//创建文件接口
		FileApi.POST("/tree/mkdir", handler.Mkdirone)
		FileApi.GET("/tree/tolocation", handler.SelectUser)
		//回收站
		FileApi.GET("/day07", handler.ChecckDay7)
		FileApi.POST("/day07/0", handler.DeleteFileTO)
		FileApi.POST("/day07/1", handler.DeleteFileT1)
		//histroy

	}
	//用户模块
	UserApi := route.Group("/user")
	{
		//登录页面打开
		UserApi.GET("/login", handler.Login)
		//判断登录
		UserApi.POST("/login", handler.CheckLogin)
		//获取注册页面
		UserApi.GET("/register", handler.NewUser)
		//接收注册信息
		UserApi.POST("/register", handler.CreateUser)
		//打开用户页面
		UserApi.GET("/index", handler.Home)
		//第一次检查token
		UserApi.POST("/check", handler.FirstCheckToken)
		//个人主页
		UserApi.GET("/info", handler.GetUserInfo)
		//获取个人信息
		UserApi.POST("/meta", handler.GetUserMeta)
		//修改个人信息
		UserApi.POST("/meatupdate", handler.PsotUserUMeta)
		//密码
		UserApi.POST("/password", handler.ChangePassword)
		//获取修改页面
		UserApi.GET("/password", handler.GetUserPassword)
		//history
		UserApi.GET("/history", handler.GetHistory)
		//
		UserApi.GET("historylist", handler.UserHistory)

	}
	AdminApi := route.Group("/admin")
	{
		//登录页面
		//登录页面打开
		AdminApi.GET("/login", handler.AdminGet)
		//判断登录
		AdminApi.POST("/login", handler.CheckLogin)
		AdminApi.GET("/index", handler.AdminIndex)
		AdminApi.GET("/user/list", handler.GetUserList)
		AdminApi.POST("/user/edit", handler.HAdminUsermeta)
		AdminApi.POST("/user/status", handler.HAdminUserStatus)
		AdminApi.POST("/user/password", handler.HAdminPassword)
		AdminApi.GET("/user/search", handler.HAdminSerch)
		AdminApi.POST("/user/createuser", handler.HAdminCreateUser)
		//打开
		AdminApi.GET("/user/info", handler.HAdminUserInfo)
		AdminApi.GET("/user/password", handler.HAdminPsaaword)
		AdminApi.GET("/history", handler.HAdminHistroy)
		AdminApi.GET("/historylist", handler.HAdminAllHsitory)
		//集群
		AdminApi.GET("/k8s", handler.HAdminK8s)
		AdminApi.GET("/k8slist", handler.HAdminK8sTop)
		AdminApi.POST("/user/delete", handler.HAdminDelete)

	}
	//测试登录
	route.POST("/check", handler.CheckLoginJs)
	//退出登录
	route.POST("/exit", handler.Exit)
	//访问跟网址直接访问
	route.GET("/", handler.Login)
	//测试路由
	//服务端口
	route.Run(":8080")

	//直接上测试生成的函数

}
