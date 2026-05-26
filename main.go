package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	tiingo "github.com/jazzboME/tiin-go"
)

var markets = []string{ "c", "y", "s" }

func main() {
	var today = time.Now().Truncate(24 * time.Hour)
	var isoFormat = "2006-01-02"

	// command line flags
	var market = flag.StringP("market", "m", "s", "name of market to query, s (security), c (currency) or y (crypto)")
	var ticker = flag.StringP("ticker", "t", "", "name of ticker to pull values for")
	var startDate = flag.TimeP("start", "s", today, []string{isoFormat}, "first date to pull prices for, in YYYY-MM-DD")
	var endDate = flag.TimeP("end", "e", today, []string{isoFormat}, "last day to pull prices for, in YYYY-MM-DD")
	flag.Parse()
	
	if *ticker == "" {
		fail("No ticker supplied", nil)
	}

	if !slices.Contains(markets, *market) {
		fail(fmt.Sprintf("Market must be in %v", markets), nil)
	}
 	// internal variables
	//var pr []tiingo.EodPrice
	ctx := context.Background()
	// get Tiingo Token
	token := os.Getenv("TIINGO_TOKEN")

	if token == "" {
		fail("No token found, set TIINGO_TOKEN", nil)
	}

	c := tiingo.NewClient(token)

	switch *market {
	case "s":
		//var pr []tiingo.EodPrice
		columns := []string{"date", "adjClose", "divCash", "splitFactor"}
		var e = tiingo.EodPriceParams{
			StartDate: *startDate,
			EndDate: *endDate,
			ResampleFreq: tiingo.Daily,
			Sort: tiingo.DateAsc,
			RespFormat: tiingo.JSON,
			Columns: columns,
		}
		fmt.Println(e)
		resp, err := c.EodPrice(ctx, *ticker, &e)
		if err != nil {
			fail("failed to get security prices", err)
		}
		for _, price := range resp {
		fmt.Println(price)
		var cmt string
		if price.DivCash > 0 { cmt += fmt.Sprintf(" Ex-dividend of %f.", price.DivCash) }
		if price.SplitFactor != 1.0 { cmt += fmt.Sprintf(" Split at %.2f ratio.", price.SplitFactor) }
		if len(cmt) > 0 {
			cmt = " ;" + cmt
		}
		fmt.Printf("P %s %s %.8f USD%s\n", price.Date.Format(isoFormat), *ticker, price.AdjClose, cmt)
		}
	case "y":
		var e = tiingo.CryptoPriceParams{
			Exchanges: []string{"POLONIEX","GDAX"},
			StartDate: *startDate,
			EndDate: *endDate,
			ResampleFreq: tiingo.OneDay,
		}
		resp, err := c.CryptoPrice(ctx, []string{*ticker}, &e)
		if err != nil {
			fail("failed to get crypto prices", err)
		}
		for _, ticker := range resp {
			for _, price := range ticker.PriceData {
				fmt.Printf("P %s %s %.8f %s\n", price.Date.Format(isoFormat), strings.ToUpper(ticker.BaseCurrency), 
							price.Close, strings.ToUpper(ticker.QuoteCurrency))

			}
		}
	}
		//err = json.Unmarshal(resp, &pr)
		//if err != nil {
		//	fail("failed to unmarshal securities, err")
		//}
	
	/*
	columns := []string{"date", "adjClose", "divCash", "splitFactor"}
	resp, err := c.EodPrice(ctx, *ticker, *startDate, *endDate, tiingo.Daily, tiingo.DateAsc, tiingo.JSON, columns)
	if err != nil {
		fail("failed to get prices", err)
	}
	
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		fail("failed to unmarshal results", err)
	}
	*/
/*
	for _, price := range resp {
		fmt.Println(price)
		var cmt string
		if price.DivCash > 0 { cmt += fmt.Sprintf(" Ex-dividend of %f.", price.DivCash) }
		if price.SplitFactor != 1.0 { cmt += fmt.Sprintf(" Split at %.2f ratio.", price.SplitFactor) }
		if len(cmt) > 0 {
			cmt = " ;" + cmt
		}
		fmt.Printf("P %s %s %.8f USD%s\n", price.Date.Format(isoFormat), *ticker, price.AdjClose, cmt)
	}
*/
}

func fail(reason string, err error) {
	if err != nil { reason = reason + fmt.Sprintf(": %v", err) }
	fmt.Fprintf(os.Stderr, "%s\n", reason)
	os.Exit(1)
}