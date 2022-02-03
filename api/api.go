package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rombintu/gotinkoff/models"
)

type API struct {
	Timeout time.Duration
	Token   string
	Store   *Store
}

type Stock struct {
	Name     string
	Ticker   string
	Price    float32
	Currency string
}

func NewAPI(path string) *API {
	return &API{
		Timeout: time.Second * 3,
		Store:   NewStore(path),
	}
}

func NewStock(name, ticker, cur string, price float32) Stock {
	return Stock{
		Name:     name,
		Ticker:   ticker,
		Price:    price,
		Currency: cur,
	}
}

func (a *API) http_req(url string, method string) ([]byte, error) {

	client := &http.Client{
		Timeout: a.Timeout,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+a.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"register, bad response code '%s' from '%s'", resp.Status, url,
		)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil

}

func (a *API) Portfolio() (models.Portfolio, error) {
	var data models.Portfolio
	// var errPayload models.NotFoundError

	respBody, err := a.http_req("https://api-invest.tinkoff.ru/openapi/sandbox/portfolio", "GET")
	if err != nil {
		return models.Portfolio{}, err
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return models.Portfolio{}, err
	}
	// if err := json.Unmarshal(respBody, &errPayload); err != nil {
	// 	return models.Portfolio{}, err
	// }
	// log.Printf("%+v", errPayload)
	// if errPayload.PayloadError.Code != "200" || errPayload.PayloadError.Code == "" {
	// 	return models.Portfolio{}, fmt.Errorf(
	// 		"%s : %s", errPayload.PayloadError.Message, errPayload.PayloadError.Code,
	// 	)
	// }

	if strings.ToUpper(data.Status) != "OK" {
		return models.Portfolio{}, fmt.Errorf(
			"register failed, trackingId: '%s'", data.TrackingID,
		)
	}

	return data, nil
}

func (a *API) StocksAll() (models.Stocks, error) {
	var data models.Stocks
	// var errPayload models.NotFoundError

	respBody, err := a.http_req("https://api-invest.tinkoff.ru/openapi/sandbox/market/stocks", "GET")
	if err != nil {
		return models.Stocks{}, err
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return models.Stocks{}, err
	}
	// if err := json.Unmarshal(respBody, &errPayload); err != nil {
	// 	return models.Stocks{}, err
	// }

	// if errPayload.PayloadError.Code != "200" {
	// 	return models.Stocks{}, fmt.Errorf(
	// 		"%s : %s", errPayload.PayloadError.Message, errPayload.PayloadError.Code,
	// 	)
	// }

	if strings.ToUpper(data.Status) != "OK" {
		return models.Stocks{}, fmt.Errorf(
			"register failed, trackingId: '%s'", data.TrackingID,
		)
	}

	return data, nil

}

func (a *API) StockByName(name string) (models.Stocks, error) {
	var data models.Stocks
	// var errPayload models.NotFoundError

	respBody, err := a.http_req("https://api-invest.tinkoff.ru/openapi/sandbox/market/search/by-ticker?ticker="+name, "GET")
	if err != nil {
		return models.Stocks{}, err
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return models.Stocks{}, err
	}
	// if err := json.Unmarshal(respBody, &errPayload); err != nil {
	// 	return models.Stocks{}, err
	// }

	// if errPayload.PayloadError.Code != "200" {
	// 	return models.Stocks{}, fmt.Errorf(
	// 		"%s : %s", errPayload.PayloadError.Message, errPayload.PayloadError.Code,
	// 	)
	// }

	if strings.ToUpper(data.Status) != "OK" {
		return models.Stocks{}, fmt.Errorf(
			"register failed, trackingId: '%s'", data.TrackingID,
		)
	}

	return data, nil

}

func (a *API) FigiByName(figi string) (models.StocksByFigi, error) {
	var data models.StocksByFigi
	// var errPayload models.NotFoundError

	respBody, err := a.http_req(
		"https://api-invest.tinkoff.ru/openapi/sandbox/market/orderbook?figi="+figi+"&depth="+"1", "GET",
	)
	if err != nil {
		return models.StocksByFigi{}, err
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return models.StocksByFigi{}, err
	}
	// if err := json.Unmarshal(respBody, &errPayload); err != nil {
	// 	return models.StocksByFigi{}, err
	// }

	// if errPayload.PayloadError.Code != "200" {
	// 	return models.StocksByFigi{}, fmt.Errorf(
	// 		"%s : %s", errPayload.PayloadError.Message, errPayload.PayloadError.Code,
	// 	)
	// }

	if strings.ToUpper(data.Status) != "OK" {
		return models.StocksByFigi{}, fmt.Errorf(
			"register failed, trackingId: '%s'", data.TrackingID,
		)
	}

	return data, nil

}

func (a *API) GetStock(name string) (Stock, error) {
	stockByName, err := a.StockByName(name)
	if err != nil {
		return Stock{}, err
	}

	if len(stockByName.Payload.Instruments) == 0 {
		return Stock{}, fmt.Errorf("error. Скорее всего такой акции не существует")
	}
	figiByName, err := a.FigiByName(stockByName.Payload.Instruments[0].Figi)
	if err != nil {
		return Stock{}, err
	}

	return Stock{
		Name:     stockByName.Payload.Instruments[0].Name,
		Ticker:   stockByName.Payload.Instruments[0].Ticker,
		Currency: stockByName.Payload.Instruments[0].Currency,
		Price:    figiByName.Payload.LastPrice,
	}, nil
}
