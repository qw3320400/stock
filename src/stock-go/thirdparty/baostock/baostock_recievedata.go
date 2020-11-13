package baostock

import (
	"encoding/json"
	"stock-go/utils"
	"strconv"
	"strings"
)

type RecieveData struct {
	MessageType       string
	MessageBodyLength int64
	ErrorCode         string
	ErrorMessage      string
	DataList          []string
	Response          interface{}
}

func (rd *RecieveData) GetResponse() error {
	if rd == nil {
		return utils.Errorf(nil, "rd is nil")
	}
	if rd.ErrorCode != BSERR_SUCCESS {
		return nil
	}
	var err error
	switch rd.MessageType {
	case MESSAGE_TYPE_LOGIN_RESPONSE:
		rd.Response, err = rd.SetLoginResponse()
	case MESSAGE_TYPE_GETKDATA_RESPONSE:
		rd.Response, err = rd.SetQueryHistoryKDataResponse()
	case MESSAGE_TYPE_QUERYTRADEDATES_RESPONSE:
		rd.Response, err = rd.SetQueryTradeDatesResponse()
	case MESSAGE_TYPE_QUERYALLSTOCK_RESPONSE:
		rd.Response, err = rd.SetQueryAllStockResponse()
	case MESSAGE_TYPE_QUERYSTOCKINDUSTRY_RESPONSE:
		rd.Response, err = rd.SetQueryStockIndustryResponse()
	}
	if err != nil {
		return utils.Errorf(err, "set %s data process fail", rd.MessageType)
	}
	return nil
}

type LoginResponse struct {
	Method string
	User   string
}

func (rd *RecieveData) SetLoginResponse() (*LoginResponse, error) {
	if rd == nil || rd.DataList == nil {
		return nil, utils.Errorf(nil, "rd or rd.DataList is nil")
	}
	if rd.DataList == nil || len(rd.DataList) < 2 {
		return nil, utils.Errorf(nil, "rd.DataList error %+v", rd.DataList)
	}
	data := &LoginResponse{
		Method: rd.DataList[0],
		User:   rd.DataList[1],
	}
	return data, nil
}

func (rd *RecieveData) GetLoginResponse() (*LoginResponse, error) {
	if rd == nil || rd.Response == nil {
		return nil, utils.Errorf(nil, "rd or rd.Response is nil")
	}
	data, ok := rd.Response.(*LoginResponse)
	if !ok {
		return nil, utils.Errorf(nil, "Response is not LoginData")
	}
	return data, nil
}

type QueryHistoryKDataResponse struct {
	Method       string
	User         string
	CurPageNum   int64
	PerPageCount int64
	Rows         *QueryHistoryKDataResponseRows
	Code         string
	Fields       []string
	StartDate    string
	EndDate      string
	Frequency    string
	AdjustFlag   string
}

type QueryHistoryKDataResponseRows struct {
	Recode [][]string `json:"record"`
}

func (rd *RecieveData) SetQueryHistoryKDataResponse() (*QueryHistoryKDataResponse, error) {
	if rd == nil || rd.DataList == nil {
		return nil, utils.Errorf(nil, "rd or rd.DataList is nil")
	}
	if rd.DataList == nil || len(rd.DataList) < 11 {
		return nil, utils.Errorf(nil, "rd.DataList error %+v", rd.DataList)
	}
	data := &QueryHistoryKDataResponse{
		Method:     rd.DataList[0],
		User:       rd.DataList[1],
		Code:       rd.DataList[5],
		StartDate:  rd.DataList[7],
		EndDate:    rd.DataList[8],
		Frequency:  rd.DataList[9],
		AdjustFlag: rd.DataList[10],
	}
	var err error
	data.CurPageNum, err = strconv.ParseInt(rd.DataList[2], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fail")
	}
	data.PerPageCount, err = strconv.ParseInt(rd.DataList[3], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fail")
	}
	data.Rows = &QueryHistoryKDataResponseRows{}
	if len(rd.DataList[4]) > 0 {
		err = json.Unmarshal([]byte(rd.DataList[4]), data.Rows)
		if err != nil {
			return nil, utils.Errorf(err, "json.Unmarshal fail")
		}
	}
	rd.DataList[6] = strings.TrimSpace(rd.DataList[6])
	data.Fields = strings.Split(rd.DataList[6], ",")
	return data, nil
}

func (rd *RecieveData) GetQueryHistoryKDataResponse() (*QueryHistoryKDataResponse, error) {
	if rd == nil || rd.Response == nil {
		return nil, utils.Errorf(nil, "rd or rd.Response is nil")
	}
	data, ok := rd.Response.(*QueryHistoryKDataResponse)
	if !ok {
		return nil, utils.Errorf(nil, "Response is not LoginData")
	}
	return data, nil
}

type QueryTradeDatesResponse struct {
	Method       string
	User         string
	CurPageNum   int64
	PerPageCount int64
	Rows         *QueryTradeDatesResponseRows
	StartDate    string
	EndDate      string
	Fields       []string
}

type QueryTradeDatesResponseRows struct {
	Recode [][]string `json:"record"`
}

func (rd *RecieveData) SetQueryTradeDatesResponse() (*QueryTradeDatesResponse, error) {
	if rd == nil || rd.DataList == nil {
		return nil, utils.Errorf(nil, "rd or rd.DataList is nil")
	}
	if rd.DataList == nil || len(rd.DataList) < 8 {
		return nil, utils.Errorf(nil, "rd.DataList error %+v", rd.DataList)
	}
	data := &QueryTradeDatesResponse{
		Method:    rd.DataList[0],
		User:      rd.DataList[1],
		StartDate: rd.DataList[5],
		EndDate:   rd.DataList[6],
	}
	var err error
	data.CurPageNum, err = strconv.ParseInt(rd.DataList[2], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fail")
	}
	data.PerPageCount, err = strconv.ParseInt(rd.DataList[3], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fail")
	}
	data.Rows = &QueryTradeDatesResponseRows{}
	if len(rd.DataList[4]) > 0 {
		err = json.Unmarshal([]byte(rd.DataList[4]), data.Rows)
		if err != nil {
			return nil, utils.Errorf(err, "json.Unmarshal fail")
		}
	}
	rd.DataList[7] = strings.TrimSpace(rd.DataList[7])
	data.Fields = strings.Split(rd.DataList[7], ",")
	return data, nil
}

func (rd *RecieveData) GetQueryTradeDatesResponse() (*QueryTradeDatesResponse, error) {
	if rd == nil || rd.Response == nil {
		return nil, utils.Errorf(nil, "rd or rd.Response is nil")
	}
	data, ok := rd.Response.(*QueryTradeDatesResponse)
	if !ok {
		return nil, utils.Errorf(nil, "Response is not LoginData")
	}
	return data, nil
}

type QueryAllStockResponse struct {
	Method       string
	User         string
	CurPageNum   int64
	PerPageCount int64
	Rows         *QueryAllStockResponseRows
	Date         string
	Fields       []string
}

type QueryAllStockResponseRows struct {
	Recode [][]string `json:"record"`
}

func (rd *RecieveData) SetQueryAllStockResponse() (*QueryAllStockResponse, error) {
	if rd == nil || rd.DataList == nil {
		return nil, utils.Errorf(nil, "rd or rd.DataList is nil")
	}
	if rd.DataList == nil || len(rd.DataList) < 7 {
		return nil, utils.Errorf(nil, "rd.DataList error %+v", rd.DataList)
	}
	data := &QueryAllStockResponse{
		Method: rd.DataList[0],
		User:   rd.DataList[1],
		Date:   rd.DataList[5],
	}
	var err error
	data.CurPageNum, err = strconv.ParseInt(rd.DataList[2], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fail")
	}
	data.PerPageCount, err = strconv.ParseInt(rd.DataList[3], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fail")
	}
	data.Rows = &QueryAllStockResponseRows{}
	if len(rd.DataList[4]) > 0 {
		err = json.Unmarshal([]byte(rd.DataList[4]), data.Rows)
		if err != nil {
			return nil, utils.Errorf(err, "json.Unmarshal fail")
		}
	}
	rd.DataList[6] = strings.TrimSpace(rd.DataList[6])
	data.Fields = strings.Split(rd.DataList[6], ",")
	return data, nil
}

func (rd *RecieveData) GetQueryAllStockResponse() (*QueryAllStockResponse, error) {
	if rd == nil || rd.Response == nil {
		return nil, utils.Errorf(nil, "rd or rd.Response is nil")
	}
	data, ok := rd.Response.(*QueryAllStockResponse)
	if !ok {
		return nil, utils.Errorf(nil, "Response is not LoginData")
	}
	return data, nil
}

type QueryStockIndustryResponse struct {
	Method       string
	User         string
	CurPageNum   int64
	PerPageCount int64
	Rows         *QueryStockIndustryResponseRows
	Code         string
	Date         string
	Fields       []string
}

type QueryStockIndustryResponseRows struct {
	Recode [][]string `json:"record"`
}

func (rd *RecieveData) SetQueryStockIndustryResponse() (*QueryStockIndustryResponse, error) {
	if rd == nil || rd.DataList == nil {
		return nil, utils.Errorf(nil, "rd or rd.DataList is nil")
	}
	if rd.DataList == nil || len(rd.DataList) < 8 {
		return nil, utils.Errorf(nil, "rd.DataList error %+v", rd.DataList)
	}
	data := &QueryStockIndustryResponse{
		Method: rd.DataList[0],
		User:   rd.DataList[1],
		Code:   rd.DataList[5],
		Date:   rd.DataList[6],
	}
	var err error
	data.CurPageNum, err = strconv.ParseInt(rd.DataList[2], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fail")
	}
	data.PerPageCount, err = strconv.ParseInt(rd.DataList[3], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fai")
	}
	data.Rows = &QueryStockIndustryResponseRows{}
	if len(rd.DataList[4]) > 0 {
		err = json.Unmarshal([]byte(rd.DataList[4]), data.Rows)
		if err != nil {
			return nil, utils.Errorf(err, "json.Unmarshal fail")
		}
	}
	rd.DataList[7] = strings.TrimSpace(rd.DataList[7])
	data.Fields = strings.Split(rd.DataList[7], ",")
	return data, nil
}

func (rd *RecieveData) GetQueryStockIndustryResponse() (*QueryStockIndustryResponse, error) {
	if rd == nil || rd.Response == nil {
		return nil, utils.Errorf(nil, "rd or rd.Response is nil")
	}
	data, ok := rd.Response.(*QueryStockIndustryResponse)
	if !ok {
		return nil, utils.Errorf(nil, "Response is not QueryStockIndustryResponse")
	}
	return data, nil
}
