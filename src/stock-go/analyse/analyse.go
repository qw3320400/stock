package analyse

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"stock-go/exportdata"
	"stock-go/utils"
	"strings"
	"time"
)

const (
	ResPath           = "/Users/k/Desktop/code/stock/res"
	LineChartFileName = "linechart.html"
)

type StockData struct {
	Code  string
	Name  string
	Price []*StockDataPrice
}

type StockDataPrice struct {
	Date  string
	Price string
}

func RelativeIndustry() error {
	file := filepath.Join(exportdata.DataPath, exportdata.AllStockIndustryPath, fmt.Sprintf(exportdata.AllStockIndustryFileName, exportdata.AllStockDate))
	fileData, err := utils.ReadCommonCSVFile(file)
	if err != nil {
		return fmt.Errorf("[RelativeIndustry] utils.ReadCommonCSVFile fail\n\t%s", err)
	}
	var codeIdx, nameIdx, industryIdx int
	for i := 0; i < len(fileData.Column); i++ {
		if fileData.Column[i] == "code" {
			codeIdx = i
		}
		if fileData.Column[i] == "industry" {
			industryIdx = i
		}
		if fileData.Column[i] == "code_name" {
			nameIdx = i
		}
	}
	codeList := []StockData{}
	for i := 0; i < len(fileData.Data); i++ {
		if len(fileData.Data[i]) <= codeIdx || len(fileData.Data[i]) <= industryIdx {
			return fmt.Errorf("[RelativeIndustry] data error %s", fileData.Data[i])
		}
		if fileData.Data[i][industryIdx] == "食品饮料" {
			codeList = append(codeList, StockData{
				Code: fileData.Data[i][codeIdx],
				Name: fileData.Data[i][nameIdx],
			})
		}
	}
	dateStr := exportdata.AllStockDate
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("[RelativeIndustry] time.Parse fail\n\t%s", err)
	}
	// 读取数据
	for i := 0; i < len(codeList); i++ {
		codeFile := filepath.Join(exportdata.DataPath,
			fmt.Sprintf(exportdata.StockPath, date.Year(), date.Month()),
			fmt.Sprintf(exportdata.StockFileName, codeList[i].Code, date.Format("2006-01"), "5", "1"))
		codeData, err := utils.ReadCommonCSVFile(codeFile)
		if err != nil {
			if err == utils.ErrFileNotExist {
				continue
			}
			return fmt.Errorf("[RelativeIndustry] utils.ReadCommonCSVFile fail\n\t%s", err)
		}
		var openIdx, dateIdx int
		for j := 0; j < len(codeData.Column); j++ {
			if codeData.Column[j] == "open" {
				openIdx = j
			}
			if codeData.Column[j] == "date" {
				dateIdx = j
			}
		}
		codeList[i].Price = []*StockDataPrice{}
		for j := 0; j < len(codeData.Data); j++ {
			codeList[i].Price = append(codeList[i].Price, &StockDataPrice{
				Date:  codeData.Data[j][dateIdx],
				Price: codeData.Data[j][openIdx],
			})
		}
	}
	// 基数
	basePriceMap := map[string]string{}
	for i := 0; i < len(codeList); i++ {
		if len(codeList[i].Price) > 0 {
			basePriceMap[codeList[i].Code] = codeList[i].Price[0].Price
		}
	}
	// 日期处理
	dateCodeMap := map[string]map[string]*StockData{}
	for i := 0; i < len(codeList); i++ {
		tmpCode := codeList[i].Code
		for j := 0; j < len(codeList[i].Price); j++ {
			tmpDate := codeList[i].Price[j].Date
			if dateCodeMap[tmpDate] == nil {
				dateCodeMap[tmpDate] = map[string]*StockData{}
			}
			if dateCodeMap[tmpDate][tmpCode] == nil {
				dateCodeMap[tmpDate][tmpCode] = &StockData{
					Code:  tmpCode,
					Name:  codeList[i].Name,
					Price: []*StockDataPrice{},
				}
			}
			dateCodeMap[tmpDate][tmpCode].Price = append(dateCodeMap[tmpDate][tmpCode].Price, codeList[i].Price[j])
		}
	}
	// 生成数据
	chartDataStr := "[\n\t['日期',"
	for i := 0; i < len(codeList); i++ {
		chartDataStr += fmt.Sprintf("'%s',", codeList[i].Name)
	}
	chartDataStr = chartDataStr[:len(chartDataStr)-1] + "],"
	startTime := time.Date(date.Year(), date.Month(), int(1), int(0), int(0), int(0), int(0), time.UTC)
	endTime := time.Date(date.Year(), date.Month()+1, int(1), int(0), int(0), int(0), int(0), time.UTC).Add(time.Hour * 24)
	for {
		if startTime.After(endTime) {
			break
		}
		dateStr := startTime.Format("2006-01-02")
		codeMap := dateCodeMap[dateStr]
		if codeMap != nil {
			for i := 0; i < 48; i++ {
				priceList := []string{}
				for j := 0; j < len(codeList); j++ {
					dateCode := codeMap[codeList[j].Code]
					if dateCode == nil || len(dateCode.Price) <= 0 {
						priceList = append(priceList, "0")
					} else {
						if len(dateCode.Price) != 48 || basePriceMap[codeList[j].Code] == "" {
							return fmt.Errorf("[RelativeIndustry] data error %+v", dateCode)
						}
						detaPrice, err := utils.GetDeltaPriceString(basePriceMap[codeList[j].Code], dateCode.Price[i].Price)
						if err != nil {
							return fmt.Errorf("[RelativeIndustry] utils.GetDeltaPriceString fail\n\t%s", err)
						}
						priceList = append(priceList, detaPrice)
					}
				}
				chartDataStr += fmt.Sprintf("\n\t['%s',%s],", dateStr, strings.Join(priceList, ","))
			}
		}
		startTime = startTime.Add(time.Hour * 24)
	}
	chartDataStr = chartDataStr[:len(chartDataStr)-1] + "\n\t]"
	chartData, err := ioutil.ReadFile(filepath.Join(ResPath, LineChartFileName))
	if err != nil {
		return fmt.Errorf("[RelativeIndustry] ioutil.ReadFile fail\n\t%s", err)
	}
	chartData = bytes.Replace(chartData, []byte("{data}"), []byte(chartDataStr), -1)
	err = ioutil.WriteFile(filepath.Join(ResPath, "test.html"), chartData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("[RelativeIndustry] ioutil.WriteFile fail\n\t%s", err)
	}
	return nil
}
