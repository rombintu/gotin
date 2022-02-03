package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/rombintu/gotin/api"
)

func TestGet(t *testing.T) {
	api := api.NewAPI("test.db")
	api.Token = os.Getenv("TOKEN")
	stockByName, err := api.StockByName("SPCE")
	if err != nil {
		t.Fatal(err)
	}
	figiByName, err := api.FigiByName(stockByName.Payload.Instruments[0].Figi)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(stockByName.Payload.Instruments[0].Name)
	fmt.Println(figiByName.Payload.LastPrice)
	fmt.Println(stockByName.Payload.Instruments[0].Currency)
}
