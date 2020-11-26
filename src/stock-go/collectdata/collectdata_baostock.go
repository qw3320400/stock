package collectdata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"stock-go/common"
	"stock-go/data/mysql/stock"
	"stock-go/thirdparty/baostock"
	"stock-go/utils"
	"strings"
	"time"
)

const (
	DataSourceBaostock = "baostock"
)

func CollectBaostockData(request *CollectDataRequest) error {
	if request == nil || request.startTime.Unix() <= 0 || request.endTime.Unix() <= 0 {
		return utils.Errorf(nil, "request param error %+v", request)
	}
	// all stock code
	err := loadAllStockCode()
	if err != nil {
		return utils.Errorf(err, "loadAllStockCode fail")
	}
	// trade date
	err = loadStockTradeDates()
	if err != nil {
		return utils.Errorf(err, "loadStockTradeDates fail")
	}
	// k data
	err = loadStockKData(request)
	if err != nil {
		return utils.Errorf(err, "loadStockKData fail")
	}
	return nil
}

func loadAllStockCode() error {
	tmpTime := time.Now().In(time.FixedZone("CST", 8*60*60))
	for {
		bc, err := baostock.GetBaostockConnection()
		if err != nil {
			return utils.Errorf(err, "baostock.GetBaostockConnection fail")
		}
		tmpDate := tmpTime.Format("2006-01-02")
		allStockResponse, err := bc.QueryAllStock(tmpDate)
		if err != nil || allStockResponse == nil || allStockResponse.Rows == nil {
			return utils.Errorf(err, "bc.QueryAllStock fail")
		}
		if len(allStockResponse.Rows.Recode) <= 0 {
			tmpTime = tmpTime.AddDate(0, 0, -1)
			continue
		} else {
			// response to mysql
			codeList, err := allStockCodeResponseToData(allStockResponse)
			if err != nil {
				return utils.Errorf(err, "allStockCodeResponseToData fail")
			}
			exsitCode, err := stock.GetAllStockCode()
			if err != nil {
				return utils.Errorf(err, "mysql.GetAllStockCode fail")
			}
			exsitCodeMap := map[string]bool{}
			for _, code := range exsitCode.StockCodeList {
				exsitCodeMap[code.Code] = true
			}
			insertCodeList := []*common.StockCode{}
			for _, code := range codeList {
				if exsitCodeMap[code.Code] {
					continue
				}
				insertCodeList = append(insertCodeList, code)
			}
			if len(insertCodeList) > 0 {
				err = stock.InsertStockCode(&stock.InsertStockCodeRequest{
					StockCodeList: insertCodeList,
				})
				if err != nil {
					return utils.Errorf(err, "stock.InsertStockCode fail")
				}
			}
			break
		}
	}
	return nil
}

func allStockCodeResponseToData(response *baostock.QueryAllStockResponse) ([]*common.StockCode, error) {
	result := []*common.StockCode{}
	fieldMap := map[string]int{}
	for idx, field := range response.Fields {
		field = strings.Replace(field, " ", "", -1)
		fieldMap[field] = idx
	}
	for _, recode := range response.Rows.Recode {
		if fieldMap["code"] < 0 || fieldMap["code"] >= len(recode) ||
			fieldMap["code_name"] < 0 || fieldMap["code_name"] >= len(recode) {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", recode)
		}
		tmp := &common.StockCode{
			Code: recode[fieldMap["code"]],
			Name: recode[fieldMap["code_name"]],
		}
		if tmp.Code == "" {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", recode)
		}
		result = append(result, tmp)
	}
	return result, nil
}

func loadStockTradeDates() error {
	tmpTime := time.Now().In(time.FixedZone("CST", 8*60*60))
	for {
		bc, err := baostock.GetBaostockConnection()
		if err != nil {
			return utils.Errorf(err, "baostock.GetBaostockConnection fail")
		}
		tmpDate := tmpTime.Format("2006-01-02")
		tradeDatesResponse, err := bc.QueryTradeDates(EarliestDate, tmpDate)
		if err != nil || tradeDatesResponse == nil || tradeDatesResponse.Rows == nil {
			return utils.Errorf(err, "bc.QueryTradeDates fail")
		}
		if len(tradeDatesResponse.Rows.Recode) <= 0 {
			tmpTime = tmpTime.AddDate(0, 0, -1)
			continue
		} else {
			// response to mysql
			tradeDateList, err := stockTradeDatesResponseToData(tradeDatesResponse)
			if err != nil {
				return utils.Errorf(err, "stockTradeDatesResponseToData fail")
			}
			exsitTradeDate, err := stock.GetAllStockTradeDate()
			if err != nil {
				return utils.Errorf(err, "stock.GetAllStockCode fail")
			}
			exsitTradeDateMap := map[string]bool{}
			for _, tradeDate := range exsitTradeDate.StockTradeDateList {
				exsitTradeDateMap[tradeDate.DateCST.Format("2006-01-02")] = true
			}
			insertTradeDateList := []*common.StockTradeDate{}
			for _, tradeDate := range tradeDateList {
				if exsitTradeDateMap[tradeDate.DateCST.Format("2006-01-02")] {
					continue
				}
				insertTradeDateList = append(insertTradeDateList, tradeDate)
			}
			if len(insertTradeDateList) > 0 {
				err = stock.InsertStockTradeDate(&stock.InsertStockTradeDateRequest{
					StockTradeDateList: insertTradeDateList,
				})
				if err != nil {
					return utils.Errorf(err, "stock.InsertStockTradeDate fail")
				}
			}
			break
		}
	}
	return nil
}

func stockTradeDatesResponseToData(response *baostock.QueryTradeDatesResponse) ([]*common.StockTradeDate, error) {
	result := []*common.StockTradeDate{}
	fieldMap := map[string]int{}
	for idx, field := range response.Fields {
		field = strings.Replace(field, " ", "", -1)
		fieldMap[field] = idx
	}
	var (
		err error
	)
	for _, recode := range response.Rows.Recode {
		if fieldMap["calendar_date"] < 0 || fieldMap["calendar_date"] >= len(recode) ||
			fieldMap["is_trading_day"] < 0 || fieldMap["is_trading_day"] >= len(recode) {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", recode)
		}
		tmp := &common.StockTradeDate{}
		tmp.DateCST, err = baostockDateToData(recode[fieldMap["calendar_date"]])
		if err != nil {
			return nil, utils.Errorf(err, "返回数据错误 %+v", recode)
		}
		tmp.IsTradingDay, err = baostockIsTradingDateToData(recode[fieldMap["is_trading_day"]])
		if err != nil {
			return nil, utils.Errorf(err, "返回数据错误 %+v", recode)
		}
		result = append(result, tmp)
	}
	return result, nil
}

func loadStockKData(request *CollectDataRequest) error {
	tmpTime := request.endTime
	for {
		// check if break
		if request.startTime.After(tmpTime) {
			break
		}
		if request.isAllDataCode {
			exsitCode, err := stock.GetAllStockCode()
			if err != nil {
				return utils.Errorf(err, "stock.GetAllStockCode fail")
			}
			request.dataCodeList = []string{}
			for _, code := range exsitCode.StockCodeList {
				request.dataCodeList = append(request.dataCodeList, code.Code)
			}
		}
		// each code
		for _, code := range request.dataCodeList {
			err := loadStockKDataByMonthAndCode(code, tmpTime)
			if err != nil {
				return utils.Errorf(err, "loadStockKDataByMonthAndCode fail")
			}
		}
		// sub 1 month
		tmpTime = tmpTime.AddDate(0, -1, 0)
	}
	return nil
}

var (
	frequencyMap = map[string]string{
		"d":  "date,code,open,high,low,close,preclose,volume,amount,adjustflag,turn,tradestatus,pctChg,peTTM,psTTM,pcfNcfTTM,pbMRQ,isST",
		"w":  "date,code,open,high,low,close,volume,amount,adjustflag,turn,pctChg",
		"m":  "date,code,open,high,low,close,volume,amount,adjustflag,turn,pctChg",
		"5":  "date,code,open,high,low,close,volume,amount,adjustflag",
		"15": "date,code,open,high,low,close,volume,amount,adjustflag",
		"30": "date,code,open,high,low,close,volume,amount,adjustflag",
		"60": "date,code,open,high,low,close,volume,amount,adjustflag",
	}
)

func loadStockKDataByMonthAndCode(code string, date time.Time) error {
	startTime := time.Date(date.Year(), date.Month(), int(1), int(0), int(0), int(0), int(0), time.UTC)
	endTime := time.Date(date.Year(), date.Month()+1, int(1), int(0), int(0), int(0), int(0), time.UTC).Add(time.Hour * -24)
	for frequency := range frequencyMap {
		// 后复权
		adjustFlag := "post"
		err := loadStockKDataAndSave(code, startTime, endTime, frequency, adjustFlag)
		if err != nil {
			return utils.Errorf(err, "loadStockKDataAndSave fail")
		}

		// 不复权
		adjustFlag = "no"
		err = loadStockKDataAndSave(code, startTime, endTime, frequency, adjustFlag)
		if err != nil {
			return utils.Errorf(err, "loadStockKDataAndSave fail")
		}

	}
	return nil
}

func loadStockKDataAndSave(code string, startTime, endTime time.Time, frequency string, adjustFlag string) error {
	baostockAdjustFlag, err := dataAdjustFlagToBaostock(adjustFlag)
	if err != nil {
		return utils.Errorf(err, "dataAdjustFlagToBaostock %+v", adjustFlag)
	}

	bc, err := baostock.GetBaostockConnection()
	if err != nil {
		return utils.Errorf(err, "baostock.GetBaostockConnection fail")
	}
	stockKDataResponse, err := bc.QueryHistoryKDataPlusWithTimeOut(code, frequencyMap[frequency], startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), frequency, baostockAdjustFlag, 60)
	if err != nil || stockKDataResponse == nil || stockKDataResponse.Rows == nil {
		if err != baostock.QueryTimeoutErr {
			return utils.Errorf(err, "bc.QueryHistoryKDataPlusWithTimeOut fail")
		}
	}
	if err == baostock.QueryTimeoutErr {
		// 读取超时不中断
		fileData := []byte("query data timeout error")
		// 记录error的
		err = os.MkdirAll(filepath.Join(DataPath, ErrorPath), os.ModePerm)
		if err != nil {
			return utils.Errorf(err, "os.MkdirAll fail")
		}
		stockFileName := fmt.Sprintf(StockFileName, code, startTime.Format("2006-01"), frequency, baostockAdjustFlag)
		err = ioutil.WriteFile(filepath.Join(DataPath, ErrorPath, stockFileName), fileData, os.ModePerm)
		if err != nil {
			return utils.Errorf(err, "ioutil.WriteFile fail")
		}
		// 连接已经bock 需要重连
		bc, err = baostock.ReconnectBaostock()
		if err != nil {
			return utils.Errorf(err, "baostock.ReconnectBaostock fail")
		}
	} else {
		dataList, err := stockKDataResponseToData(stockKDataResponse, frequency)
		if err != nil {
			return utils.Errorf(err, "stockKDataResponseToData fail")
		}
		if len(dataList) > 0 {
			err = stock.InsertStockKData(&stock.InsertStockKDataRequest{
				StockKDataList: dataList,
			})
			if err != nil {
				return utils.Errorf(err, "stock.InsertStockKData fail")
			}
		}
	}
	return nil
}

func stockKDataResponseToData(response *baostock.QueryHistoryKDataResponse, frequency string) ([]*common.StockKData, error) {
	result := []*common.StockKData{}
	fieldMap := map[string]int{}
	for idx, field := range response.Fields {
		field = strings.Replace(field, " ", "", -1)
		fieldMap[field] = idx
	}
	var (
		err error
	)
	for _, recode := range response.Rows.Recode {
		if fieldMap["date"] < 0 || fieldMap["date"] >= len(recode) ||
			fieldMap["code"] < 0 || fieldMap["code"] >= len(recode) {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", recode)
		}
		tmp := &common.StockKData{
			Code:      recode[fieldMap["code"]],
			Frequency: frequency,
		}
		if tmp.Code == "" {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", recode)
		}
		tmp.TimeCST, err = baostockDateToData(recode[fieldMap["date"]])
		if err != nil {
			return nil, utils.Errorf(err, "返回数据错误 fail %+v", recode)
		}
		if idx, ok := fieldMap["open"]; ok && idx < len(recode) {
			tmp.Open = recode[idx]
		}
		if idx, ok := fieldMap["high"]; ok && idx < len(recode) {
			tmp.High = recode[idx]
		}
		if idx, ok := fieldMap["low"]; ok && idx < len(recode) {
			tmp.Low = recode[idx]
		}
		if idx, ok := fieldMap["close"]; ok && idx < len(recode) {
			tmp.Close = recode[idx]
		}
		if idx, ok := fieldMap["preclose"]; ok && idx < len(recode) {
			tmp.Preclose = recode[idx]
		}
		if idx, ok := fieldMap["volume"]; ok && idx < len(recode) {
			tmp.Volume = recode[idx]
		}
		if idx, ok := fieldMap["amount"]; ok && idx < len(recode) {
			tmp.Amount = recode[idx]
		}
		if idx, ok := fieldMap["adjustflag"]; ok && idx < len(recode) {
			tmp.AdjustFlag, err = baostockAdjustFlagToData(recode[idx])
			if err != nil {
				return nil, utils.Errorf(err, "返回数据错误 %+v", recode)
			}
		} else {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", recode)
		}
		if idx, ok := fieldMap["turn"]; ok && idx < len(recode) {
			tmp.Turn = recode[idx]
		}
		if idx, ok := fieldMap["tradestatus"]; ok && idx < len(recode) {
			tmp.TradeStatus, err = baostockTradeStatusToData(recode[idx])
			if err != nil {
				return nil, utils.Errorf(err, "返回数据错误 %+v", recode)
			}
		}
		if idx, ok := fieldMap["pctChg"]; ok && idx < len(recode) {
			tmp.PctChg = recode[idx]
		}
		if idx, ok := fieldMap["peTTM"]; ok && idx < len(recode) {
			tmp.PeTTM = recode[idx]
		}
		if idx, ok := fieldMap["psTTM"]; ok && idx < len(recode) {
			tmp.PsTTM = recode[idx]
		}
		if idx, ok := fieldMap["pcfNcfTTM"]; ok && idx < len(recode) {
			tmp.PcfNcfTTM = recode[idx]
		}
		if idx, ok := fieldMap["pbMRQ"]; ok && idx < len(recode) {
			tmp.PbMRQ = recode[idx]
		}
		if idx, ok := fieldMap["isST"]; ok && idx < len(recode) {
			tmp.IsST, err = baostockIsSTToData(recode[idx])
			if err != nil {
				return nil, utils.Errorf(err, "返回数据错误 %+v", recode)
			}
		}
		result = append(result, tmp)
	}
	return result, nil
}

func baostockDateToData(date string) (time.Time, error) {
	tmp, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, utils.Errorf(err, "time.Parse fail")
	}
	return tmp, nil
}

func baostockAdjustFlagToData(adjustFlag string) (string, error) {
	switch adjustFlag {
	case "1":
		return "post", nil
	case "2":
		return "pre", nil
	case "3":
		return "no", nil
	default:
		return "", utils.Errorf(nil, "返回数据错误 %+v", adjustFlag)
	}
}

func dataAdjustFlagToBaostock(adjustFlag string) (string, error) {
	switch adjustFlag {
	case "post":
		return "1", nil
	case "pre":
		return "2", nil
	case "no":
		return "3", nil
	default:
		return "", utils.Errorf(nil, "返回数据错误 %+v", adjustFlag)
	}
}

func baostockTradeStatusToData(tradeStatus string) (string, error) {
	switch tradeStatus {
	case "1":
		return "normal", nil
	case "0":
		return "stop", nil
	default:
		return "", utils.Errorf(nil, "返回数据错误 %+v", tradeStatus)
	}
}

func baostockIsTradingDateToData(isTradingDate string) (bool, error) {
	if isTradingDate == "1" {
		return true, nil
	}
	return false, nil
}

func baostockIsSTToData(isST string) (bool, error) {
	if isST == "1" {
		return true, nil
	}
	return false, nil
}
