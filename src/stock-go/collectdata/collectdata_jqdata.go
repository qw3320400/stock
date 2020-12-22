package collectdata

import (
	"stock-go/common"
	"stock-go/data/mysql/stock"
	"stock-go/thirdparty/jqdata"
	"stock-go/utils"
	"strconv"
	"strings"
	"time"
)

const (
	DataSourceJQData = "jqdata"
)

func loadJQDataPricePeriodAndSave(code string, startTime, endTime time.Time, frequency string, adjustFlag string) error {
	jqdataAdjustFlag, err := dataAdjustFlagToJQData(adjustFlag)
	if err != nil {
		return utils.Errorf(err, "dataAdjustFlagToJQData %+v", adjustFlag)
	}
	jqdataFrequency, err := dataFrequencyToJQData(frequency)
	if err != nil {
		return utils.Errorf(err, "dataFrequencyToJQData %+v", frequency)
	}

	request := map[string]string{}
	request["code"] = code
	request["unit"] = jqdataFrequency
	request["date"] = startTime.Format("2006-01-02 15:04:05")
	request["end_date"] = endTime.Format("2006-01-02 15:04:05")
	request["fq_ref_date"] = jqdataAdjustFlag
	response, err := jqdata.GetPricePeriod(request)
	if err != nil {
		return utils.Errorf(err, "jqdata.GetPricePeriod")
	}
	dataList, err := jqdataPricePeriodResponseToData(response, code, adjustFlag, frequency)
	if err != nil {
		return utils.Errorf(err, "jqdataPricePeriodResponseToData fail")
	}
	if len(dataList) > 0 {
		err = stock.InsertStockKData(&stock.InsertStockKDataRequest{
			StockKDataList: dataList,
		})
		if err != nil {
			return utils.Errorf(err, "stock.InsertStockKData fail")
		}
	}

	return nil
}

func dataAdjustFlagToJQData(adjustFlag string) (string, error) {
	switch adjustFlag {
	case "post":
		return "2006-01-01", nil
	case "pre":
		return time.Now().In(time.FixedZone("CST", 8*60*60)).Format("2006-01-02"), nil
	case "no":
		return "", nil
	default:
		return "", utils.Errorf(nil, "返回数据错误 %+v", adjustFlag)
	}
}

/*
d, w, m, 5, 15, 30, 60
1m, 5m, 15m, 30m, 60m, 120m, 1d, 1w, 1M
*/

func dataFrequencyToJQData(frequency string) (string, error) {
	_, err := strconv.ParseInt(frequency, 10, 64)
	if err != nil {
		switch frequency {
		case "d":
			return "1" + "d", nil
		case "w":
			return "1" + "w", nil
		case "m":
			return "1" + "M", nil
		default:
			return "", utils.Errorf(nil, "返回数据错误 %+v", frequency)
		}
	}
	return frequency + "m", nil
}

func jqDataDateToData(date string) (time.Time, error) {
	tmp, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		tmp, err = time.Parse("2006-01-02", date)
		if err != nil {
			return time.Time{}, utils.Errorf(err, "time.Parse fail")
		} else {
			return tmp, nil
		}
	}
	return tmp, nil
}

func jqdataPricePeriodResponseToData(response *jqdata.GetPricePeriodResponse, code, adjustFlag, frequency string) ([]*common.StockKData, error) {
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
		if fieldMap["date"] < 0 || fieldMap["date"] >= len(recode) {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", recode)
		}
		tmp := &common.StockKData{
			Code:       code,
			AdjustFlag: adjustFlag,
			Frequency:  frequency,
		}
		if tmp.Code == "" {
			return nil, utils.Errorf(nil, "返回数据错误 %+v", recode)
		}
		tmp.TimeCST, err = jqDataDateToData(recode[fieldMap["date"]])
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
		if idx, ok := fieldMap["volume"]; ok && idx < len(recode) {
			tmp.Volume = recode[idx]
		}
		if idx, ok := fieldMap["money"]; ok && idx < len(recode) {
			tmp.Amount = recode[idx]
		}
		result = append(result, tmp)
	}
	return result, nil
}
