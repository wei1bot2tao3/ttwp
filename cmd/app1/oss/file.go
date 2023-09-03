package ossapi

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

//

// 创建目录
func Mkdir(name string) bool {
	//获取存储空间
	bucket, _ := CreateOSSClientBucket()
	//创建目录 如果存在会覆盖原始文件
	//按照格式拼目录
	exampledir := name + "/"
	fmt.Println(exampledir)
	err := bucket.PutObject(exampledir, bytes.NewReader([]byte("")))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// 删除目录
func RmMkdir(name string) bool {
	//获取存储空间
	bucket, _ := CreateOSSClientBucket()

	//按照格式拼目录
	exampledir := name + "/"
	//：初始化一个Marker对象，用于标记从哪里开始列举对象。在这里，我们将其设置为空字符串，表示从头开始列举。
	marker := oss.Marker("")
	//初始化一个Prefix对象，用于指定前缀。在这里，我们将其设置为“src”，表示列举所有前缀为“src”的文件。
	//如果您仅需要删除src目录及目录下的所有文件，则prefix设置为src/
	//for {...}：使用for循环列举所有包含指定前缀的文件并删除。
	prefix := oss.Prefix(exampledir)
	isExist, err := bucket.IsObjectExist(exampledir)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	if !isExist {
		return false
	}

	for {
		//从头开始列举 ListObjects：从哪开始，指定前缀 返回给ListObjectsResult
		//ListObjectsResult包含 ：Objects：表示列举到的所有对象，是一个ObjectInfo的切片
		count := 0
		//因为bucket.ListObjects默认一千个所
		lor, err := bucket.ListObjects(marker, prefix)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}

		var objects []string
		for _, object := range lor.Objects {
			//object.Key对象名称
			objects = append(objects, object.Key)
		}

		// 将oss.DeleteObjectsQuiet设置为true，表示不返回删除结果。
		delRes, err := bucket.DeleteObjects(objects, oss.DeleteObjectsQuiet(true))
		if err != nil {
			fmt.Println("Error报错是:", err)
			os.Exit(-1)
		} //成功删除的对象信息delRes.DeletedObjects

		//判断是否删除干净
		if len(delRes.DeletedObjects) > 0 {
			fmt.Println("these objects deleted failure,", delRes.DeletedObjects)
			fmt.Println("删除失败")
			break
		}
		//进行分页查询
		count += len(objects)
		prefix = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}
	return true
}

// 上传文件
func UpFile(filelocation string, osslocation string) bool {
	bucket, _ := CreateOSSClientBucket()
	fmt.Println(osslocation)
	fmt.Println(filelocation)
	err := bucket.PutObjectFromFile(osslocation, filelocation)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return true
}

// 下载文件
func GetOSSFileUrl(filename string) string {
	bucket, _ := CreateOSSClientBucket()
	url, err := bucket.SignURL(filename, oss.HTTPGet, 3600)
	if err != nil {
		// 处理错误
	}

	return url
}

func GetOSSFileUrlTime(filename string, time int64) string {
	bucket, _ := CreateOSSClientBucket()
	url, err := bucket.SignURL(filename, oss.HTTPGet, time)
	if err != nil {

	}

	return url
}

// 删除文件
func RmFile(filename string) bool {
	//获取存储空间
	bucket, _ := CreateOSSClientBucket()
	err := bucket.DeleteObject(filename)
	if err != nil {
		HandleError(err)
		return false
	}
	return true
}

// 修改文件名
func NewFile(old string, new string) bool {
	bucket, _ := CreateOSSClientBucket()
	object, err := bucket.CopyObject(old, new)
	if err != nil {
		fmt.Println(object.XMLName)
		return false

	}
	RmFile(old)
	return true

}
