package market

import (
	"fmt"
	"NewAllMarket/model"
	"github.com/wonderivan/logger"
	"golang.org/x/net/websocket"
	"strconv"
	"encoding/json"
	"sync"
)

type Bian struct {
	Tick  string `json:"s"`
	Close string `json:"c"`
	Open  string `json:"o"`
	High  string `json:"h"`
	Low   string `json:"l"`
	Vol   string `json:"v"`
	Count string `json:"q"`
}

func swapBianMarket(b *model.Bian, tmpB *Bian) {
	b.Open,_ =strconv.ParseFloat(tmpB.Open, 64)
	b.Close,_ = strconv.ParseFloat(tmpB.Close, 64)
	b.High,_ = strconv.ParseFloat(tmpB.High, 64)
	b.Low,_ = strconv.ParseFloat(tmpB.Low, 64)
	b.Count,_ = strconv.ParseFloat(tmpB.Count, 64)
	b.Vol,_ =  strconv.ParseFloat(tmpB.Vol, 64)
	b.Rose, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (b.Open-b.Close)/b.Open*100+0.005), 64)
}

func GetBianMarket(wg *sync.WaitGroup) {

	defer wg.Done()

	srcMarket, ok := model.SrcMarketMap["bian"]
	if !ok {
		logger.Error("can not find bian")
		return
	}

	ws, err := websocket.Dial(srcMarket.WsUrl, "", Origin)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var msg = make([]byte, 512000)

	for {
		m, err := ws.Read(msg)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		newmsg := msg[:m]

		var data []Bian
		err = json.Unmarshal(newmsg, &data)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		for _, v := range data {
			modelBian ,ok:= model.BianMap[v.Tick]
			if !ok{
				logger.Error("can not find ticke:",v.Tick)
				continue
			}
			swapBianMarket(&modelBian,&v)
			modelBian.Update()
		}
	}
}
