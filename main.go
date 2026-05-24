package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	tiingo "github.com/the-trader-dev/tiin-go"
)

func main() {
	var today = time.Now().Truncate(24 * time.Hour)
	var isoFormat = "2006-01-02"

	// command line flags
	var ticker = flag.StringP("ticker", "t", "", "name of ticker to pull values for")
	var startDate = flag.TimeP("start", "s", today, []string{isoFormat}, "first date to pull prices for, in YYYY-MM-DD")
	var endDate = flag.TimeP("end", "e", today, []string{isoFormat}, "last day to pull prices for, in YYYY-MM-DD")
	flag.Parse()
	
	if *ticker == "" {
		fail("No ticker supplied", nil)
	}

 	// internal variables
	var pr []tiingo.EodPrice
	ctx := context.Background()
	// get Tiingo Token
	token := os.Getenv("TIINGO_TOKEN")

	if token == "" {
		fail("No token found, set TIINGO_TOKEN", nil)
	}

	c := tiingo.NewClient(token)
	columns := []string{"date", "adjClose", "divCash", "splitFactor"}
	resp, err := c.EodPrice(ctx, *ticker, *startDate, *endDate, tiingo.Daily, tiingo.DateAsc, tiingo.JSON, columns)
	if err != nil {
		fail("failed to get prices", err)
	}
	
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		fail("failed to unmarshal results", err)
	}

	for _, price := range pr {
		fmt.Printf("P %s %s %.8f USD\n", price.Date.Format(isoFormat), *ticker, price.AdjClose)
	}
}

func fail(reason string, err error) {
	if err != nil { reason = reason + ":"}
	fmt.Fprintf(os.Stderr, "%s %v", reason, err)
	os.Exit(1)
}