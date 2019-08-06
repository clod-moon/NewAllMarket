package model

import (
	"time"
	"strings"
	"github.com/wonderivan/logger"
)

type Huobi struct {
	Id           int       `gorm:"primary_key;type:int(11);AUTO_INCREMENT`
	Ticker       string    `json:"ticker"`
	DealBiId     int       `gorm:"type:int;not null"`
	StandardBiId int       `gorm:"type:int;not null"`
	Amount       float64   `json:"amount" gorm:"type:float;not null"`
	Open         float64   `json:"open" gorm:"type:float;not null"`
	Close        float64   `json:"close" gorm:"type:float;not null"`
	Rose         float64   `json:"rose" gorm:"type:float;not null"`
	High         float64   `json:"high" gorm:"type:float;not null"`
	Low          float64   `json:"low" gorm:"type:float;not null"`
	Count        float64   `json:"count" gorm:"type:float;not null"`
	Vol          float64   `json:"vol" gorm:"type:float;not null"`
	CreateTime   time.Time `gorm:"type:datetime;not null;"`
	UpdateTime   time.Time `gorm:"type:datetime;not null;"`
}

func GetHuoBi() {
	var list []Huobi
	DBHd.Find(&list)

	for _, v := range list {
		ticker := DMap[v.DealBiId] + SMap[v.StandardBiId]
		ticker = strings.ToLower(ticker)
		HuobiMap[ticker] = v
	}
	//fmt.Println("--------》HuobiMap：",len(HuobiMap))
}

func (h *Huobi) Update() {
	//fmt.Printf("huobi:",*h)
	h.UpdateTime = time.Now()
	DBHd.Model(&h).Update(h)
}

func QueryHuobiMarket(tick string) (list []Huobi,b bool) {

	if tick == "all"{
		DBHd.Find(&list)

		for i, v := range list {
			if v.Open == 0.0{
				list = append(list[:i],list[i:]...)
				continue
			}
			ticker := DMap[v.DealBiId] + "-" + SMap[v.StandardBiId]
			v.Ticker = ticker
		}
	}else {


		ticks := strings.Split(tick,"-")
		if len(ticks) < 2{
			logger.Error("非法的ticker:",tick)
			return list,false
		}

		dealBiId := DealBiMap[ticks[0]]
		standardBiId := StandardBiMap[ticks[1]]

		DBHd.Find(&list,"where deal_bi_id = ? and standard_bi_id = ?",dealBiId,standardBiId)
	}

	return list,true
}
