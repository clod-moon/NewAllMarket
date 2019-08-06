package model

import (
	"time"
	"github.com/wonderivan/logger"
	"strings"
)

type Bian struct {
	Id           int       `gorm:"primary_key;type:int(11);AUTO_INCREMENT`
	Ticker       string    `json:"tick" `
	DealBiId     int       `gorm:"type:int;not null"`
	StandardBiId int       `gorm:"type:int;not null"`
	Amount       float64   `json:"amount" gorm:"type:float;not null"`
	Open         float64   `json:"open" gorm:"type:float;not null"`
	Close        float64   `json:"close" gorm:"type:float;not null"`
	Rose         float64   `json:"rose" gorm:"type:float;not null"`
	High         float64   `json:"high" gorm:"type:float;not null"`
	Count        float64   `json:"count" gorm:"type:float;not null"`
	Low          float64   `json:"low" gorm:"type:float;not null"`
	Vol          float64   `json:"vol" gorm:"type:float;not null"`
	CreateTime   time.Time `gorm:"type:datetime;not null;"`
	UpdateTime   time.Time `gorm:"type:datetime;not null;"`
}

//func GetBian()  {
//	var list []Bian
//	DBHd.Find(&list)
//
//	for _,v := range list{
//		ticker := DealBiMap[v.DealBiId]+StandardBiMap[v.StandardBiId]
//		BianMap[ticker] = v
//	}
//}

func GetBian() {
	var list []Bian
	DBHd.Find(&list)

	for _, v := range list {
		ticker := DMap[v.DealBiId] + SMap[v.StandardBiId]
		//ticker =strings.ToLower(ticker)
		//logger.Debug("bian tick:",ticker)
		BianMap[ticker] = v
	}
	//fmt.Println("--------》HuobiMap：",len(HuobiMap))
}

func (b *Bian) Update() {
	b.UpdateTime = time.Now()
	DBHd.Model(b).Update(b)
}


func QueryBianMarket(tick string) (list []Bian,b bool) {

	if tick == "all"{
		DBHd.Find(&list)

		i := 0
		for ;i<len(list) ;{
			if list[i].Open <= 0.0001 && list[i].Open > -0.0001{
				list = append(list[:i],list[i+1:]...)
				continue
			}
			ticker := DMap[list[i].DealBiId] + "-" + SMap[list[i].StandardBiId]
			list[i].Ticker = ticker
			i++
		}
	}else {


		ticks := strings.Split(tick,"-")
		if len(ticks) < 2{
			logger.Error("非法的ticker:",tick)
			return list,false
		}

		dealBiId := DealBiMap[ticks[0]]
		standardBiId := StandardBiMap[ticks[1]]

		DBHd.Find(&list," deal_bi_id = ? and standard_bi_id = ?",dealBiId,standardBiId)

		if len(list) > 0{
			list[0].Ticker = tick
		}else{
			return list,false
		}
	}

	return list,true
}