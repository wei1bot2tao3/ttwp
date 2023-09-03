package handler

import (
	"bysj0.3/cmd/app1/db"
	"bysj0.3/cmd/app1/meta"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 临时测试
func CheckOnlyFile(c *gin.Context) {
	//解析参数
	//获取用户名
	username := c.Query("username")
	tablename := UserTable(fmt.Sprint(username))

	filehash := c.Query("filehash")
	fmt.Println(filehash)
	//唯一表检查是否存在
	onemeta, ok := db.CheckOnlyFiles(filehash)
	//存在返回true 否则false
	if ok {
		//存在执行秒传
		fmt.Printf("存在文件，直接写入当前用:%v\n", onemeta)
		//存在直接写入用户表
		filemeta := meta.CretaeFileMeta{
			Filehash:       filehash,
			Filelocation:   onemeta.Filelocation,
			Userid:         1,
			Filename:       onemeta.Filename,
			Filesize:       onemeta.Filesize,
			Filecreatetime: time.Now().Format("2006-01-02 15:04:06"),
			Filestatus:     onemeta.Filestatus,
			Filetype:       onemeta.Filetype,
		}
		//获取返回值，然后写入到这用户表
		ok := db.WriteFileMeta(tablename, filemeta)
		if ok {
			fmt.Println("写入成功")
			c.JSON(http.StatusOK, gin.H{"code": 200})
		}
		c.Abort()

	} else {
		//继续执行
		fmt.Println("唯一文件不存在，继续执行写入数据库")
		//不写在这里了
		c.Next()
	}

}
