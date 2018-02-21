package kraken

import "encoding/json"

type KrakenOrderBookResponse map[string]KrakenOrderBook

type KrakenOrderBook struct {
	Asks []KrakenOrderBookItem
	Bids []KrakenOrderBookItem
}

type KrakenOrderBookItem struct {
	Price     string
	Amount    string
	Timestamp int64
}

func (self *KrakenOrderBookItem) UnmarshalJSON(data []byte) error {

	// https://github.com/beldur/kraken-go-api-client/blob/master/types.go#L434
	tmpStruct := struct {
		price  string
		amount string
		ts     int64
	}{}
	tmpArr := []interface{}{&tmpStruct.price, &tmpStruct.amount, &tmpStruct.ts}
	err := json.Unmarshal(data, &tmpArr)

	if err != nil {
		return err
	}

	self.Price = *tmpArr[0].(*string)
	self.Amount = *tmpArr[1].(*string)
	self.Timestamp = *tmpArr[2].(*int64)

	if err != nil {
		return err
	}

	return nil
}
