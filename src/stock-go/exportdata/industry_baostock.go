package exportdata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"stock-go/thirdparty/baostock"
)

const (
	AllStockIndustryPath     = "allstock"
	AllStockIndustryFileName = "allstockindustry:%s.csv"
)

func ExportBaostockAllCodeIndustry() error {
	// 连接
	bc, err := baostock.NewBaostockConnection()
	if err != nil {
		return fmt.Errorf("[ExportBaostockAllCodeIndustry] baostock.NewBaostockConnection fail\n\t%s", err)
	}
	defer func() {
		bc.CloseConnection()
	}()
	// 登陆
	err = bc.Login("", "", 0)
	if err != nil {
		return fmt.Errorf("[ExportBaostockAllCodeIndustry] bc.Login fail\n\t%s", err)
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
		return fmt.Errorf("[ExportBaostockAllCodeIndustry] os.MkdirAll fail\n\t%s", err)
	}
	allStockResponse, err := bc.QueryAllStock(AllStockDate)
	if err != nil || allStockResponse.Rows == nil || allStockResponse.Rows.Recode == nil || len(allStockResponse.Fields) <= 0 {
		return fmt.Errorf("[ExportBaostockAllCodeIndustry] bc.QueryAllStock fail\n\t%s", err)
	}
	if allStockResponse.Fields[0] != "code" {
		return fmt.Errorf("[ExportBaostockAllCodeIndustry] allStockData error %+v", allStockResponse.Fields)
	}
	allStockResponse.Fields = append(allStockResponse.Fields, "industry", "industryClassification")
	for i := 0; i < len(allStockResponse.Rows.Recode); i++ {
		if len(allStockResponse.Rows.Recode[i]) <= 0 {
			return fmt.Errorf("[ExportBaostockAllCodeIndustry] allStockData error %+v", allStockResponse.Rows.Recode[i])
		}
		code := allStockResponse.Rows.Recode[i][0]
		stockIndustryRespose, err := bc.QueryStockIndustry(code, AllStockDate)
		if err != nil || stockIndustryRespose.Rows == nil {
			return fmt.Errorf("[ExportBaostockAllCodeIndustry] bc.QueryStockIndustry fail\n\t%s", err)
		}
		if len(stockIndustryRespose.Rows.Recode) <= 0 {
			allStockResponse.Rows.Recode[i] = append(allStockResponse.Rows.Recode[i], "", "")
			continue
		}
		if len(stockIndustryRespose.Rows.Recode[0]) != len(stockIndustryRespose.Fields) || len(stockIndustryRespose.Fields) <= 4 {
			return fmt.Errorf("[ExportBaostockAllCodeIndustry] industry data error %+v %+v", stockIndustryRespose.Rows.Recode, stockIndustryRespose.Fields)
		}
		allStockResponse.Rows.Recode[i] = append(allStockResponse.Rows.Recode[i], stockIndustryRespose.Rows.Recode[0][3], stockIndustryRespose.Rows.Recode[0][4])
	}
	fileData, err := baostockResponseToFileByte(allStockResponse.Fields, allStockResponse.Rows.Recode)
	if err != nil {
		return fmt.Errorf("[ExportBaostockAllCodeIndustry] baostockResponseToFileByte fail\n\t%s", err)
	}
	err = ioutil.WriteFile(allStockFile, fileData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("[ExportBaostockAllCodeIndustry] ioutil.WriteFile fail\n\t%s", err)
	}
	return nil
}
