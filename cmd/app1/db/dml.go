package db

import (
	"bysj0.3/cmd/app1/meta"
	"fmt"
)

//文件系统

// 查询类
// 查询返回
// 通过哈希查文件meta信息
func SelectOneMeta(tablename string, Filehash string) meta.CretaeFileMeta {
	var one meta.CretaeFileMeta
	Dbgorm().Table(tablename).Where("Filehash=?", Filehash).Take(&one)
	return one

}

// 查询当前路径下全部内容
func Showtabel(tablename string, parent int) ([]meta.CretaeFileMeta, error) {
	var tables []meta.CretaeFileMeta
	filemeta := Dbgorm().Table(tablename).Where("parent=? AND filestatus=?", parent, 1).Find(&tables)
	if filemeta == nil {
		fmt.Printf("读取失败：%v", err)
		return nil, nil
	}
	return tables, nil
}

// 根据文件id获取文件信息
func GetParentLocation(tablename string, parent int) meta.CretaeFileMeta {
	var tables meta.CretaeFileMeta
	filemeta := Dbgorm().Table(tablename).Where("fileid=?", parent).Find(&tables)
	if filemeta == nil {
		fmt.Println(filemeta.Error)
	}
	return tables

}

// 查询全部失效文件
func ShowDay7(tablename string) ([]meta.CretaeFileMeta, error) {
	var tables []meta.CretaeFileMeta
	filemeta := Dbgorm().Table(tablename).Where("filestatus=?", 0).Find(&tables)
	if filemeta == nil {
		fmt.Printf("读取失败：%v", err)
		return nil, nil
	}
	return tables, nil
}

// 查询全部文件
func AllFile(tablename string) ([]meta.CretaeFileMeta, error) {
	var tables []meta.CretaeFileMeta
	filemeta := Dbgorm().Table(tablename).Where("filestatus=?", 1).Find(&tables)
	if filemeta == nil {
		fmt.Printf("读取失败：%v", err)
		return nil, nil
	}
	return tables, nil
}

// 模糊查询文件名
func SearchFile(tablename string, key string) (metas []meta.CretaeFileMeta) {
	var mates []meta.CretaeFileMeta
	ok := Dbgorm().Table(tablename).Where("filename LIKE ?", "%"+key+"%").Find(&mates)
	if ok.Error != nil {
		fmt.Println("检索失败:", ok.Error)
	}
	return mates
}

// 根据类型查文件
func Searchtype(tablename string, peach string) (metas []meta.CretaeFileMeta) {
	var tables []meta.CretaeFileMeta
	filemeta := Dbgorm().Table(tablename).Where("filetype=?", peach).Find(&tables)
	if filemeta == nil {
		fmt.Printf("读取失败：%v", err)
		return nil
	}
	return tables
}

// 检查类
// 检查是否存在唯一哈希值
func CheckOnlyFiles(Filehash string) (meta.Onlyfile, bool) {
	one := meta.Onlyfile{}
	fmt.Println(Filehash)
	ok := Dbgorm().Where("Filehash=?", Filehash).Find(&one)
	fmt.Printf("查询结果为：%v\n", ok.RowsAffected)
	if ok.RowsAffected == 0 {
		//不存在
		return one, false
	}
	//存在
	return one, true
}

// 修改类
// 写入
// 上传文件meta信息写入用户数据库
func WriteCretaeFile(table string, m meta.CretaeFileMeta) bool {
	metas := &meta.CretaeFileMeta{
		Userid:         m.Userid,
		Filehash:       m.Filehash,
		Filename:       m.Filename,
		Filesize:       m.Filesize,
		Filecreatetime: m.Filecreatetime,
		Filelocation:   m.Filelocation,
		Filestatus:     1,
		Filetype:       m.Filetype,
		Parent:         m.Parent,
		Osslocation:    m.Osslocation,
	}
	fmt.Println(metas)
	fmt.Println("开始写入数据库")
	if err := Dbgorm().Table(table).Create(metas).Error; err != nil {
		fmt.Printf("插入失败:%v", err)
		return false
	}

	return true
}

// 直接写入用户文件表 秒传
func WriteFileMeta(tablename string, fileMeta meta.CretaeFileMeta) bool {
	if err := Dbgorm().Table(tablename).Create(fileMeta).Error; err != nil {
		fmt.Printf("插入失败:%v", err)
		return false
	}
	return true
}

// 上传文件写入唯一文件表格
func WriteOnlyFiles(onemeta meta.Onlyfile) bool {
	onemeta.Have = 1
	//插入
	ok := Dbgorm().Create(&onemeta)
	if ok.Error != nil {
		fmt.Printf("插入失败错误是：%v", err)
		return false
	}
	fmt.Println("唯一文件插入成功")
	return true

}

// 修改
// 修改文件名字以及获取修改时间
func EditFilename(tablename string, Filehash string, Filenewname string, Ossloaction string) (bool, error) {
	//先基于哈希值查询
	table := meta.CretaeFileMeta{
		Filehash: Filehash,
	}

	//根据文件哈希值去修改文件名
	Dbgorm().Table(tablename).Model(&table).Where("filehash=?", Filehash).Update("Filename", Filenewname)
	Dbgorm().Table(tablename).Model(&table).Where("filehash=?", Filehash).Update("Osslocation", Ossloaction)

	return true, nil
	// 使用 map 更新多个属性，只会更新其中有变化的属性
	//db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	//UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

}

// 删除
// 删除文件在数据库的信息
func Deletemeta(tablename string, delhash string) (bool, error) {
	hash := meta.CretaeFileMeta{}
	filedelete := Dbgorm().Table(tablename).Where("Filehash=?", delhash).Delete(&hash)
	if filedelete.Error != nil {
		fmt.Printf("删除失败，err:%v", err)
		return false, err
	}
	return true, nil
}

func Delete07(tablename string, Filehash string, value int) (bool, error) {
	//先基于哈希值查询
	table := meta.CretaeFileMeta{
		Filehash: Filehash,
	}
	//根据文件哈希值去修改文件名
	ok := Dbgorm().Table(tablename).Model(&table).Where("filehash=?", Filehash).Update("filestatus", value)
	return true, ok.Error
}

//管理员 修改用户数据 修改 密码 or 邮箱 or 用户名
//用户修改自己的信息
//管理员重置用户密码

//
