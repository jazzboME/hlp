package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	tiingo "github.com/the-trader-dev/tiin-go"
)

func main() {
	var pr []tiingo.EodPrice
	token := os.Getenv("YOUR_TIINGO_TOKEN")
	ctx := context.Background()

	c := tiingo.NewClient(token)
	security := "AAPL"

	resp, err := c.DefaultEodPrice(ctx, security)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		log.Panic(err)
	}

	for _, price := range pr {
		fmt.Printf("P %s %s $%.2f\n", price.Date.Format("2006-01-02"), security, price.Close)
	}

}

