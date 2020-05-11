package utils

import (
	"fmt"
	"strconv"
)

func GetDeltaPriceString(basePrice, price string) (string, error) {
	if basePrice == "" || price == "" {
		return "", fmt.Errorf("[GetDeltaPriceString] basePrice or price is nil")
	}
	basePriceF, err := strconv.ParseFloat(basePrice, 64)
	if err != nil {
		return "", fmt.Errorf("[GetDeltaPriceString] strconv.ParseFloat fail\n\t%s", err)
	}
	priceF, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", fmt.Errorf("[GetDeltaPriceString] strconv.ParseFloat fail\n\t%s", err)
	}
	return strconv.FormatFloat(priceF/basePriceF-1, 'f', -1, 64), nil
}
