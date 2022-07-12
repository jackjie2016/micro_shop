package api

import (
	"fmt"
	validator "github.com/go-playground/validator/v10"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"goods_web/global"
)

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

func HandleValitorError(c *gin.Context, err error) {

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
