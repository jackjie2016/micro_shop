package main

import (
	"log"
	"net/http"

	"github.com/namsral/flag"

	"go.uber.org/zap"

	"OldPackageTest/hardware/test2/server"
)
var wsAddr = flag.String("ws_addr", ":9091", "websocket address to listen")

func main()  {

	flag.Parse()

	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	addr := *wsAddr
	logger.Info("Ws server started.", zap.String("addr", addr))
	logger.Sugar().Fatal(
		http.ListenAndServe(addr, nil),
	)

}
