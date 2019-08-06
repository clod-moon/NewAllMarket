package model

import (
	"time"
	//"github.com/wonderivan/logger"
)



type Bian struct {
	Id 			 int       `gorm:"primary_key;type:int(11);AUTO_INCREMENT`
	//Ticker     string    `json:"ticker" gorm:"type:varchar(30);not null"`
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

func GetBian()  {
	var list []Bian
	DBHd.Find(&list)

	for _,v := range list{
		ticker := DMap[v.DealBiId]+SMap[v.StandardBiId]
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
