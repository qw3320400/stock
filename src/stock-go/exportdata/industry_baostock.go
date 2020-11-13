package exportdata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"stock-go/thirdparty/baostock"
	"stock-go/utils"
)

const (
	AllStockIndustryPath     = "allstock"
	AllStockIndustryFileName = "allstockindustry:%s.csv"
)

func ExportBaostockAllCodeIndustry() error {
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
	allStockFileName := fmt.Sprintf(AllStockIndustryFileName, AllStockDate)
	allStockFilePath := filepath.Join(DataPath, AllStockIndustryPath)
	allStockFile := filepath.Join(allStockFilePath, allStockFileName)
	err = os.MkdirAll(allStockFilePath, os.ModePerm)
	if err != nil {
		return utils.Errorf(err, "os.MkdirAll fail")
	}
	allStockResponse, err := bc.QueryAllStock(AllStockDate)
	if err != nil || allStockResponse.Rows == nil || allStockResponse.Rows.Recode == nil || len(allStockResponse.Fields) <= 0 {
		return utils.Errorf(err, "bc.QueryAllStock fail")
	}
	if allStockResponse.Fields[0] != "code" {
		return utils.Errorf(nil, "allStockData error %+v", allStockResponse.Fields)
	}
	allStockResponse.Fields = append(allStockResponse.Fields, "industry", "industryClassification")
	for i := 0; i < len(allStockResponse.Rows.Recode); i++ {
		if len(allStockResponse.Rows.Recode[i]) <= 0 {
			return utils.Errorf(nil, "allStockData error %+v", allStockResponse.Rows.Recode[i])
		}
		code := allStockResponse.Rows.Recode[i][0]
		stockIndustryRespose, err := bc.QueryStockIndustry(code, AllStockDate)
		if err != nil || stockIndustryRespose.Rows == nil {
			return utils.Errorf(err, "bc.QueryStockIndustry fail")
		}
		if len(stockIndustryRespose.Rows.Recode) <= 0 {
			allStockResponse.Rows.Recode[i] = append(allStockResponse.Rows.Recode[i], "", "")
			continue
		}
		if len(stockIndustryRespose.Rows.Recode[0]) != len(stockIndustryRespose.Fields) || len(stockIndustryRespose.Fields) <= 4 {
			return utils.Errorf(nil, "industry data error %+v %+v", stockIndustryRespose.Rows.Recode, stockIndustryRespose.Fields)
		}
		allStockResponse.Rows.Recode[i] = append(allStockResponse.Rows.Recode[i], stockIndustryRespose.Rows.Recode[0][3], stockIndustryRespose.Rows.Recode[0][4])
	}
	fileData, err := baostockResponseToFileByte(allStockResponse.Fields, allStockResponse.Rows.Recode)
	if err != nil {
		return utils.Errorf(err, "baostockResponseToFileByte fail")
	}
	err = ioutil.WriteFile(allStockFile, fileData, os.ModePerm)
	if err != nil {
		return utils.Errorf(err, "ioutil.WriteFile fail")
	}
	return nil
}
