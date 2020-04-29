package baostock

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type RecieveData struct {
	MessageType       string
	MessageBodyLength int64
	ErrorCode         string
	ErrorMessage      string
	DataList          []string
	StructData        interface{}
}

type LoginData struct {
	Method string
	User   string
}

func (rd *RecieveData) GetStructData() error {
	if rd == nil {
		return fmt.Errorf("[GetStructData] rd is nil")
	}
	if rd.ErrorCode != BSERR_SUCCESS {
		return nil
	}
	var err error
	switch rd.MessageType {
	case MESSAGE_TYPE_LOGIN_RESPONSE:
		rd.StructData, err = rd.SetLoginData()
	}
	if err != nil {
		return fmt.Errorf("[GetStructData] set %s data process fail\n\t%s", rd.MessageType, err)
	}
	return nil
}

func (rd *RecieveData) SetLoginData() (*LoginData, error) {
	if rd == nil || rd.DataList == nil {
		return nil, fmt.Errorf("[GetLoginData] rd or rd.DataList is nil")
	}
	if rd.DataList == nil || len(rd.DataList) < 2 {
		return nil, fmt.Errorf("[GetLoginData] rd.DataList error %+v", rd.DataList)
	}
	data := &LoginData{
		Method: rd.DataList[0],
		User:   rd.DataList[1],
	}
	return data, nil
}

func (rd *RecieveData) GetLoginData() (*LoginData, error) {
	if rd == nil || rd.StructData == nil {
		return nil, fmt.Errorf("[GetLoginData] rd or rd.StructData is nil")
	}
	data, ok := rd.StructData.(*LoginData)
	if !ok {
		return nil, fmt.Errorf("[GetLoginData] StructData is not LoginData")
	}
	return data, nil
}

type QueryHistoryKData struct {
	Method       string
	User         string
	CurPageNum   int64
	PerPageCount int64
	Rows         *QueryHistoryKDataRows
	Code         string
	Fields       []string
	StartDate    string
	EndDate      string
	Frequency    string
	AdjustFlag   string
}

type QueryHistoryKDataRows struct {
	Recode [][]string `json:"record"`
}

func (rd *RecieveData) SetQueryHistoryKData() (*QueryHistoryKData, error) {
	if rd == nil || rd.DataList == nil {
		return nil, fmt.Errorf("[SetQueryHistoryKData] rd or rd.DataList is nil")
	}
	if rd.DataList == nil || len(rd.DataList) < 11 {
		return nil, fmt.Errorf("[GetLoginData] rd.DataList error %+v", rd.DataList)
	}
	data := &QueryHistoryKData{
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
		return nil, fmt.Errorf("[GetLoginData] strconv.ParseInt fail %s", err)
	}
	data.PerPageCount, err = strconv.ParseInt(rd.DataList[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[GetLoginData] strconv.ParseInt fail %s", err)
	}
	data.Rows = &QueryHistoryKDataRows{}
	err = json.Unmarshal([]byte(rd.DataList[4]), data.Rows)
	if err != nil {
		return nil, fmt.Errorf("[GetLoginData] json.Unmarshal fail %s", err)
	}
	data.Fields = strings.Split(rd.DataList[6], ",")
	return data, nil
}

func (rd *RecieveData) GetQueryHistoryKData() (*QueryHistoryKData, error) {
	if rd == nil || rd.StructData == nil {
		return nil, fmt.Errorf("[SetQueryHistoryKData] rd or rd.StructData is nil")
	}
	data, ok := rd.StructData.(*QueryHistoryKData)
	if !ok {
		return nil, fmt.Errorf("[SetQueryHistoryKData] StructData is not LoginData")
	}
	return data, nil
}
