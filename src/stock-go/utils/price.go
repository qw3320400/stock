package utils

import (
	"strconv"
)

func GetDeltaPriceString(basePrice, price string) (string, error) {
	if basePrice == "" || price == "" {
		return "", Errorf(nil, "basePrice or price is nil")
	}
	basePriceF, err := strconv.ParseFloat(basePrice, 64)
	if err != nil {
		return "", Errorf(err, "strconv.ParseFloat fail")
	}
	priceF, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", Errorf(err, "strconv.ParseFloat fail")
	}
	return strconv.FormatFloat(priceF/basePriceF-1, 'f', -1, 64), nil
}
