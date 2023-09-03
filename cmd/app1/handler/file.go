package handler

import (
	"bysj0.3/cmd/app1/db"
	"bysj0.3/cmd/app1/meta"
	ossapi "bysj0.3/cmd/app1/oss"
	"bysj0.3/cmd/app1/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

func DeleteFileTO(c *gin.Context) {
	//接收参数
	username := c.PostForm("username")
	filehash := c.PostForm("Filehash")
	tablename := UserTable(fmt.Sprint(username))

	//修改数据库参数 文件状态为0
	ok, err := db.Delete07(tablename, filehash, 0)
	if !ok {
		fmt.Println("临时删除失败")
		fmt.Println(err)
	}
	//然后返回
	onemeta := db.SelectOneMeta(tablename, filehash)
	db.WriteDo(onemeta.Userid, "把文件"+onemeta.Filename+"放入回收站")
	c.JSON(http.StatusOK, nil)
}

// 恢复
func DeleteFileT1(c *gin.Context) {
	//接收参数
	username := c.PostForm("username")
	filehash := c.PostForm("Filehash")
	tablename := UserTable(fmt.Sprint(username))
	//修改数据库参数 文件状态为0
	ok, err := db.Delete07(tablename, filehash, 1)
	if !ok {
		fmt.Println("恢复成功")
		fmt.Println(err)
	}
	//然后返回
	onemeta := db.SelectOneMeta(tablename, filehash)
	db.WriteDo(onemeta.Userid, "把文件"+onemeta.Filename+"取出回收站")
	c.JSON(http.StatusOK, nil)
}

// 分享
func ShareFile(c *gin.Context) {
	//获取当前用户文件id 用户名
	//
}

// 上传文件
// 接收文件
func UpfilePost(c *gin.Context) {
	//获取当前用户名
	username := c.Query("username")
	tablename := UserTable(fmt.Sprint(username))
	user, _ := db.GetUserMeta(username)
	//获取当前的parentid
	parent := c.Query("parent")
	parentid := ParentId(parent)
	//获取用户表名
	//接收文件
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("上传文件失败：%v", err)
		return
	}
	//组装meta
	//通过parentid 获取父亲的osslocation+filename
	parentmeta := db.GetParentLocation(tablename, parentid)
	osslocation := parentmeta.Osslocation + "/" + file.Filename
	//更改接口
	filelocation := fmt.Sprintf("/root/file/" + strconv.Itoa(user.Userid) + "/home/" + file.Filename)
	metas := meta.CretaeFileMeta{
		Userid:         user.Userid,
		Filename:       file.Filename,
		Filesize:       file.Size,
		Filecreatetime: time.Now().Format("2006-01-02 15:04:06"),
		Filelocation:   filelocation,
		Filestatus:     1,
		Parent:         parentid,
		Osslocation:    osslocation,
	}
	c.SaveUploadedFile(file, metas.Filelocation)
	//通过文件路径计算哈希值
	metas.Filehash, _ = util.FileSha1(metas.Filelocation)
	//判断文件类型
	metas.Filetype = meta.Filetype(file.Filename)
	//上传到oss
	if ossapi.UpFile(metas.Filelocation, metas.Osslocation) {
		//删除文本地文件
		err := os.Remove(metas.Filelocation)
		if err != nil {
			fmt.Print("本地删除失败：%v", err)
			return
		}
	}
	//写入数据库根据文件表
	//db.CreateFileMeta(tablename, userid.Userid, CretaeMeta, "home")
	db.WriteCretaeFile(tablename, metas)
	//返回json
	data := meta.Json{
		Code: http.StatusOK,
		Data: nil,
	}
	db.WriteDo(user.Userid, "把文件"+metas.Filename+"上传")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": data})

}

// 下载文件
func DownFile(c *gin.Context) {
	//获取文件名
	filehash := c.PostForm("Filehash")
	username := c.PostForm("username")
	tablename := UserTable(fmt.Sprint(username))
	//通过数据库查询meta
	onemeta := db.SelectOneMeta(tablename, filehash)
	url := ossapi.GetOSSFileUrl(onemeta.Osslocation)

	db.WriteDo(onemeta.Userid, "请求下载"+onemeta.Filename+"")
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "url": url})

}

// 修改文件名
func InstallFileMeta(c *gin.Context) {
	//获取文件哈希值和新名字
	username := c.PostForm("username")
	tablename := UserTable(fmt.Sprint(username))
	filehash := c.PostForm("Filehash")
	filenewname := c.PostForm("Filenewname")
	onemeta := db.SelectOneMeta(tablename, filehash)
	fmt.Println()
	parent := c.PostForm("parent")
	parentid := ParentId(parent)
	//通过parentid 获取父亲的osslocation+filename
	parentmeta := db.GetParentLocation(tablename, parentid)
	newossloaction := parentmeta.Osslocation + "/" + filenewname
	fmt.Println(onemeta.Osslocation)
	fmt.Println(newossloaction)
	if ossapi.NewFile(onemeta.Osslocation, newossloaction) {
		//进去数据库修改

	}
	ok, err := db.EditFilename(tablename, filehash, filenewname, newossloaction)
	if ok != true {
		panic(err)
	}
	//返回结果

	db.WriteDo(onemeta.Userid, "把文件"+onemeta.Filename+"修改为："+filenewname)
	c.JSON(http.StatusOK, gin.H{"code": 200})

}

// 删除文件
// 彻底删除TODO删除文件
func DeleteFile(c *gin.Context) {
	//获取form表单参数
	username := c.PostForm("username")
	filehash := c.PostForm("Filehash")
	tablename := UserTable(fmt.Sprint(username))
	//数据库查询返回全部信息
	onemeta := db.SelectOneMeta(tablename, filehash)
	if onemeta.Filetype != "wjj" {
		if ossapi.RmFile(onemeta.Osslocation) {
			//删除数据库
			ok, err := db.Deletemeta(tablename, filehash)
			if ok != true {
				panic(err)
			}
		}
	} else {
		if ossapi.RmMkdir(onemeta.Osslocation) {
			//删除数据库
			ok, err := db.Deletemeta(tablename, filehash)
			if ok != true {
				panic(err)
			}
		}
	}

	db.WriteDo(onemeta.Userid, "把文件"+onemeta.Filename+"彻底删除")
	//返回值
	c.JSON(http.StatusOK, nil)
}
