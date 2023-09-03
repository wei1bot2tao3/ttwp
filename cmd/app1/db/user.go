package db

import (
	"bysj0.3/cmd/app1/meta"
	"fmt"
	"time"
)

// 创建数据库token信息
func CreateToken(username string, token string) bool {
	tokens := &meta.Token{
		Username:   username,
		Token:      token,
		Singuptime: time.Now(),
	}
	//把token写入数据库
	result := Dbgorm().FirstOrCreate(&tokens, &meta.Token{Username: username})
	if err := result.Error; err != nil {
		return false
	}
	return true
}

// 通过用户名来返回所有数据
func GetUserMeta(name string) (meta.User, bool) {
	Usermeta := meta.User{}
	Dbgorm().Where("username=?", name).Find(&Usermeta)
	return Usermeta, true

}

// 通过id返回
func GetIdUserMeta(userid int) (meta.User, bool) {
	Usermeta := meta.User{}
	Dbgorm().Where("username=?", userid).Find(&Usermeta)
	return Usermeta, true

}

// 检查是否重名
// 修改登录密码
func NewPassword(name string, password string) bool {
	Usermeta := meta.User{}
	Dbgorm().Model(&Usermeta).Where("username=?", name).Update("password", password)
	return true
}

// 用户相关内容
// 更新用户信息
func UpdateUserMeta(userid int, updates ...interface{}) error {
	updateMap := make(map[string]interface{})
	for i := 0; i < len(updates); i += 2 {
		updateMap[updates[i].(string)] = updates[i+1]
	}

	err := Dbgorm().Model(&meta.User{}).Where("userid = ?", userid).Updates(updateMap).Error
	return err
}

// 检查用户
func CheckUser(name string) bool {
	Usermeta := meta.User{}
	ok := Dbgorm().Where("username=?", name).First(&Usermeta)
	if ok.RowsAffected > 0 {
		return false
	} else {
		return true
	}
}

// TODO
// 从DML移动用户相关
// 获取全部用户
func AllUser(privilege string) []meta.User {
	var alluser []meta.User
	ok := Dbgorm().Where("privilege = ?", privilege).Find(&alluser).Error
	if ok != nil {
		fmt.Println("用户查询失败")
		return nil
	}
	return alluser
}

// 获取全部封禁用户
func AllXhw(privilege int) []meta.User {
	var alluser []meta.User
	ok := Dbgorm().Where("status= ?", privilege).Find(&alluser).Error
	if ok != nil {
		fmt.Println("用户查询失败")
		return nil
	}
	return alluser
}

// 用户系统
// 注册用户 写入数据库用户表
func WriteUser(user meta.User) int {
	//获取前端信息
	user.CreatedAt = time.Now()
	result := Dbgorm().FirstOrCreate(&user, meta.User{Username: user.Username})
	fmt.Print("数据库开始操作 判断返回值")
	if err := result.Error; err != nil {
		return 0
	}
	return 1
}

// 注册用户创建用户文件表
func CreateUserTable(tablename string) {
	//创建用户表
	err := Dbgorm().Table(tablename).AutoMigrate(&meta.CretaeFileMeta{})
	fmt.Println(err)

}

func WriteDo(userid int, some string) {
	Do := meta.Do{
		Userid: userid,
		Do:     some,
	}
	Dbgorm().Create(&Do)

}

// history
func UserHsitory(userid string) []meta.Do {
	var History []meta.Do
	Dbgorm().Where("userid=?", userid).Find(&History)
	return History
}

// 全部
func AllUserHsitory(userid string) []meta.Do {
	var History []meta.Do
	Dbgorm().Where("user_id=?", userid).Find(&History)
	return History
}
