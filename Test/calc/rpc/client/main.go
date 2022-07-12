package main

import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	)

type  ResponseData struct {
	Data int `json:"data"`
}

func Add(a,b int) int {
	req:=HttpRequest.NewRequest()
	res,_:=req.Get(fmt.Sprintf("http://127.0.0.1:8001/%s?a=%d&b=%d","add",a,b))
	body,_:=res.Body()
	repData := ResponseData{}
	_=json.Unmarshal(body,&repData)

	return repData.Data
}
/*
 最原始的rpc 场景demo 使用网络协议：http，数据协议：url
 */
func main() {
	fmt.Println(Add(2,2))
}