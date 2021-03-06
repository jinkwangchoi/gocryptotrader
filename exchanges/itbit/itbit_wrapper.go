package itbit

import (
	"log"
	"strconv"
	"time"

	"github.com/thrasher-/gocryptotrader/exchanges"
	"github.com/thrasher-/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-/gocryptotrader/exchanges/stats"
	"github.com/thrasher-/gocryptotrader/exchanges/ticker"
)

func (i *ItBit) Start() {
	go i.Run()
}
func (i *ItBit) Run() {
	if i.Verbose {
		log.Printf("%s polling delay: %ds.\n", i.GetName(), i.RESTPollingDelay)
		log.Printf("%s %d currencies enabled: %s.\n", i.GetName(), len(i.EnabledPairs), i.EnabledPairs)
	}

	for i.Enabled {
		for _, x := range i.EnabledPairs {
			currency := x
			go func() {
				ticker, err := i.GetTickerPrice(currency)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("ItBit %s: Last %f High %f Low %f Volume %f\n", currency, ticker.Last, ticker.High, ticker.Low, ticker.Volume)
				stats.AddExchangeInfo(i.GetName(), currency[0:3], currency[3:], ticker.Last, ticker.Volume)
			}()
		}
		time.Sleep(time.Second * i.RESTPollingDelay)
	}
}

func (i *ItBit) GetTickerPrice(currency string) (ticker.TickerPrice, error) {
	tickerNew, err := ticker.GetTicker(i.GetName(), currency[0:3], currency[3:])
	if err == nil {
		return tickerNew, nil
	}

	var tickerPrice ticker.TickerPrice
	tick, err := i.GetTicker(currency)
	if err != nil {
		return tickerPrice, err
	}

	tickerPrice.Ask = tick.Ask
	tickerPrice.Bid = tick.Bid
	tickerPrice.FirstCurrency = currency[0:3]
	tickerPrice.SecondCurrency = currency[3:]
	tickerPrice.Last = tick.LastPrice
	tickerPrice.High = tick.High24h
	tickerPrice.Low = tick.Low24h
	tickerPrice.Volume = tick.Volume24h
	ticker.ProcessTicker(i.GetName(), tickerPrice.FirstCurrency, tickerPrice.SecondCurrency, tickerPrice)
	return tickerPrice, nil
}

func (i *ItBit) GetOrderbookEx(currency string) (orderbook.OrderbookBase, error) {
	ob, err := orderbook.GetOrderbook(i.GetName(), currency[0:3], currency[3:])
	if err == nil {
		return ob, nil
	}

	var orderBook orderbook.OrderbookBase
	orderbookNew, err := i.GetOrderbook(currency)
	if err != nil {
		return orderBook, err
	}

	for x, _ := range orderbookNew.Bids {
		data := orderbookNew.Bids[x]
		price, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			log.Println(err)
		}
		amount, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			log.Println(err)
		}
		orderBook.Bids = append(orderBook.Bids, orderbook.OrderbookItem{Amount: amount, Price: price})
	}

	for x, _ := range orderbookNew.Asks {
		data := orderbookNew.Asks[x]
		price, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			log.Println(err)
		}
		amount, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			log.Println(err)
		}
		orderBook.Asks = append(orderBook.Asks, orderbook.OrderbookItem{Amount: amount, Price: price})
	}
	orderBook.FirstCurrency = currency[0:3]
	orderBook.SecondCurrency = currency[3:]
	orderbook.ProcessOrderbook(i.GetName(), orderBook.FirstCurrency, orderBook.SecondCurrency, orderBook)
	return orderBook, nil
}

//TODO Get current holdings from ItBit
//GetExchangeAccountInfo : Retrieves balances for all enabled currencies for the ItBit exchange
func (e *ItBit) GetExchangeAccountInfo() (exchange.ExchangeAccountInfo, error) {
	var response exchange.ExchangeAccountInfo
	response.ExchangeName = e.GetName()
	return response, nil
}
