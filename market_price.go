package main

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Upbit Struct
type Upbit struct {
	conn *websocket.Conn
}

// Connect - Upbit WebSocket Connect
func (u *Upbit) Connect() error {
	makeURL := url.URL{Scheme: "wss", Host: "api.upbit.com", Path: "websocket/v1"}
	log.Printf("connect %s", makeURL.String())

	c, _, err := websocket.DefaultDialer.Dial(makeURL.String(), nil)
	u.conn = c
	return err
}

// Request - Upbit WebSocket Coin Data Request
func (u *Upbit) Request() error {
	jsonData := `[{"ticket":"UNIQUE_TICKET"},{"type":"ticker", "codes":["KRW-PXL","KRW-BTC"]}]`
	var data []interface{}
	json.Unmarshal([]byte(jsonData), &data)

	err := u.conn.WriteJSON(data)
	return err
}

// Recv - Recv from Upbit Coin Price
func (u *Upbit) Recv() error {
	var recvJSON map[string]interface{}
	for {

		err := u.conn.ReadJSON(&recvJSON)
		if err != nil {
			log.Println("read:", err)
			return err
		}
		price := recvJSON["trade_price"].(float64)
		coinName := recvJSON["code"]
		if price >= 100 { // 100원 이상은 호가단위 변경으로 소수점 없음
			log.Printf("[%s] : %v", coinName, int(price))
		} else {
			log.Printf("[%s] : %v", coinName, price)
		}

	}
}

func main() {
	upbit := Upbit{}

	err := upbit.Connect()
	defer upbit.conn.Close()

	if err != nil {
		log.Fatal("dial:", err)
	}

	err = upbit.Request()
	if err != nil {
		log.Fatal("err", err)
	}

	upbit.Recv()
}
