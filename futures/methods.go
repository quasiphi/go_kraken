package futures

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	// ? Taken from: https://futures.kraken.com/derivatives/api/v3/instruments
	PrecisionRates = map[string]int{
		"PF_XBTUSD":   4,
		"PF_ETHUSD":   3,
		"PF_SOLUSD":   2,
		"PF_UNIUSD":   1,
		"PF_ATOMUSD":  1,
		"PF_BCHUSD":   2,
		"PF_DOTUSD":   1,
		"PF_FILUSD":   1,
		"PF_LINKUSD":  1,
		"PF_LTCUSD":   2,
		"PF_DEFIUSD":  2,
		"PF_AVAXUSD":  2,
		"PF_XMRUSD":   2,
		"PF_AAVEUSD":  2,
		"PF_YFIUSD":   4,
		"PF_MKRUSD":   3,
		"PF_ETCUSD":   1,
		"PF_ZECUSD":   2,
		"PF_EGLDUSD":  2,
		"PF_OCEANUSD": 1,
		"PF_BANDUSD":  1,
		"PF_BALUSD":   1,
		"PF_ENSUSD":   1,
		"PF_QNTUSD":   2,
		"PF_APTUSD":   1,
		"PF_GMXUSD":   2,
		"PF_INJUSD":   1,
	}
)

func (api *Kraken_Futures) Get_OHLC(Interval, Tradeable string, FromTS int) (OHLC_Response Kraken_Futures_OHLC_Response, Err error) {
	// * Gets OHLC data from the Kraken Futures API

	params := url.Values{}
	if FromTS != 0 {
		params.Add("from", strconv.Itoa(FromTS))
	}

	URL := fmt.Sprintf("https://futures.kraken.com/api/charts/v1/spot/%s/%s?", Tradeable, Interval)

	response_body, err := api.Request(
		"GET",
		URL,
		"",
		"",
		false,
	)
	if err != nil {
		return Kraken_Futures_OHLC_Response{}, err
	}

	ohlc_response := Kraken_Futures_OHLC_Response{}
	err = json.Unmarshal(response_body, &ohlc_response)
	if err != nil {
		log.Println("Reading body failed")
		return Kraken_Futures_OHLC_Response{}, err
	}

	return ohlc_response, nil
}

func (api *Kraken_Futures) Get_Current_Price(Pair, Type string) (Price float64, Error error) {
	// * Gets the current prices of the pair

	response_body, err := api.Request(
		"GET",
		"https://futures.kraken.com/derivatives/api/v3/tickers",
		"",
		"",
		false,
	)
	if err != nil {
		return 0, err
	}

	current_prices_response := Current_Price_Response{}
	err = json.Unmarshal(response_body, &current_prices_response)
	if err != nil {
		log.Println("Get_current_price, reading body failed")
		return 0, err
	}

	for i := 0; i < len(current_prices_response.Tickers); i++ {
		ticker := current_prices_response.Tickers[i]

		if ticker.Symbol == strings.ToLower(Pair) {
			// * Print current market conditions
			log.Println(color.YellowString(fmt.Sprint("Mark Price: ", ticker.MarkPrice)))
			log.Println(color.YellowString(fmt.Sprint("Current Ask Size: ", ticker.AskSize)))
			log.Println(color.YellowString(fmt.Sprint("Current Bid Size: ", ticker.BidSize)))
			log.Println(color.YellowString(fmt.Sprint("Funding Rate: ", strconv.FormatFloat(ticker.FundingRate, 'f', -1, 64))))

			switch Type {
			case "Mark":
				return current_prices_response.Tickers[i].MarkPrice, nil
			case "Bid":
				return current_prices_response.Tickers[i].Bid, nil
			case "Ask":
				return current_prices_response.Tickers[i].Ask, nil
			case "Last":
				return current_prices_response.Tickers[i].Last, nil
			default:
				return current_prices_response.Tickers[i].Ask, nil
			}
		}
	}

	return 0, errors.New("could not find the current price of the specified pair")
}

func Get_Precision(Pair string) int {
	// * Returns a precision of the symbol
	if val, ok := PrecisionRates[Pair]; ok {
		return val
	} else {
		return 0
	}
}
