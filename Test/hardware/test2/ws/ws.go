package ws

import (
	"OldPackageTest/hardware/test2/global"
	"encoding/json"
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"net"
	"net/http"
	"strconv"
	"time"
	"log"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)
type WsServer struct {
	listener net.Listener
	addr     string
	upgrade  *websocket.Upgrader
	Logger *zap.Logger
}
type MsgStruct struct {
	AccessToken string `json:"access_token"`
	ExpireIn int64 `json:"expire_in"`
	Openid string  `json:"openid"`
	Unionid string `json:"unionid"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
}

func NewWsServer(addr string ,Logger *zap.Logger) *WsServer {
	ws := new(WsServer)
	ws.addr = addr
	ws.Logger = Logger
	ws.upgrade = &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if r.Method != "GET" {
				fmt.Println("method is not GET")
				return false
			}
			if r.URL.Path != "/ws" {
				fmt.Println("path error")
				return false
			}
			return true
		},
	}
	return ws
}

func (self *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ws" {
		httpCode := http.StatusInternalServerError
		reasePhrase := http.StatusText(httpCode)
		fmt.Println("path error ", reasePhrase)
		http.Error(w, reasePhrase, httpCode)
		return
	}
	conn, err := self.upgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("websocket error:", err)
		return
	}
	fmt.Println("client connect :", conn.RemoteAddr())
	go self.connHandle(conn)

}
func (self *WsServer) connHandle(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()
	stopCh := make(chan int)
	go self.send(conn, stopCh)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(5000)))
		_, msg, err := conn.ReadMessage()

		if err != nil {
			fmt.Println(err)
			close(stopCh)
			// 判断是不是超时
			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() {
					fmt.Printf("ReadMessage timeout remote: %v\n", conn.RemoteAddr())
					return
				}
			}
			// 其他错误，如果是 1001 和 1000 就不打印日志
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Printf("ReadMessage other remote:%v error: %v \n", conn.RemoteAddr(), err)
			}
			return
		}


		log.Printf("recv: %s", msg)


		var MsgS MsgStruct
		if err := json.Unmarshal(msg, &MsgS); err == nil {
			fmt.Println("收到消息MsgS：",MsgS)
			fmt.Println("收到消息MsgS：",MsgS.AccessToken)
		} else {
			fmt.Println(err)
		}

	}
}



func (self *WsServer) send(conn *websocket.Conn, stopCh chan int) {
	//self.send10(conn)

	redisClient:=global.GlobalRedi.Get()
	defer redisClient.Close()
	for {
		select {
		case <-stopCh:
			fmt.Println("connect closed")
			return
		case <-time.After(time.Second * 2):
			data:=map[string]string{"ping": "1"}
			sendData,err:=json.Marshal(data)
			if err!=nil{
				panic(err)
			}
			err = conn.WriteMessage(1, sendData)
			fmt.Println("sending....")
			if err != nil {
				fmt.Println("send msg faild ", err)
				return
			}
			//自己的逻辑 连接池冲突的解决
			//var Tailresult map[string]int
			//Tailresult=make(map[string]int)
			Tailresult,Ttotal,err:=GetListStatistics("Tailstatistics")
			if err!=nil{
				panic(err)
			}

			//var HeadResult map[string]int
			//HeadResult=make(map[string]int)
			HeadResult,Htotal,err:=GetListStatistics("Headstatistics")
			if err!=nil{
				panic(err)
			}

			data2:=map[string]interface{}{"HeadTotal": Htotal,"HeadResult":HeadResult,"TailTotal": Ttotal,"Tailresult": Tailresult}
			sendData2,err:=json.Marshal(data2)
			if err!=nil{
				panic(err)
			}

			fmt.Println("sending....",Tailresult,HeadResult)
			err = conn.WriteMessage(1, sendData2)
			if err != nil {
				fmt.Println("send msg faild ", err)
				return
			}

			time.Sleep(1500*time.Millisecond)
		}
	}
}

func (w *WsServer) Start() (err error) {

	w.listener, err = net.Listen("tcp", w.addr)
	if err != nil {
		w.Logger.Error("net listen error:"+string(err.Error()))
		return
	}
	err = http.Serve(w.listener, w)
	if err != nil {
		w.Logger.Error("http serve error:"+string(err.Error()))
		return
	}
	w.Logger.Info("Ws server started.", zap.String("addr", w.addr))
	return nil
}


func GetListStatistics( RedisKey string)(map[string]int,int,error)  {
	redisClient:=global.GlobalRedi.Get()
	defer redisClient.Close()
	var keys map[int]string
	var values map[int]int
	var response map[string]int
	keys = make(map[int]string)
	values = make(map[int]int)
	response = make(map[string]int)
	result, err := redigo.Values(redisClient.Do("hgetall", RedisKey))
	total:=0
	if err != nil {
		fmt.Println("hgetall failed", err.Error())
		return nil ,total,err
	} else {
		h:=0
		j:=0
		for i, v := range result {
			if i%2==0{
				keys[h]=string(v.([]byte))
				h++
			}else{
				vv, err := strconv.Atoi(string(v.([]byte)))

				if err != nil {
					values[j]=0
				} else {
					values[j]=vv
				}
				j++
			}

		}

		for k, vv := range keys {
			response[vv]=values[k]
			total=values[k]
		}
	}

	return response,total,nil
}
