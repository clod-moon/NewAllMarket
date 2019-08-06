package market

import (
	"fmt"
	"AllMarket/model"
	"github.com/wonderivan/logger"
	"golang.org/x/net/websocket"

)

func GetBianMarket() {

	srcMarket,ok:=model.SrcMarketMap["bian"]
	if !ok{
		logger.Error("can not find bian")
		return
	}

	ws, err := websocket.Dial(srcMarket.WsUrl, "", Origin)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	fmt.Println("===============")
	var msg = make([]byte, 512000)

	for {
		m, err := ws.Read(msg)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		newmsg := msg[:m]

		fmt.Println("newmsg:",string(newmsg))
	}
}
