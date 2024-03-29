package main
import (
	"fmt"
	"os"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)
func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}
func main() {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyId := "LTAI5tFc8pbqU7WyheF8LkdE"
	accessKeySecret := "MgRvwEIYecE7xnCe3t5atsmaSvM9gu"
	bucketName := "ceshixiaa2"

	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	objectName := "goods/1.bmp"

	// 创建OSSClient实例。
	//client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	//if err != nil {
	//	handleError(err)
	//}
	//// 创建存储空间。
	//err = client.CreateBucket(bucketName)
	//if err != nil {
	//	handleError(err)
	//}



	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := `F:\code\21022250.bmp`
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		//panic(err.Error())
		handleError(err)
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		handleError(err)
	}
}