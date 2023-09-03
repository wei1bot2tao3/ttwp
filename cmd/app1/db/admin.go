package db

import (
	"bysj0.3/cmd/app1/meta"
)

func AdminPassword(userid int) bool {
	Usermeta := meta.User{}
	Dbgorm().Model(&Usermeta).Where("userid=?", userid).Update("password", 123456)
	return true
}
func AdminSatus(userid int, status int) bool {
	Usermeta := meta.User{}
	Dbgorm().Model(&Usermeta).Where("userid=?", userid).Update("status", status)
	return true
}

func AdminSearchUser(name string) []meta.User {
	Usermeta := []meta.User{}
	Dbgorm().Where("username LIKE ?", "%"+name+"%").Find(&Usermeta)
	return Usermeta
}

func AdminGetAllHistory() []meta.Do {
	history := []meta.Do{}
	Dbgorm().Find(&history)
	return history
}

func AdminDeleteUser(user int) bool {
	Usermeta := []meta.User{}
	Dbgorm().Where("userid=?", user).Delete(&Usermeta)

	return true

}
