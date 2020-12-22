package jqdata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"stock-go/utils"
	"strings"
)

type GetCurrenctTokenResponse struct {
	Token string
}

func getCurrenctToken() (*GetCurrenctTokenResponse, error) {
	request := map[string]string{
		"method": JQDataMethodGetCurrentToken,
		"mob":    JQDataMobile,
		"pwd":    JQDataPassword,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, utils.Errorf(err, "json.Marshal fail")
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", JQDataAPIURL, strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, utils.Errorf(err, "http.NewRequest fail")
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK || resp.Body == nil {
		return nil, utils.Errorf(err, "client.Do fail %+v", resp)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.Errorf(err, "ioutil.ReadAll fail")
	}
	return &GetCurrenctTokenResponse{
		Token: string(responseBody),
	}, nil
}

type GetPricePeriodResponse struct {
	Fields []string
	Rows   *GetPricePeriodResponseRows
}

type GetPricePeriodResponseRows struct {
	Recode [][]string `json:"record"`
}

func GetPricePeriod(request map[string]string) (*GetPricePeriodResponse, error) {
	if request == nil || request["code"] == "" || request["unit"] == "" || request["date"] == "" || request["end_date"] == "" {
		return nil, utils.Errorf(nil, "param error %+v", request)
	}
	tokenResp, err := getCurrenctToken()
	if err != nil {
		return nil, utils.Errorf(err, "getCurrenctToken fail")
	}
	request["method"] = JQDataMethodGetPricePeriod
	request["token"] = tokenResp.Token
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, utils.Errorf(err, "json.Marshal fail")
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", JQDataAPIURL, strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, utils.Errorf(err, "http.NewRequest fail")
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK || resp.Body == nil {
		return nil, utils.Errorf(err, "client.Do fail %+v", resp)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.Errorf(err, "ioutil.ReadAll fail")
	}
	response, err := pricePeriodToResponse(string(responseBody))
	if err != nil {
		return nil, utils.Errorf(err, "pricePeriodToResponse fail")
	}
	return response, nil
}

func pricePeriodToResponse(csvData string) (*GetPricePeriodResponse, error) {
	ret := &GetPricePeriodResponse{
		Fields: []string{},
		Rows: &GetPricePeriodResponseRows{
			Recode: [][]string{},
		},
	}
	lines := strings.Split(csvData, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}
		words := strings.Split(line, ",")
		if i == 0 {
			for j := 0; j < len(words); j++ {
				word := words[j]
				if word == "" {
					return nil, utils.Errorf(nil, "返回数据错误 %+v", csvData)
				}
				ret.Fields = append(ret.Fields, word)
			}
			if len(ret.Fields) <= 0 {
				return nil, utils.Errorf(nil, "返回数据错误 %+v", csvData)
			}
			continue
		}
		if len(words) != len(ret.Fields) {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", csvData)
		}
		newRecode := []string{}
		for j := 0; j < len(words); j++ {
			word := words[j]
			newRecode = append(newRecode, word)
		}
		ret.Rows.Recode = append(ret.Rows.Recode, newRecode)
	}
	return ret, nil
}
