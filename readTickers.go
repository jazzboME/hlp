package main

import (
	"bufio"
	"os"
)

func readTickers(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tickers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tickers = append(tickers, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tickers, err
}