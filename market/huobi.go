package market

import (
	"fmt"
	"NewAllMarket/model"
	"regexp"
	"github.com/wonderivan/logger"
	"encoding/json"
	"strconv"
	"golang.org/x/net/websocket"
	"github.com/bitly/go-simplejson"
	"time"
	"sync"
)

var (
	regx = regexp.MustCompile(`market.([a-zA-Z]+).detail`)
)

func swapHuobiMarket(h, tmph *model.Huobi) {
	h.Open = tmph.Open
	h.Close = tmph.Close
	h.Amount = tmph.Amount
	h.High = tmph.High
	h.Low = tmph.Low
	h.Count = tmph.Count
	h.Vol = tmph.Vol
	h.Rose, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (h.Open-h.Close)/h.Open*100+0.005), 64)
}


func HuoBiPing(ws *websocket.Conn,pong int64){
	send([]byte(fmt.Sprintf( `{"pong":%d}`,pong)),ws)
	//time.Now().UnixNano() / int64(time.Millisecond))
}

func GetHuobiMarket(wg *sync.WaitGroup) {
	defer wg.Done()
	srcMarket,ok:=model.SrcMarketMap["huobi"]
	if !ok{
		logger.Error("can not find huobi")
		return
	}

	ws, err := websocket.Dial(srcMarket.WsUrl, "", Origin)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ping := time.Now().UnixNano() / int64(time.Millisecond)
	send([]byte(fmt.Sprintf( `{"ping":%d}`,ping)),ws)

	for _, ticker := range HuobiTickerList {
		go func(ticker string){
			message := []byte(fmt.Sprintf(srcMarket.Template,ticker,ticker))
			send(message, ws)
		}(ticker)
	}
	var msg = make([]byte, 512000)

	for {
		m, err := ws.Read(msg)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		newmsg := msg[:m]

		unzipmsg, _ := ParseGzip(newmsg, true)

		//fmt.Printf("Receive[UNZIP]: [%d:%d] %s\n", m, len(unzipmsg), unzipmsg[:])

		resp, err := simplejson.NewJson(unzipmsg)
		if err != nil {
			//logger.Debug(err.Error())
			continue
		}

		if ping := resp.Get("ping").MustInt64(); ping > 0 {
			HuoBiPing(ws,ping)
			continue
		}

		if ch := resp.Get("ch").MustString(); ch != "" {
			tick := regx.FindStringSubmatch(ch)
			if len(tick) < 2{
				logger.Error("订阅返回tick错误:",ch)
				continue
			}
			modelHuoBi ,ok:= model.HuobiMap[tick[1]]
			if !ok{
				logger.Error("can not find ticke:",tick[1])
				continue
			}
			strTick,_ := resp.Get("tick").Encode()
			var tmpHuobi model.Huobi
			err = json.Unmarshal([]byte(strTick), &tmpHuobi)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			swapHuobiMarket(&modelHuoBi,&tmpHuobi)
			modelHuoBi.Update()
		}
	}
}