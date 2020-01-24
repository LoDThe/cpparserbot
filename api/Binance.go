package api

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance"

	"github.com/lodthe/cpparserbot/models"
)

//Binance provides with binance.com API methods
type Binance struct {
	client *binance.Client
}

//Init initializes Binance client with API keys
func (b *Binance) Init(apiKey, secretKey string) {
	b.client = binance.NewClient(apiKey, secretKey)
}

//GetPrice returns Binance price for given pair
func (b *Binance) GetPrice(pair *models.Pair) (float64, error) {
	prices, err := b.client.NewListPricesService().Symbol(pair.ToBinanceFormat()).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
		return 0, err
	}
	return strconv.ParseFloat(prices[0].Price, 64)
}

//Kline holds information about Binance kline record
type Kline struct {
	Price     float64
	Timestamp int64
}

//GetKlines returns information about how pair price was changing during the day
func (b *Binance) GetKlines(pair *models.Pair) ([]Kline, error) {
	klines, err := b.client.
		NewKlinesService().Symbol(pair.ToBinanceFormat()).
		Interval("1h").
		StartTime(int64(1000) * (time.Now().Add(-time.Hour * 24).Unix())).
		Do(context.Background())
	if err != nil {
		log.Fatalln(err)
		return make([]Kline, 0), err
	}

	var result []Kline

	for _, i := range klines {
		price, err := strconv.ParseFloat(i.Close, 64)
		if err != nil {
			log.Fatalln(err)
			return make([]Kline, 0), err
		}

		result = append(result, Kline{price, i.CloseTime})
	}

	return result, nil
}
