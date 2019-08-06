package webclient

import (
	"MarketTranspond/config"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/net/websocket"
	"sync"
	"compress/flate"

)



//var url = "wss://www.hbdm.com/ws"











func GzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)

}

func GetOKExMarket(wg sync.WaitGroup) {
	ws, err := websocket.Dial(config.OKEX_WS_URL, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	var msg = make([]byte, 512000)

	for {
		m, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		newmsg := msg[:m]

		unzipmsg, _ := GzipDecode(newmsg)

		fmt.Printf("Receive[UNZIP]: [%d:%d] %s\n", m, len(unzipmsg), unzipmsg[:])

		if len(unzipmsg) > 21 {
			pingcmd := string(unzipmsg[2:6])
			if "ping" == pingcmd {
				pingtime := string(unzipmsg[8:21])
				pongstr := fmt.Sprintf("{\"pong\":%s}", pingtime)
				message := []byte(pongstr)

				send(message, ws)
			}
		}
	}

	ws.Close() //关闭连接
	wg.Done()
}







func GetBianMarket(wg sync.WaitGroup) {
	ws, err := websocket.Dial(config.BIAN_WS_URL, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	var msg = make([]byte, 512000)

	for {
		m, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		newmsg := msg[:m]

		if len(newmsg) > 21 {
			pingcmd := string(newmsg[2:6])
			if "ping" == pingcmd {
				pingtime := string(newmsg[8:21])
				pongstr := fmt.Sprintf("{\"pong\":%s}", pingtime)
				message := []byte(pongstr)

				send(message, ws)
			}
		}
	}

	ws.Close() //关闭连接
	wg.Done()
}
