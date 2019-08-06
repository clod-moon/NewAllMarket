package main

import (
	"NewAllMarket/model"
	"NewAllMarket/market"
	"sync"
	"NewAllMarket/httpserver"
)

var(
	Wg sync.WaitGroup
)

func init() {

	model.Init()

	market.Init()


}

func main() {

	Wg.Add(1)
	go httpserver.GetMarketHandler(&Wg)

	Wg.Add(1)
	go market.GetHuobiMarket(&Wg)

	//wg.Add(1)
	//go market.GetOkexMarket()

	Wg.Add(1)
	go market.GetBianMarket(&Wg)

	Wg.Wait()
	return
}
