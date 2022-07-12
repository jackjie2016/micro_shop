package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"io"
	"strings"
)

func genMd5(code ,salt string) string {
	Md5:=md5.New()
	_,_=io.WriteString(Md5,code+salt)
	return  hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	//dsn := "root:admin123@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second, // 慢 SQL 阈值
	//		LogLevel:      logger.Info, // Log level
	//		Colorful:      true,        // 禁用彩色打印
	//	},
	//)
	//
	//// 全局模式
	////NamingStrategy和Tablename不能同时配置，
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		//TablePrefix: "mxshop_",
	//		SingularTable: true,
	//	},
	//	Logger: newLogger,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//options := &password.Options{6, 100, 30, sha512.New}
	//salt, encodedPwd := password.Encode("admin123", options)
	//NewPassword:=fmt.Sprintf("$zifeng-sha512$%s$%s",salt, encodedPwd)
	//for i:=0;i<10;i++ {
	//	var user model.User
	//	user.NickName=fmt.Sprintf("zifeng %d",i)
	//	user.Mobile=fmt.Sprintf("1688888888%d",i)
	//	user.Password=NewPassword;
	//	db.Save(&user)
	//}
	////
	////db.AutoMigrate(&model.User{})
	//fmt.Println(genMd5("12345","44"))
	////34827ccb0eea8a706c4c34a16891f84e7b

	// Using custom options
	options := &password.Options{6, 100, 30, sha512.New}
	salt, encodedPwd := password.Encode("zifeng234", options)
	NewPassword:=fmt.Sprintf("$zifeng-sha512$%s$%s",salt, encodedPwd)
	fmt.Println(NewPassword) // true


	passwordinfo:=strings.Split(NewPassword,"$")
	fmt.Println(passwordinfo[2],passwordinfo[3])

	check := password.Verify("zifeng234", passwordinfo[2], passwordinfo[3], options)
	fmt.Println(check) // true
}
