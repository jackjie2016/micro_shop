package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
)

type LoginForm struct {
	User string `json:"user" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=5,max=10"`
}
type SignupForm struct {
	Age uint8 `json:"age" binding:"required,gte=1,lte=130"`
	Name string `json:"name" binding:"required,min=3,max=10"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5,max=10"`
	RePassword string `json:"repassword" binding:"required,eqfield=Password"`
}
var trans ut.Translator
func InitTrans(locale string)(err error)  {
	//修改gin框架中的validator 引擎属性，实现定制
	if v,ok:=binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name:=strings.SplitN(fld.Tag.Get("json"),",",2)[0]
			if name=="-"{
				return ""
			}
			return name
		})
		zhT:=zh.New()//中文翻译器
		enT:=en.New()//英文翻译器
		//第一个参数是备用的语言环境，后面的应该支持的语言环境
		uni:=ut.New(enT,zhT,zhT)
		trans,ok=uni.GetTranslator(locale)
		if !ok{
			return fmt.Errorf("uni.GetTranslator(%s)",locale)
		}
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v,trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v,trans)
		default:
			en_translations.RegisterDefaultTranslations(v,trans)
		}
		return
	}
	return
}

func removeTopStruct(fields map[string]string) map[string]string  {

	rsp:=map[string]string{}
	for field,err:=range fields{
		fmt.Printf("field:[%s],err:[%s]\n",field,err)

		fmt.Printf("位置:[%s]\n",field[strings.Index(field,".")+1:])
		rsp[field[strings.Index(field,".")+1:]]=err
	}
	return rsp
}
func main(){
	if err:=InitTrans("zh");err!=nil{
		fmt.Printf("初始化翻译器错误")
		return
	}
	 router:=gin.Default()
	 
	 router.POST("/loginJSON", func(c *gin.Context) {
		 var loginForm LoginForm
		 if err:=c.ShouldBind(&loginForm);err!=nil{
		 	fmt.Println(err.Error())
		 	errs,ok:=err.(validator.ValidationErrors)
		 	if !ok{
		 		c.JSON(http.StatusOK,gin.H{
		 			"msg":err.Error(),
				})
			}

		 	c.JSON(http.StatusBadRequest,gin.H{
		 		"error":removeTopStruct(errs.Translate(trans)),
			})
			 return
		 }
		 c.JSON(http.StatusBadRequest,gin.H{
			 "msg":"登录成功",
		 })
	 })

	router.POST("/signup", func(c *gin.Context) {
		var SignupForm SignupForm

		if err:=c.ShouldBind(&SignupForm);err!=nil{

			errs,ok:=err.(validator.ValidationErrors)
			if !ok{
				c.JSON(http.StatusOK,gin.H{
					"msg":err.Error(),
				})
			}
			
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest,gin.H{
				"error":errs.Translate(trans),
			})
			return
		}
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":"注册成功",
		})
	})

	 router.Run(":8082")
}

