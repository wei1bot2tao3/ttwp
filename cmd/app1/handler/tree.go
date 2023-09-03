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
	"time"
)

// 目录模块
// 生成目录
// 用树来存储目录结构

type FileMetadata struct {
	FileName string
	FileID   int
	ParentID int
}

type Directory struct {
	Name    string
	Subdirs []*Directory
	Files   []*FileMetadata
}

func BuildDirectory(username string, fileMetadataList []meta.CretaeFileMeta) *Directory {
	metadataMap := make(map[int]*FileMetadata)
	for _, metadata := range fileMetadataList {
		metadataMap[metadata.Fileid] = &FileMetadata{
			FileName: metadata.Filename,
			FileID:   metadata.Fileid,
			ParentID: metadata.Parent,
		}
	}
	username = "qwe"
	home := &Directory{
		Name:    username,
		Subdirs: []*Directory{},
		Files:   []*FileMetadata{},
	}

	for _, metadata := range metadataMap {
		parent := findDirectory(home, metadata.ParentID)
		if parent == nil {
			parent = home
		}

		subdir := findSubdir(parent, metadata.ParentID)
		if subdir == nil {
			subdir = &Directory{
				Name:    fmt.Sprintf("dir_%d", metadata.ParentID),
				Subdirs: []*Directory{},
				Files:   []*FileMetadata{},
			}
			parent.Subdirs = append(parent.Subdirs, subdir)
		}

		subdir.Files = append(subdir.Files, metadata)
	}

	return home
}

func findDirectory(dir *Directory, id int) *Directory {
	if id == 0 {
		return dir
	}
	for _, subdir := range dir.Subdirs {
		if subdir.Name == fmt.Sprintf("dir_%d", id) {
			return subdir
		}
		subdirResult := findDirectory(subdir, id)
		if subdirResult != nil {
			return subdirResult
		}
	}
	return nil
}

func findSubdir(dir *Directory, id int) *Directory {
	for _, subdir := range dir.Subdirs {
		if subdir.Name == fmt.Sprintf("dir_%d", id) {
			return subdir
		}
	}
	return nil
}

func printDirectory(dir *Directory, level int) {
	indent := ""
	for i := 0; i < level; i++ {
		indent += "    "
	}
	fmt.Printf("%s%s\n", indent, dir.Name)
	for _, file := range dir.Files {
		fmt.Printf("%s    %s\n", indent, file.FileName)
	}
	for _, subdir := range dir.Subdirs {
		printDirectory(subdir, level+1)
	}
}

// 创建目录
func Mkdirone(c *gin.Context) {
	//获取参数 用户名 当前路径
	username := c.Query("username")
	nowtree := c.Query("tree")
	maketree := c.Query("maketree")
	user, _ := db.GetUserMeta(username)
	tablename := UserTable(fmt.Sprint(username))
	pwd := fmt.Sprintf("/root/file/%d/", user.Userid)
	location := pwd + nowtree + "/" + maketree

	parent := c.Query("parentid")
	fmt.Println("获取到创建文件")
	fmt.Println(parent)
	parentid := ParentId(parent)
	fmt.Println(parentid)
	fmt.Printf("文件父亲id%d\n", parentid)
	//创建文件夹
	err := os.Mkdir(location, 0755)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": "0"})
		panic(err)
	}
	//oss
	parentmeta := db.GetParentLocation(tablename, parentid)
	ossloaction := parentmeta.Osslocation + "/" + maketree
	//计算哈希值
	metas := meta.CretaeFileMeta{
		Userid:         user.Userid,
		Filename:       maketree,
		Filelocation:   location,
		Filecreatetime: time.Now().Format("2006-01-02 15:04:06"),
		Filetype:       "文件夹",
		Parent:         parentid,
		Filestatus:     1,
		Osslocation:    ossloaction,
	}
	metas.Filehash = util.FolderHash(metas.Filelocation)

	db.WriteCretaeFile(tablename, metas)
	ossapi.Mkdir(ossloaction)

	//返回到前端
	c.JSON(http.StatusOK, gin.H{"ok": "ok"})
}

// 目录访问
func SelecTree(c *gin.Context) {
	//通过查询当前用户表格名
	//处理请求
	username := c.Query("username")
	//获取本目录id 然后检索父亲id是这个的
	parent := c.Query("tolocation")
	tree := c.Query("tree")
	filename := c.PostForm("filename")
	//获取 本身id tree路径 和文件名 把文件名当作tree返回
	parentid := ParentId(parent)
	location := tree + "/" + filename
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
	//目前parent，id

	c.JSON(200, gin.H{"code": forms.Code, "data": forms.Data, "token": token, "edlocation": location})
	//c.Writer.Write(jsons)
}

// TODO

//

//删除目录//判断目录里面是否有文件 循环删除
//修改目录名字
// 目录模块 检索目录模块
