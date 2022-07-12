package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"user_web/forms"
	"user_web/global"
	reponse "user_web/global/response"
	"user_web/middlewares"
	"user_web/models"
	"user_web/proto"
)

var userSrvClient proto.UserClient
var conn *grpc.ClientConn

func removeTopStruct(fields map[string]string) map[string]string {

	rsp := map[string]string{}
	for field, err := range fields {
		fmt.Printf("field:[%s],err:[%s]\n", field, err)

		fmt.Printf("位置:[%s]\n", field[strings.Index(field, ".")+1:])
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func HandleValitor(c *gin.Context, err error) {

	fmt.Println(err.Error())
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})

}

func GetUserlist(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	r, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{Pn: int32(pnInt), PSize: int32(pSizeInt)})
	if err != nil {
		zap.S().Errorw("【GetUserlist】查询【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	claims, _ := ctx.Get("claims")
	userId, _ := ctx.Get("userId")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户：%d", currentUser.ID)
	zap.S().Infof("访问用户2：%d", userId)
	result := make([]interface{}, 0)
	for _, value := range r.Data {
		user := reponse.UserResponse{
			Id:       value.Id,
			Mobile:   value.Mobile,
			NickName: value.NickName,
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
		}

		//data:=make(map[string]interface{})
		//
		//data["id"]=user.Id
		//data["mobile"]=user.Mobile
		//data["name"]=user.NickName
		//data["gender"]=user.Gender
		//data["birthday"]=user.BirthDay

		result = append(result, user)

	}
	ctx.JSON(http.StatusOK, result)

	fmt.Println(r.Data)
}

func PassWordLogin(c *gin.Context) {
	//InitGrpc()
	PassloginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&PassloginForm); err != nil {
		HandleValitor(c, err)
		return
	}

	//图像验证码验证，同一个包下面的变量可以共用
	//if !store.Verify(PassloginForm.CaptchaId, PassloginForm.Captcha, true) {
	//	c.JSON(http.StatusBadRequest, map[string]string{
	//		"captcha": "验证码错误",
	//	})
	//	return
	//}

	//登录的逻辑
	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: PassloginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"msg": e.Code().String(),
				})
			}
			return
		}
	} else {
		//只是查询到用户了而已，并没有检查密码
		if passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.CheckInfo{
			Password:          PassloginForm.Password,
			EncryptedPassword: rsp.Password,
		}); pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if passRsp.Success {

				//生成token
				j := middlewares.NewJWT()
				claim := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityID: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						Audience:  "",                              //观众
						ExpiresAt: time.Now().Unix() + 30*60*60*24, //一个月
						Issuer:    "zifeng6257",                    //观众
						NotBefore: time.Now().Unix(),
						Subject:   "hhaha",
					},
				}
				var token string
				if token, err = j.CreateToken(claim); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
				}

				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登录失败密码错误",
				})
			}
		}
	}

}

func RegisterForm(c *gin.Context) {
	//InitGrpc()
	RegisterForm := forms.RegisterForm{}
	if err := c.ShouldBind(&RegisterForm); err != nil {
		HandleValitor(c, err)
		return
	}

	//验证码
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	val2, err := rdb.Get(context.Background(), RegisterForm.Mobile).Result()
	if err == redis.Nil {
		zap.S().Errorf("Redis中的验证码错误: %s", val2)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		zap.S().Infof("Redis中的验证码: %s", val2)
		if val2 != RegisterForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}

	//注册的逻辑 查询手机是否存在
	if _, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: RegisterForm.Mobile,
	}); err == nil {
		c.JSON(http.StatusConflict, map[string]string{
			"msg": "手机号已经被注册",
		})
		return
	} else {
		if e, ok := status.FromError(err); ok {
			fmt.Println(e.Code())
		}
	}

	//只是查询到用户了而已，并没有检查密码
	if regRsp, pasErr := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   RegisterForm.Mobile,
		NickName: RegisterForm.Mobile,
		Password: RegisterForm.PassWord,
	}); pasErr != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	} else {
		//生成token
		j := middlewares.NewJWT()
		claim := models.CustomClaims{
			ID:          uint(regRsp.Id),
			NickName:    regRsp.NickName,
			AuthorityID: uint(regRsp.Role),
			StandardClaims: jwt.StandardClaims{
				Audience:  "",                              //观众
				ExpiresAt: time.Now().Unix() + 30*60*60*24, //一个月
				Issuer:    "zifeng6257",                    //观众
				NotBefore: time.Now().Unix(),
				Subject:   "hhaha",
			},
		}
		var token string
		if token, err = j.CreateToken(claim); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "生成token失败",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"id":         regRsp.Id,
			"nick_name":  regRsp.NickName,
			"token":      token,
			"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
		})

	}

}
