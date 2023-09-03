package util

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func FolderHash(folderPath string) string {
	var files []string
	var folderInfo os.FileInfo
	var err error

	// 获取文件夹内所有文件的路径
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		} else if folderInfo == nil {
			folderInfo = info
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	// 对文件路径进行排序
	sort.Strings(files)

	// 计算文件哈希值
	h := sha256.New()
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		h.Write(data)
	}

	// 添加文件夹信息
	h.Write([]byte(folderInfo.Name()))
	h.Write([]byte(folderInfo.ModTime().String()))
	h.Write([]byte(folderInfo.Mode().String()))

	// 返回哈希值
	return fmt.Sprintf("%x", h.Sum(nil))
}
