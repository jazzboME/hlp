package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	tiingo "github.com/the-trader-dev/tiin-go"
)

func main() {
	// command line flags
	var tf = flag.StringP("tickerfile", "f", "tickers.hlp", "file which contains tickers to fetch, one line per ticker")
	var debug = flag.BoolP("debug", "d", false, "writed debugging information to stderr")
	flag.Parse()

	// internal variables
	var pr []tiingo.EodPrice
	ctx := context.Background()
	// get Tiingo Token
	token := os.Getenv("YOUR_TIINGO_TOKEN")

	if token == "" {
		fmt.Fprintf(os.Stderr, "No token found. Set YOUR_TIINGO_TOKEN\n")
		os.Exit(1)
	}

	// read Tickers file
	tickers, err := readTickers(*tf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read tickerfile: %s\n", err)
		os.Exit(1)
	}

	if *debug {
		fmt.Fprintf(os.Stderr, "%s\n", tickers)
	}
	
	c := tiingo.NewClient(token)

	for _, security := range tickers {
		resp, err := c.DefaultEodPrice(ctx, security)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", security, err)
			continue
		}
		err = json.Unmarshal(resp, &pr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", security, err)
			continue
		}

		for _, price := range pr {
			fmt.Printf("P %s %s $%.2f\n", price.Date.Format("2006-01-02"), security, price.AdjClose)
		}
	}
}