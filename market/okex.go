package market

import (
	"AllMarket/model"
	//"regexp"
	"github.com/wonderivan/logger"
	"encoding/json"
	"strconv"
	"time"
	"fmt"
	"golang.org/x/net/websocket"
)

var (
//regx = regexp.MustCompile(`market.([a-zA-Z]+).detail`)

//OkexEndpoint = "wss://real.okex.com:10442/ws/v3"
)

type Okex struct {
	Open      string    `json:"open_24h"`
	Close     string    `json:"last" `
	Low       string    `json:"low_24h"`
	Tick      string    `json:"instrument_id"`
	Timestamp time.Time `json:"timestamp"`
}

type OkexResp struct {
	Table string `json:"table"`
	Data  []Okex `json:"data"`
}

func swapOkexMarket(o *model.Okex, tmpO *Okex) {
	o.Open, _ = strconv.ParseFloat(tmpO.Open, 64)
	o.Close, _ = strconv.ParseFloat(tmpO.Close, 64)
	o.Low, _ = strconv.ParseFloat(tmpO.Low, 64)
	o.Rose, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (o.Open-o.Close)/o.Open*100+0.005), 64)
}

func OkexPing(ws *websocket.Conn,pong int64){
	send([]byte(fmt.Sprintf( `{"pong":%d}`,pong)),ws)
	//time.Now().UnixNano() / int64(time.Millisecond))
}

func GetOkexMarket() {

	srcMarket,ok:=model.SrcMarketMap["okex"]
	if !ok{
		logger.Error("can not find okex")
		return
	}

	ws, err := websocket.Dial(srcMarket.WsUrl, "", Origin)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ping := time.Now().UnixNano() / int64(time.Millisecond)
	send([]byte(fmt.Sprintf( `{"ping":%d}`,ping)),ws)

	logger.Error("OkexTickerList:",len(OkexTickerList))

	//for _, ticker := range OkexTickerList {
	ticker := "ETH-USD"
		go func(ticker string){
			message := []byte(fmt.Sprintf(srcMarket.Template,ticker))
			send(message, ws)
		}(ticker)
	//}
	var msg = make([]byte, 512000)

	for {
		m, err := ws.Read(msg)
		if err != nil {
			logger.Error(err)
			continue
		}
		newmsg := msg[:m]

		unzipmsg, _ := GzipDecode(newmsg)

		fmt.Printf("Receive[UNZIP]: [%d:%d] %s\n", m, len(unzipmsg), unzipmsg[:])
		var okexResp OkexResp
		err = json.Unmarshal(unzipmsg,&okexResp)
		if  err != nil{
			logger.Error(err.Error())
			continue
		}

		if okexResp.Table == "index/ticker"{
			for _,o:= range okexResp.Data{
				modelOkex, ok := model.OkexMap[o.Tick]
				if !ok {
					logger.Error("can not find ticke:", o.Tick)
					continue
				}
				swapOkexMarket(&modelOkex, &o)
				//fmt.Println("----------------->modelHuoBi:",modelHuoBi)
				modelOkex.Update()
			}
		}
	}
}