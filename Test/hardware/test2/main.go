package main

import (
	"OldPackageTest/hardware/test2/ws"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/namsral/flag"
	modbus "github.com/thinkgos/gomodbus"
	"github.com/thinkgos/gomodbus/mb"
	"go.uber.org/zap"

	"OldPackageTest/hardware/test2/global"
	"OldPackageTest/hardware/test2/initialize"
	"OldPackageTest/hardware/test2/server"
)

var Haddr = flag.String("Haddr", "COM3", "address to listen")
var Taddr = flag.String("Taddr", "COM4", "address to listen")
var wsAddr = flag.String("ws_addr", ":9090", "websocket address to listen")
var addr = flag.String("addr", ":8080", "address to listen")
func main(){
	flag.Parse()

	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	//redis Init
	initialize.InitRedis()



	//调用RTUClientProvider的构造函数,返回结构体指针
	p := modbus.NewRTUClientProvider()





	//fmt.Println(p)



		//waiting := make(chan int)
		logger.Info("Ws开始工作====")
		go wsStart(*wsAddr,logger)

		logger.Info("计数器开始工作====")
		//go Location(*Haddr,p,true)
		go Location(*Haddr,p,true)
		//<-waiting

	r:=gin.Default()
	r.LoadHTMLGlob("template/**/*")//加载当前目录下所有的文件
	r.Static("/assets","./assets")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK,"index/index.html",gin.H{
			"title":"首页",
		})
	})
	r.GET("api/index/getPlanData", getPlanData)

	r.Run(":8082")


}
func wsStart(wsAddr string ,logger *zap.Logger)  {
	wsServer := ws.NewWsServer(wsAddr,logger)
	err := wsServer.Start()
	if err!=nil{
		logger.Error(err.Error())
	}
}

func getPlanData(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{
		"title":"商品",
	})
}

type RedisKey struct {
	LockName string
	QueueName string
	StatisticsName string
}

func mux()  {
	
}
func Location(Address  string,p *modbus.RTUClientProvider,isHead bool)  {
	var KeyTable RedisKey
	fmt.Println("====read======")
	if isHead{
		KeyTable=RedisKey{
			LockName: "HeadLock",
			QueueName: "HeadRfidList",
			StatisticsName: "Headstatistics",
		}
	}else{
		KeyTable=RedisKey{
			LockName: "TailLock",
			QueueName: "TailRfidList",
			StatisticsName: "Tailstatistics",
		}
	}
	fmt.Println("====key======",KeyTable)
	go Calculator(Address,p,KeyTable)
	go GetNowRFID(KeyTable)
}

//计算器
func Calculator(Address  string,p *modbus.RTUClientProvider,KeyTable RedisKey)  {

	p.Address = Address

	p.BaudRate = 9600

	p.DataBits = 8

	p.Parity = "N"

	p.StopBits = 1

	p.Timeout = 100 * time.Millisecond

	client := mb.NewClient(p)

	client.LogMode(false)

	err := client.Start()

	if err != nil {
		panic(err)
	}


	redisClient:=global.GlobalRedi.Get()
	defer redisClient.Close()
	for {
		value, err := client.ReadHoldingRegisters(1, 0, 1)

		if err != nil {
			fmt.Println("readHoldErr,", err)
		} else {
			for _,v:= range value {
				lock_name,err := redisClient.Do("Get",KeyTable.LockName)
				if err!=nil{
					fmt.Println(err.Error())
				}

				if lock_name!=nil{
					_,err=redisClient.Do("hset",KeyTable.StatisticsName,fmt.Sprintf("%s",lock_name),
						v)
					if err != nil {
						panic(err)
					}

				}
			}
		}
	}
}

//RFID 锁
func GetNowRFID(KeyTable RedisKey)  {
	redisClient:=global.GlobalRedi.Get()
	defer redisClient.Close()
	for  {
		//当前rfid值
		node,err := redigo.String(redisClient.Do("rpop",KeyTable.QueueName))
		if err!=nil{
			if err == redigo.ErrNil {
				fmt.Println("rfidList nil")
				continue
			}else{
				panic(err)
			}
		}

		if node!=""{

			lock_name,err := redisClient.Do("Get",KeyTable.LockName)
			if err!=nil{
				fmt.Println(err.Error())
			}

			if lock_name==nil {
				//5秒过期
				_,err = redisClient.Do("Set",KeyTable.LockName,node,"EX",5)
				if err != nil {
					panic(err)
				}
			}
			if fmt.Sprintf("%s",lock_name)==node{
				//5秒过期
				_,err = redisClient.Do("EXPIRE",KeyTable.LockName,5)
				if err != nil {
					panic(err)
				}
			}
		}else{
			fmt.Println("当前队列空")
		}
		time.Sleep(time.Millisecond*100)
	}
}