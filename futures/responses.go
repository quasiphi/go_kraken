package futures

type Kraken_Futures_OHLC_Response struct {
	Candles      []Candle `json:"candles"`
	More_Candles bool     `json:"more_candles"`
}

type Candle struct {
	Time   int     `json:"time"`
	Close  float64 `json:"close,string"`
	High   float64 `json:"high,string"`
	Low    float64 `json:"low,string"`
	Open   float64 `json:"open,string"`
	Volume float32 `json:"volume,string"`
}

type Balances_Response struct {
	Result string
	// ? Because this will or could change at any time due to specific accounts
	// ? I will leave it as interface and use response as blueprint rather than
	// ? dynamically assinging names to the accounts
	Accounts interface{}
}

type Leverage_Response struct {
	Result              string `json:"result"`
	ServerTime          string `json:"serverTime"`
	LeveragePreferences []struct {
		Symbol     string  `json:"symbol"`
		MaxLeverge float32 `json:"maxLeverage"`
	} `json:"leveragePreferences"`
}

type Send_Order_Response struct {
	Result     string `json:"result"`
	ServerTime string `json:"serverTime"`
	SendStatus struct {
		OrderID      string `json:"orderId"`
		Status       string `json:"status"`
		ReceivedTime string `json:"receivedTime"`
		OrderEvents  []struct {
			ExecutionID          string  `json:"executionId"`
			Price                float64 `json:"price"`
			Amount               float64 `json:"amount"`
			OrderPriorEdit       string  `json:"orderPriorEdit"`
			TakerReducedQuantity string  `json:"takerReducedQuantity"`
			Type                 string  `json:"type"`
			Reason               string  `json:"reason"`
			OrderPriorExecution  struct {
				OrderID             string  `json:"orderId"`
				Cli0rID             string  `json:"cliOrdId"`
				Type                string  `json:"type"`
				Symbol              string  `json:"symbol"`
				Side                string  `json:"side"`
				Quantity            float64 `json:"quantity"`
				Filled              float64 `json:"filled"`
				LimitPrice          float64 `json:"limitPrice"`
				ReduceOnly          bool    `json:"reduceOnly"`
				Timestamp           string  `json:"timestamp"`
				LastUpdateTimestamp string  `json:"lastUpdateTimestamp"`
			}
		}
	}
}

type Current_Price_Response struct {
	Result     string `json:"result"`
	ServerTime string `json:"serverTime"`
	Tickers    []struct {
		Tag                   string  `json:"Tag"`
		Pair                  string  `json:"pair"`
		Symbol                string  `json:"symbol"`
		MarkPrice             float64 `json:"markPrice"`
		Bid                   float64 `json:"bid"`
		BidSize               float64 `json:"bidSize"`
		Ask                   float64 `json:"ask"`
		AskSize               float64 `json:"askSize"`
		Vol24h                float64 `json:"vol24h"`
		OpenInterest          float64 `json:"openInterest"`
		Open24h               float64 `json:"open24h"`
		IndexPrice            float64 `json:"indexPrice"`
		Last                  float64 `json:"last"`
		LastTime              string  `json:"lastTime"`
		LastSize              float64 `json:"lastSize"`
		Suspended             bool    `json:"suspended"`
		FundingRate           float64 `json:"fundingRate"`
		FundingRatePrediction float64 `json:"fundingRatePrediction"`
		PostOnly              bool    `json:"postOnly"`
	}
}

type Dead_Mans_Switch_Response struct {
	Result string `json:"result"`
	Status struct {
		CurrentTime string `json:"currentTime"`
		TriggerTime string `json:"triggerTime"`
	}
	ServerTime string `json:"serverTime"`
	Error      string `json:"error"`
}
