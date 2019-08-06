package httpserver

import (
	"net/http"
	"encoding/json"
	"time"
	"io/ioutil"
	"NewAllMarket/model"
	"sync"
)

var messageChan = make(chan []byte, 1024)

type Requst struct {
	Platform string `json:"platform"`
	Tick     string `json:"tick"`
}

type Resp struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data []Market `json:"data"`
}

type Market struct {
	Name       string    `json:"name"`
	Close      float64   `json:"Close"`
	Open       float64   `json:"Open"`
	Rose       float64   `json:"Roes"`
	UpdateTime time.Time `json:"update_time"`
	CreateTime time.Time `json:"create_time"`
}

func MarketServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域

	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Header().Set("content-type", "application/json") //返回数据格式是json

	var request Requst
	var resp Resp
	requestBoyd, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(requestBoyd, &request)

	if request.Platform == "huobi" {

		list , b :=model.QueryHuobiMarket(request.Tick)
		if b == false{
			return
		}
		strList ,_ := json.Marshal(list)
		json.Unmarshal(strList,&resp.Data)

	} else if request.Platform == "bian" {

	}
	//fmt.Println(DB.Find(&Market{}).Value)


	resp.Code = 200

	resp.Msg = "success"

	ret, _ := json.Marshal(resp)

	w.Write(ret)
}

func GetMarketHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	http.HandleFunc("/get_market", MarketServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
