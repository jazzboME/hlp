## hlp
Tool to pull stock prices from Tiingo and put them into a hledger compatible format.

# TIINGO token
Your Tiingo token should be in the environment variable `TIINGO_TOKEN`

# Use in getprices
hledger2 has a helper script getprices to pull prices, and uses pricehist by default, to substitute hlp use:

Replace the definition of `CMD' with:
```
CMD="hlp -t $PAIR -s $START -e $END | sed -e 's|[-/]$BASECUR||'"
```

Important note: hlp assumes the base currency is USD, so getprices will fail if that is not true. 
