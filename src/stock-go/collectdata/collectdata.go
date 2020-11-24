package collectdata

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"stock-go/data/mysql"
	"stock-go/thirdparty/baostock"
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
	default:
		return utils.Errorf(nil, "request param error %+v", request)
	}
}

func FileToMysql() error {
	// 链接数据库
	err := mysql.Connect()
	if err != nil {
		return utils.Errorf(err, "mysql.Connect fail")
	}
	defer mysql.Close()

	root := "/Users/k/Desktop/code/stock/data/baostock/stock"
	yearList, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}
	for _, year := range yearList {
		yearPath := filepath.Join(root, year.Name())
		if strings.HasPrefix(year.Name(), ".") {
			continue
		}
		if year.Name() == "2015" || year.Name() == "2016" || year.Name() == "2017" {
			continue
		}

		monthList, err := ioutil.ReadDir(yearPath)
		if err != nil {
			panic(err)
		}
		for _, month := range monthList {
			monthPath := filepath.Join(yearPath, month.Name())
			if strings.HasPrefix(month.Name(), ".") {
				continue
			}

			fileList, err := ioutil.ReadDir(monthPath)
			if err != nil {
				panic(err)
			}
			for _, file := range fileList {
				filePath := filepath.Join(monthPath, file.Name())
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}

				// read file
				fileName := file.Name()
				if !strings.HasSuffix(fileName, ".csv") {
					panic(err)
				}
				fileName = fileName[:len(fileName)-4]
				nameList := strings.Split(fileName, ":")
				if len(nameList) != 4 {
					panic(" len(nameList) != 4 ")
				}
				var (
					code       = nameList[0]
					dateStr    = nameList[1]
					frequency  = nameList[2]
					adjustFlag = nameList[3]
				)
				utils.Log(fmt.Sprintf("%+v %+v %+v %+v", code, dateStr, frequency, adjustFlag))

				fileData, err := utils.ReadCommonCSVFile(filePath)
				if err != nil {
					if err == utils.ErrFileQueryTimeout {
						continue
					}
					panic(err)
				}
				stockKDataResponse := &baostock.QueryHistoryKDataResponse{
					Fields: fileData.Column,
					Rows: &baostock.QueryHistoryKDataResponseRows{
						Recode: fileData.Data,
					},
				}
				dataList, err := stockKDataResponseToData(stockKDataResponse, frequency)
				if err != nil {
					panic(err)
				}
				if len(dataList) > 0 {
					err = mysql.InsertStockKData(&mysql.InsertStockKDataRequest{
						StockKDataList: dataList,
					})
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
	return nil
}
