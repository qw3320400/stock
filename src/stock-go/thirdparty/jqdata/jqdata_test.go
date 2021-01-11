package jqdata

import (
	"encoding/json"
	"testing"
)

func TestToken(t *testing.T) {
	t.Log(getCurrenctToken())
}

func TestFreqency(t *testing.T) {
	request := map[string]string{}
	request["code"] = "000300.XSHG"
	request["unit"] = "15m"
	request["date"] = "2020-01-01 00:00:00"
	request["end_date"] = "2020-02-01 00:00:00"
	request["fq_ref_date"] = ""
	response, err := GetPricePeriod(request)
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
