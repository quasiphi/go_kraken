package futures

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func (api *Kraken_Futures) Get_Balance() (Balances_Response, error) {
	// * Gets balance from the kraken futures account

	response_body, err := api.Request(
		"GET",
		"https://futures.kraken.com/derivatives/api/v3/accounts",
		"/api/v3/accounts",
		"",
		true,
	)
	if err != nil {
		return Balances_Response{}, err
	}

	balances_response := Balances_Response{}
	err = json.Unmarshal(response_body, &balances_response)
	if err != nil {
		log.Println("Reading body failed")
		return Balances_Response{}, err
	}

	return balances_response, nil
}

func (api *Kraken_Futures) Get_Leverage_Preferences(Tradeable string) (float32, error) {
	// * Gets the Leverage for the tradeable

	response_body, err := api.Request(
		"GET",
		"https://futures.kraken.com/derivatives/api/v3/leveragepreferences",
		"/api/v3/leveragepreferences",
		"",
		true,
	)
	if err != nil {
		return 0, err
	}

	leverage_preferences_response := Leverage_Response{}
	err = json.Unmarshal(response_body, &leverage_preferences_response)
	if err != nil {
		log.Println("Reading body failed")
		return 0, err
	}

	for i := 0; i < len(leverage_preferences_response.LeveragePreferences); i++ {
		symbol := leverage_preferences_response.LeveragePreferences[i].Symbol
		leverage := leverage_preferences_response.LeveragePreferences[i].MaxLeverge

		if symbol == Tradeable {
			// * Found our pair
			return leverage, nil
		}
	}

	log.Println(color.HiRedString("Could not find the specified pair"))

	return 1, nil
}

func (api *Kraken_Futures) Dead_Mans_Switch() {
	// * Implements a dead's man's switch for crashes
	post_body := "timeout=60"

	response_body, err := api.Request(
		"POST",
		"https://futures.kraken.com/derivatives/api/v3/cancelallordersafter",
		"/api/v3/cancelallordersafter",
		post_body,
		true,
	)
	if err != nil {
		log.Panic(color.HiRedString(fmt.Sprint(err)))
	}

	dead_mans_switch_response := Dead_Mans_Switch_Response{}
	err = json.Unmarshal(response_body, &dead_mans_switch_response)
	if err != nil {
		log.Println(color.HiRedString(fmt.Sprint("Could not parse DMS response: ", string(response_body), " Trace: ", err)))
	}

	if dead_mans_switch_response.Error != "" {
		log.Println(color.HiRedString(fmt.Sprint("Error in Dead Man's Switch: ", dead_mans_switch_response.Error)))
	}
}

func (api *Kraken_Futures) Deactivate_DMS() {
	// * Deactivates Dead Man's Switch
	post_body := "timeout=0"

	log.Println(color.RedString("Canceling Dead Man's Switch"))

	response_body, err := api.Request(
		"POST",
		"https://futures.kraken.com/derivatives/api/v3/cancelallordersafter",
		"/api/v3/cancelallordersafter",
		post_body,
		true,
	)
	if err != nil {
		log.Panic(color.HiRedString(fmt.Sprint(err)))
	}

	dead_mans_switch_response := Dead_Mans_Switch_Response{}
	err = json.Unmarshal(response_body, &dead_mans_switch_response)
	if err != nil {
		log.Panic(color.HiRedString(fmt.Sprint(err)))
	}

	if dead_mans_switch_response.Error != "" {
		log.Panic(color.HiRedString(fmt.Sprint("Error in Dead Man's Switch: ", dead_mans_switch_response.Error)))
	}
}

func (api *Kraken_Futures) Place_Order(OrderType, Side, Pair string, Limit, Quantity float64, Reduce bool) (UUID string, Price float64, Error error) {
	precision := Get_Precision(Pair)

	post_body := "orderType=" + OrderType + "&side=" + Side + "&size=" + strconv.FormatFloat(Quantity, 'f', precision, 64) + "&symbol=" + Pair + "&limitPrice=" + strconv.FormatFloat(Limit, 'f', 2, 64) + "&reduceOnly=" + strconv.FormatBool(Reduce)

	response_body, err := api.Request(
		"POST",
		"https://futures.kraken.com/derivatives/api/v3/sendorder",
		"/api/v3/sendorder",
		post_body,
		true,
	)
	if err != nil {
		return "", 0, err
	}

	send_order_response := Send_Order_Response{}
	err = json.Unmarshal(response_body, &send_order_response)
	if err != nil {
		log.Println("place_order, reading body failed")
		return "", 0, err
	}

	// TODO: Fix this, because it sometimes crashes here
	order_id := send_order_response.SendStatus.OrderEvents[0].OrderPriorExecution.OrderID

	if send_order_response.SendStatus.OrderEvents[0].Type == "EXECUTION" {
		price := send_order_response.SendStatus.OrderEvents[0].Price
		return order_id, price, nil
	} else if send_order_response.SendStatus.OrderEvents[0].Type == "REJECT" {
		log.Println("Order rejected: ", send_order_response.SendStatus.OrderEvents[0].Reason)
		if send_order_response.SendStatus.OrderEvents[0].Reason == "IOC_WOULD_NOT_EXECUTE" {
			return order_id, 0, errors.New("ioc_reject")
		}
		return order_id, 0, errors.New("order rejected")
	} else {
		log.Println("Order type:", send_order_response.SendStatus.OrderEvents[0].Type)
		return order_id, 0, errors.New("order did not execute")
	}
}

func (api *Kraken_Futures) Place_Instant_Order(Side, Pair string, Quantity float64, Reduce bool) (UUID string, Price float64, Error error) {
	var price_type string

	if Side == "buy" {
		price_type = "Ask"
	} else {
		price_type = "Bid"
	}

	current_price, err := api.Get_Current_Price(Pair, price_type)
	if err != nil {
		return "", 0, err
	}

	uuid, price, err := api.Place_Order("ioc", Side, Pair, current_price, Quantity, Reduce)
	if err != nil {
		if err.Error() == "ioc_reject" {
			log.Println("IOC was rejected, retrying...")
			time.Sleep(time.Second * 5)
			uuid, price, _ = api.Place_Instant_Order(Side, Pair, Quantity, Reduce)
			return uuid, price, nil
		}
		return "", 0, err
	}
	return uuid, price, nil
}
