package main

import (
	"tracr-client"
	"tracr-daemon/exchanges"
	"time"
	"github.com/inconshreveable/log15"
)

func main() {
	key := "qvktu9Z2HklydKJrUQK2SkS5mNkOZA3tk7QS0iInlV3uV31LT7VJTytZ"
	secret := "pvm2iGZS2mxcgtX+yY45kNlZXY0s/yicBduK/JXMlbwQ3iBA44sp9wcI0Ji6Ydzdy3uobM16eiWZ8Y02OngwIg=="

	client := tracr_client.NewApiClient(key, secret, exchanges.KRAKEN, "https://api.kraken.com", "https://api.kraken.com", 1*time.Second)

	resp, err := client.Do("POST", "/0/private/Balance", nil, nil, nil)

	if err != nil {
		log15.Error("there was an error", "error", err)
		return
	}

	log15.Info("resp", "val", string(resp))
}
