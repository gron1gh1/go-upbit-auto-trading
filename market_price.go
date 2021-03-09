package main

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// [{"ticket":"UNIQUE_TICKET"},{"type":"ticker", "codes":["KRW-EOS"]}]

func main() {
	u := url.URL{Scheme: "wss", Host: "api.upbit.com", Path: "websocket/v1"}
	log.Printf("connect %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	oriJson := `[{"ticket":"UNIQUE_TICKET"},{"type":"ticker", "codes":["KRW-PXL","KRW-BTC"]}]`

	var data []interface{}
	json.Unmarshal([]byte(oriJson), &data)

	log.Printf("해석 %v\n", data[1].(map[string]interface{})["type"].(string))

	errr := c.WriteJSON(data)
	if errr != nil {
		log.Fatal("err", errr)
	}
	var data2 map[string]interface{}
	for {

		err := c.ReadJSON(&data2)
		if err != nil {
			log.Println("read:", err)
			return
		}
		price := data2["trade_price"].(float64)
		if price >= 100 {
			log.Printf("recv: %v", int(data2["trade_price"].(float64)))
		} else {
			log.Printf("recv: %v", data2["trade_price"].(float64))
		}

	}

}
