package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func init() {
	// MySQL 配置信息
	username := "test"
	//username := "user_file_DML"                               // 账号
	password := "mutqa8-sikkiw-Vakrof"                        // 密码
	host := "bj-cynosdbmysql-grp-f1wykius.sql.tencentcdb.com" // 地址
	port := 20593                                             // 端口
	DBname := "File_meta"                                     // 数据库名称
	timeout := "10s"                                          // 连接超时，10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)
	// Open 连接
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql.")
	}

}

func Dbgorm() *gorm.DB {
	fmt.Println("数据库链接成功")

	return db
}
