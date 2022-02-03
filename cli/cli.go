package cli

import (
	"fmt"
	"log"
	"os"

	API "github.com/rombintu/gotinkoff/api"
	"github.com/urfave/cli"
)

func Cli() {
	app := cli.NewApp()

	app.Name = "Gotin"
	app.Usage = "Manage stocks of tinkoff"
	commands(app)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func commands(app *cli.App) {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "ticker",
			Value: "",
		},
		cli.StringFlag{
			Name:  "dbpath",
			Value: "./dev.db",
			// Required: true,
		},
		cli.StringFlag{
			Name: "token",
		},
		cli.BoolFlag{
			Name: "verbose",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "add",
			Usage:       "Добавить акцию в базу данных",
			Description: "Добавить акцию в базу данных для мониторинга",
			ShortName:   "a",
			Flags:       flags,
			Action: func(c *cli.Context) error {
				api := API.NewAPI(c.String("dbpath"))
				fmt.Printf("Была использована база данных: %s\n\n", c.String("dbpath"))

				token, err := api.Store.GetToken()
				if err != nil {
					return err
				}
				api.Token = token

				newStock, err := api.GetStock(c.String("ticker"))
				if err != nil {
					return err
				}

				if err := api.Store.AddStock(newStock); err != nil {
					return err
				}

				fmt.Printf("Была добавлена позиция: %s %f %s",
					newStock.Name,
					newStock.Price,
					newStock.Currency,
				)
				return nil
			},
		},
		{
			Name:  "init",
			Usage: "Создать базу данных",
			Flags: flags,
			Action: func(c *cli.Context) error {
				if c.String("token") == "" {
					fmt.Println("Token not set")
					return nil
				}

				api := API.NewAPI(c.String("dbpath"))
				if err := api.Store.Init(c.String("token")); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:  "get",
			Usage: "Получить айтем",
			Flags: flags,
			Action: func(c *cli.Context) error {
				api := API.NewAPI(c.String("dbpath"))
				token, err := api.Store.GetToken()
				if err != nil {
					return err
				}
				api.Token = token
				var stock API.Stock
				switch c.String("ticker") {
				case "rand", "":
					stock, err = api.Store.GetStockRand()
					if err != nil {
						return err
					}
				default:
					stock, err = api.Store.GetStockByTicker(c.String("ticker"))
					if err != nil {
						return err
					}

				}

				updatedStock, err := api.GetStock(stock.Ticker)
				if err != nil {
					return err
				}
				if c.Bool("verbose") {
					fmt.Printf("%s %.2f %s", updatedStock.Name, updatedStock.Price, updatedStock.Currency)
				} else {
					fmt.Printf("%s %.2f %s", updatedStock.Ticker, updatedStock.Price, updatedStock.Currency)
				}
				return nil
			},
		},
		{
			Name:      "delete",
			Usage:     "Удалить акцию из базы данных",
			ShortName: "d",
			Flags:     flags,
			Action: func(c *cli.Context) error {
				api := API.NewAPI(c.String("dbpath"))
				token, err := api.Store.GetToken()
				if err != nil {
					return err
				}
				api.Token = token
				ticker := c.String("ticker")
				if ticker == "" {
					fmt.Println("Введите имя тикера. --ticker NAME")
					return nil
				}
				if err := api.Store.DelStockByTicker(ticker); err != nil {
					return err
				}
				return nil
			},
		},
	}
}
