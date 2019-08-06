package httpserver

import (
	"net/http"
	"golang.org/x/net/websocket"
	"github.com/wonderivan/logger"
)
var messageChan  = make(chan []byte,1024)

func EchoServer(ws *websocket.Conn) {
	var err error
	for {
		//var reply string

		//if err = websocket.Message.Receive(ws, &reply); err != nil {
		//	fmt.Println(err)
		//	continue
		//}
		if err = websocket.Message.Send(ws, <-messageChan); err != nil {
			logger.Error(err.Error())
			continue
		}
	}
}

func GetMarketHandler() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}