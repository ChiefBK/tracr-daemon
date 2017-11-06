package goku_bot

import (
	"github.com/gorilla/websocket"
	"net/url"
	"log"
	"time"
	"errors"
	"fmt"
)

func websocketConnect(address string, retries int) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: address}

	var connection *websocket.Conn
	var err error
	retriesLeft := retries
	for retriesLeft > 0 {
		connection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)

		if err != nil {
			log.Printf("error connecting: %s", err)
			log.Println("Retrying connection after 5 seconds")
			retriesLeft--

			timer := time.NewTimer(time.Second * 5)
			<-timer.C
		} else {
			break
		}
	}

	if retriesLeft == 0 {
		log.Println()
		return connection, errors.New(fmt.Sprintf("Could not connect after %d attemps", retries))
	}

	return connection, nil
}
