## hlp
Tool to pull stock prices from Tiingo and put them into a hledger compatible format.

# TIINGO token
Your Tiingo token should be in the environment variable `TIINGO_TOKEN`

# Use in getprices
hledger2 has a helper script getprices to pull prices, and uses pricehist by default, to substitute hlp use:

1. Replace the definition of `CMD' with:
```
CMD="hlp -m $SOURCE -t $PAIR -s $START -e $END | sed -e 's|[-/]$BASECUR||'"
```
2. Change price provider section to be:
```
# Select an appropriate price provider and pair format.
# Customise as needed.

if isCurrency "$COMM"; then
    SOURCE="currency"
    PAIR="$COMM/$BASECUR"

elif isCryptocurrency "$COMM"; then
    SOURCE="crypto"
    PAIR="$COMM$BASECUR"

# other cases
elif [[ "$COMM" =~ ^(OMG)$ ]]; then
    SOURCE="crypto"
    PAIR="$COMM$BASECUR"

# default
else
    SOURCE="security"
    PAIR="$COMM"

fi
```
Important note: hlp assumes the base currency is USD, so getprices will fail if that is not true. 
