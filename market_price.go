package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// Upbit Struct
type Upbit struct {
	conn    *websocket.Conn
	coinArr []string
}

// Connect Upbit WebSocket
func (u *Upbit) Connect() error {
	makeURL := url.URL{Scheme: "wss", Host: "api.upbit.com", Path: "websocket/v1"}
	log.Printf("connect %s", makeURL.String())

	c, _, err := websocket.DefaultDialer.Dial(makeURL.String(), nil)
	u.conn = c
	return err
}

// CoinAppend - Add a coin to check the price
func (u *Upbit) CoinAppend(coinNames []string) {
	u.coinArr = append(u.coinArr, coinNames...)
}

// Request Upbit WebSocket Coin Data
func (u *Upbit) Request() error {
	newCoinArr := []string{}
	for _, coin := range u.coinArr {
		newCoinArr = append(newCoinArr, fmt.Sprintf(`"KRW-%s"`, coin))
	}
	jsonData := fmt.Sprintf(`[{"ticket":"UNIQUE_TICKET"},{"type":"ticker", "codes":[%s]}]`, strings.Join(newCoinArr, ", "))

	var data []interface{}
	json.Unmarshal([]byte(jsonData), &data)

	err := u.conn.WriteJSON(data)
	return err
}

// Recv from Upbit Coin Price
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

	upbit.CoinAppend([]string{"DOT", "BTC", "MBL"})

	err = upbit.Request()
	if err != nil {
		log.Fatal("err", err)
	}
	upbit.Recv()
}
