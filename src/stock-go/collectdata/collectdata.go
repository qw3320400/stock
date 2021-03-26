package collectdata

import (
	"stock-go/utils"
	"strings"
	"time"
)

const (
	EarliestDate = "2006-01-01"

	DataPath = "/Users/k/Desktop/code/stock/data/baostock"

	StockFileName = "%s:%s:%s:%s.csv"

	ErrorPath = "error"
)

type CollectDataRequest struct {
	DataSource string `json:"data_source"`
	DataCode   string `json:"data_code"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Frequency  string `json:"frequency"`
	AdjustFlag string `json:"adjust_flag"`
	// internal
	startTime     time.Time `json:"-"`
	endTime       time.Time `json:"-"`
	isAllDataCode bool      `json:"-"`
	dataCodeList  []string  `json:"-"`
}

func CollectData(request *CollectDataRequest) error {
	if request == nil || request.DataSource == "" || request.DataCode == "" || request.StartDate == "" || request.EndDate == "" {
		return utils.Errorf(nil, "request param error %+v", request)
	}
	if strings.ToLower(request.DataCode) == "all" {
		request.isAllDataCode = true
	} else {
		request.dataCodeList = []string{}
		for _, code := range strings.Split(request.DataCode, ",") {
			request.dataCodeList = append(request.dataCodeList, code)
		}
	}
	var (
		err error
	)
	request.startTime, err = time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return utils.Errorf(nil, "request param error %+v", request)
	}
	request.endTime, err = time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return utils.Errorf(nil, "request param error %+v", request)
	}
	switch request.DataSource {
	case DataSourceBaostock:
		return CollectBaostockData(request)
	case DataSourceJQData:
		return CollectJQDatakPricePeriodData(request)
	default:
		return utils.Errorf(nil, "request param error %+v", request)
	}
}
