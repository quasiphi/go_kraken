package main

import (
	"log"
	"os"

	"github.com/quasiphi/go_kraken/futures"
)

func main() {
	api := futures.New(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET"))
	data, err := api.Get_OHLC("15m", "PF_XBTUSD", 0)
	if err != nil {
		log.Panicln(err)
		return
	}
	log.Println(data)
}
