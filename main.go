package main

import (
	"AllMarket/model"
	"AllMarket/market"
	"sync"
)

var(
	Wg sync.WaitGroup
)

func init() {

	model.Init()

	market.Init()


}

func main() {

	//wg.Add(1)
	//go market.GetHuobiMarket()

	//wg.Add(1)
	//go market.GetOkexMarket()

	Wg.Add(1)
	go market.GetBianMarket()

	Wg.Wait()
	return
}
