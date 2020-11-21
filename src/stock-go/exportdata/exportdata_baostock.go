package exportdata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"stock-go/thirdparty/baostock"
	"stock-go/utils"
	"strings"
	"time"
)

const (
	StartDate = "2015-04-01"
	EndData   = "2015-04-20"

	DataPath = "/Users/k/Desktop/code/stock/data/baostock"
	// DataPath = "/root/stock/data/baostock"

	AllStockDate     = "2020-04-17"
	AllStockPath     = "allstock"
	AllStockFileName = "allstock:%s.csv"

	TradeDatePath     = "tradedate"
	TradeDateFileName = "tradedate:%s.csv"

	StockPath     = "stock/%d/%d"
	StockFileName = "%s:%s:%s:%s.csv"

	ErrorPath = "error"
)

var (
	FrequencyList = [][]string{
		{"d", "date,code,open,high,low,close,preclose,volume,amount,adjustflag,turn,tradestatus,pctChg,peTTM,psTTM,pcfNcfTTM,pbMRQ,isST"},
		{"w", "date,code,open,high,low,close,volume,amount,adjustflag,turn,pctChg"},
		{"m", "date,code,open,high,low,close,volume,amount,adjustflag,turn,pctChg"},
		{"5", "date,code,open,high,low,close,volume,amount,adjustflag"},
		{"15", "date,code,open,high,low,close,volume,amount,adjustflag"},
		{"30", "date,code,open,high,low,close,volume,amount,adjustflag"},
		{"60", "date,code,open,high,low,close,volume,amount,adjustflag"},
	}
)

func ExportBaostockData() error {
	// 连接
	bc, err := baostock.NewBaostockConnection()
	if err != nil {
		return utils.Errorf(err, "baostock.NewBaostockConnection fail")
	}
	defer func() {
		bc.CloseConnection()
	}()
	// 登陆
	err = bc.Login("", "", 0)
	if err != nil {
		return utils.Errorf(err, "bc.Login fail")
	}
	defer func() {
		bc.Logout()
	}()

	// all stock code
	allStockFileName := fmt.Sprintf(AllStockFileName, AllStockDate)
	allStockFilePath := filepath.Join(DataPath, AllStockPath)
	err = os.MkdirAll(allStockFilePath, os.ModePerm)
	if err != nil {
		return utils.Errorf(err, "os.MkdirAll fail")
	}
	allStockResponse, err := bc.QueryAllStock(AllStockDate)
	if err != nil || allStockResponse.Rows == nil {
		return utils.Errorf(err, "bc.QueryAllStock fail")
	}
	utils.Log("writing file ... " + allStockFileName)
	fileData, err := baostockResponseToFileByte(allStockResponse.Fields, allStockResponse.Rows.Recode)
	if err != nil {
		return utils.Errorf(err, "baostockResponseToFileByte fail")
	}
	err = ioutil.WriteFile(filepath.Join(allStockFilePath, allStockFileName), fileData, os.ModePerm)
	if err != nil {
		return utils.Errorf(err, "ioutil.WriteFile fail")
	}
	utils.Log("write file success " + allStockFileName)

	// trade date
	startTime, err := time.Parse("2006-01-02", StartDate)
	if err != nil {
		return utils.Errorf(err, "time.Parse fail")
	}
	nowTime := time.Now()
	tradeDateFileName := fmt.Sprintf(TradeDateFileName, nowTime.Format("2006-01-02"))
	tradeDateFilePath := filepath.Join(DataPath, TradeDatePath)
	err = os.MkdirAll(tradeDateFilePath, os.ModePerm)
	if err != nil {
		return utils.Errorf(err, "os.MkdirAll fail")
	}
	tradeDateResponse, err := bc.QueryTradeDates(startTime.Format("2006-01-02"), nowTime.Format("2006-01-02"))
	if err != nil || tradeDateResponse.Rows == nil {
		return utils.Errorf(err, "bc.QueryTradeDates fail")
	}
	utils.Log("writing file ... " + tradeDateFileName)
	fileData, err = baostockResponseToFileByte(tradeDateResponse.Fields, tradeDateResponse.Rows.Recode)
	if err != nil {
		return utils.Errorf(err, "baostockResponseToFileByte fail")
	}
	err = ioutil.WriteFile(filepath.Join(tradeDateFilePath, tradeDateFileName), fileData, os.ModePerm)
	if err != nil {
		return utils.Errorf(err, "ioutil.WriteFile fail")
	}
	utils.Log("write file success " + tradeDateFileName)

	// k data
	startTime = time.Date(startTime.Year(), startTime.Month(), int(1), int(0), int(0), int(0), int(0), time.UTC)
	endTime, err := time.Parse("2006-01-02", EndData)
	if err != nil {
		return utils.Errorf(err, "time.Parse fail")
	}
	endTime = time.Date(endTime.Year(), endTime.Month(), int(1), int(0), int(0), int(0), int(0), time.UTC)
	for {
		// check if break
		if startTime.After(endTime) {
			break
		}
		// each code
		for i := 0; i < len(allStockResponse.Rows.Recode); i++ {
			err = ExportBaostockDataByMonth(bc, allStockResponse.Rows.Recode[i][0], endTime)
			if err != nil {
				return utils.Errorf(err, "ExportBaostockDataByMonth fail")
			}
		}
		// sub 1 month
		endTime = time.Date(endTime.Year(), endTime.Month()-1, int(1), int(0), int(0), int(0), int(0), time.UTC)
	}

	return nil
}

func ExportBaostockDataByMonth(bc *baostock.BaostockConnection, code string, date time.Time) error {
	if bc == nil || code == "" {
		return utils.Errorf(nil, "bc or code is nil")
	}
	startTime := time.Date(date.Year(), date.Month(), int(1), int(0), int(0), int(0), int(0), time.UTC)
	endTime := time.Date(date.Year(), date.Month()+1, int(1), int(0), int(0), int(0), int(0), time.UTC).Add(time.Hour * -24)
	stockPath := fmt.Sprintf(StockPath, startTime.Year(), startTime.Month())
	stockFilePath := filepath.Join(DataPath, stockPath)
	err := os.MkdirAll(stockFilePath, os.ModePerm)
	if err != nil {
		return utils.Errorf(err, "os.MkdirAll fail")
	}
	for _, frequency := range FrequencyList {
		// 后复权
		err = queryAndSaveBaostockKData(bc, code, startTime, endTime, stockFilePath, frequency, "1")
		if err != nil {
			return utils.Errorf(err, "queryAndSaveBaostockKData fail")
		}

		// 不复权
		err = queryAndSaveBaostockKData(bc, code, startTime, endTime, stockFilePath, frequency, "3")
		if err != nil {
			return utils.Errorf(err, "queryAndSaveBaostockKData fail")
		}
	}
	return nil
}

func queryAndSaveBaostockKData(bc *baostock.BaostockConnection, code string, startTime, endTime time.Time, stockFilePath string, frequency []string, adjustFlag string) error {
	stockFileName := fmt.Sprintf(StockFileName, code, startTime.Format("2006-01"), frequency[0], adjustFlag)
	stockFile := filepath.Join(stockFilePath, stockFileName)
	// 如果文件已存在则跳过
	if _, err := os.Stat(stockFile); os.IsNotExist(err) {
		stockResponse, err := bc.QueryHistoryKDataPlusWithTimeOut(code, frequency[1], startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), frequency[0], adjustFlag, 60)
		if err != nil || stockResponse.Rows == nil {
			if err != baostock.QueryTimeoutErr {
				return utils.Errorf(err, "bc.QueryTradeDates fail")
			}
		}
		utils.Log("writing file ... " + stockFileName)
		var fileData []byte
		if err == baostock.QueryTimeoutErr {
			// 读取超时不中断
			fileData = []byte("query data timeout error")
			// 记录error的
			err = os.MkdirAll(filepath.Join(DataPath, ErrorPath), os.ModePerm)
			if err != nil {
				return utils.Errorf(err, "os.MkdirAll fail")
			}
			err = ioutil.WriteFile(filepath.Join(DataPath, ErrorPath, stockFileName), fileData, os.ModePerm)
			if err != nil {
				return utils.Errorf(err, "ioutil.WriteFile fail")
			}
			// 连接已经bock 需要重连
			err = bc.ReConnect()
			if err != nil {
				return utils.Errorf(err, "bc.ReConnect fail")
			}
		} else {
			fileData, err = baostockResponseToFileByte(stockResponse.Fields, stockResponse.Rows.Recode)
			if err != nil {
				return utils.Errorf(err, "baostockResponseToFileByte fail")
			}
		}
		err = ioutil.WriteFile(stockFile, fileData, os.ModePerm)
		if err != nil {
			return utils.Errorf(err, "ioutil.WriteFile fail")
		}
		utils.Log("write file success " + stockFileName)
	}
	return nil
}

func baostockResponseToFileByte(fields []string, record [][]string) ([]byte, error) {
	if fields == nil || len(fields) <= 0 {
		return nil, utils.Errorf(nil, "fields is nil")
	}
	ret := ""
	ret += strings.Join(fields, ",")
	for i := 0; i < len(record); i++ {
		ret += "\n"
		ret += strings.Join(record[i], ",")
	}
	return []byte(ret), nil
}
