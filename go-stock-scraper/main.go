package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change string
}

func main() {
	ticker := []string{
		"MSFT",
		"IBM",
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"WMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAPL",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
		"VZ",
	}

	stocks := []Stock{} // Added empty brackets to initialize an empty slice

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting:", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("something went wrong:", err)
	})

	c.OnHTML("div#quote-header-info", func(e *colly.HTMLElement) {
		stock := Stock{}
		stock.company = e.ChildText("h1")
		fmt.Println("company", stock.company)
		stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
		fmt.Println("price:", stock.price)

		stocks = append(stocks, stock)
	})

	c.Wait()

	for _, t := range ticker {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}

	fmt.Println(stocks)

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	headers := []string{
		"company",
		"price",
		"change",
	}
	writer.Write(headers)

	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
		}
		writer.Write(record)
	}

	defer writer.Flush()
}
