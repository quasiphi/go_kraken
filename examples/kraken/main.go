package main

import (
	"log"
	"os"

	kraken "github.com/quasiphi/go_kraken/kraken"
)

func main() {
	api := kraken.New(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET"))
	data, err := api.Ticker("XXBTZUSD")
	if err != nil {
		log.Panicln(err)
		return
	}
	log.Println(data)
}
