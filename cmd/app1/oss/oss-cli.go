package ossapi

import (
	"bysj0.3/cmd/app1/meta"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func HandleError(err error) {
	fmt.Println("Error:", err)
	//写入日志
}

// 直接返回一个桶
func CreateOSSClientBucket() (*oss.Bucket, error) {
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	endpoint := "https://oss-cn-beijing.aliyuncs.com"
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := "LTAI5tQ6x3CUBMCdENuJcE9E"
	accessKeySecret := "DAQWQj4lT08a85nX5aA9fBSrqn8VgK"
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	//返回存储空间
	// 使用client进行其他操作
	// 获取存储空间。
	bucketName := "ttwangpan"
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
	}
	return bucket, err

}

func GetOss() string {
	endpoint := "https://oss-cn-beijing.aliyuncs.com"
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := "LTAI5tQ6x3CUBMCdENuJcE9E"
	accessKeySecret := "DAQWQj4lT08a85nX5aA9fBSrqn8VgK"
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	// yourBucketName填写Bucket名称。
	// yourInventoryId填写清单规则名称。
	stat, err := client.GetBucketStat("ttwangpan")
	if err != nil {
		HandleError(err)
	}
	// 获取Bucket的总存储量，单位为字节。
	gbSize := meta.FormatBytes(stat.Storage)
	return gbSize
	// 获取标准存储类型Object的存储量，单位为字节。

}
