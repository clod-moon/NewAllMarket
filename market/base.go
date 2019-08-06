package market

import (
	"github.com/bitly/go-simplejson"
	"strings"
	"fmt"
	"NewAllMarket/model"
	"github.com/wonderivan/logger"
	"bytes"
	"encoding/binary"
	"compress/gzip"
	"io/ioutil"
	"golang.org/x/net/websocket"
	"compress/flate"
)

type JSON = simplejson.Json

var buffer bytes.Buffer

var ( //
	//wss://www.hbdm.com/ws
	HuobiEndpoint = "wss://api.huobi.pro/ws"

	OkexEndpoint = "wss://real.okex.com:10442/ws/v3"

	Origin          = "http://www.baidu.com"
	HuobiTickerList []string
	OkexTickerList  []string
)

func getAllTicker() {
	for value, _ := range model.StandardBiMap {
		i := 0
		for v, _ := range model.DealBiMap {
			if v != value {
				lowerValue := strings.ToLower(value)
				lowerV := strings.ToLower(v)
				ticker := fmt.Sprintf("%s%s", lowerV, lowerValue)
				HuobiTickerList = append(HuobiTickerList, ticker)
				if i < 45{
					ticker = fmt.Sprintf("%s-%s", v, value)
					OkexTickerList = append(OkexTickerList, ticker)
					i++
				}
			}
		}
	}
	logger.Debug("getAllTicker")
	return
}

func Init() {
	getAllTicker()
}

func errHandler(data []byte) []byte {
	buffer.Write(data)
	msg, err := ParseGzip(buffer.Bytes(), false)
	if err == nil {
		//fmt.Println("!!!!!!", string(msg[:]))
		return msg
	}
	return nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)

}

func ParseGzip(data []byte, handleErr bool) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		//with error
		if handleErr {
			msg := errHandler(data)
			if msg != nil{
				return msg,nil
			}
		}else{
			return nil, err
		}

	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			//with error
			if handleErr {
				errHandler(data)
			}
			return nil, err
		} else {
			//buffer.Reset()
			return undatas, nil
		}
	}
}

func send(message []byte, ws *websocket.Conn) {
	_, err := ws.Write(message)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Debug("Send: %s\n", message)
}