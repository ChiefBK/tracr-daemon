package main

import (
	"bytes"
	"encoding/gob"
	"os"
	"fmt"
)

type FullOrderBook struct {
	CurrencyPair string
	Asks         []OrderBookEntry
	Bids         []OrderBookEntry
	Sequence     float64
}

type OrderBookEntry struct {
	Price  float64
	Amount float64
}

func main() {
	a := os.Getenv("POLONIEX_API_KEYYY")
	fmt.Println(a)
	fmt.Println(len(a))

	fmt.Println("done")
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
