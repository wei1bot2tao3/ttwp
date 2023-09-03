package meta

import (
	"fmt"
	"math"
	"path/filepath"
	"time"
)

// 文件信息
type FileMeta struct {
	FileHash       string
	FileName       string
	FileSize       int64
	FileCreateTime string
	FileLocation   string
}

// 通过filehas 查找 filemeta
var filemetas map[string]FileMeta

// 初始化filemetas
func init() {
	filemetas = make(map[string]FileMeta)
}

// 映射用户文件表
type CretaeFileMeta struct {
	Fileid         int `gorm:"primaryKey;autoIncrement"`
	Filehash       string
	Userid         int
	Filename       string
	Filesize       int64
	Filecreatetime string
	Filelocation   string
	Filestatus     int
	Filetype       string
	Parent         int
	Osslocation    string
	UpdatedAt      time.Time
}

// 返回值 暂时
type FormJson struct {
	Code int              `json:"code"`
	Data []CretaeFileMeta `json:"data"`
}

// 文件列表返回值
type Json struct {
	Code  int
	Msg   string
	Count int
	Data  []CretaeFileMeta
}

func (u *CretaeFileMeta) TableName() string {
	return "namefile"
}

// 映射用户表
type User struct {
	Userid    int
	Username  string
	Table     string
	Password  string
	Email     string
	Status    string
	Lsb       string
	Privilege string
	Phone     string
	CreatedAt time.Time
}

type Token struct {
	Id         int
	Username   string
	Token      string
	Singuptime time.Time
}

// 唯一文件表
type Onlyfile struct {
	Filehash       string
	Filename       string
	Filesize       int64
	Filecreatetime string
	Filelocation   string
	Filestatus     int
	UpdatedAt      time.Time
	Have           int
	Haveuser       string
	Filetype       string
}

// 判断文件类型
func Filetype(filename string) string {
	//获取文件名
	//解析后缀
	//进行判断
	fileType := filepath.Ext(filename)
	fmt.Println(fileType)
	switch fileType {
	case ".mp3", ".wav", ".cda", ".amr", ".m4a", ".ac3", ".wv", ".ape", "aiff":
		fmt.Println("音乐")
		return "音乐"
	case ".mp4", ".avi", ".mov", ".flv", ".wmv", ".mkv", ".rmvb":
		fmt.Println("视频")
		return "视频"
	case ".doc", ".docx", ".pdf", ".txt", ".ppt", ".excel", ".xmind":
		fmt.Println("文档")
		return "文档"
	case ".jpg", ".jpeg", ".png", "gif", "bmp", "webp":
		fmt.Println("图片")
		return "图片"
	default:
		fmt.Println("其他")
		return "其他"
	}
}

type NodeResourceUsage struct {
	Name        string
	CpuUsage    float64
	MemoryUsage float64
}

// TODO
// 暂时废弃
//
//	func Checkuserlocation(username string, nowuserlocation string) string {
//		if nowuserlocation == " " {
//			return username
//		}
//		return nowuserlocation
//	}
//
// 字节转换为GB
func FormatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	base := 1024.0
	if bytes == 0 {
		return "0 B"
	}
	i := math.Floor(math.Log(float64(bytes)) / math.Log(base))
	size := float64(bytes) / math.Pow(base, i)
	if size < 1 {
		size = 1
	}
	return fmt.Sprintf("%.2f %s", size, units[int(i)])
}

type Do struct {
	Userid    int
	Do        string
	CreatedAt time.Time
}
