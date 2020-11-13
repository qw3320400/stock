package analyse

import (
	"fmt"
	"path/filepath"
	"stock-go/exportdata"
	"stock-go/utils"
	"strconv"
	"time"
)

type WeekDayData struct {
	Total int64
	Win   int64
}

func WeekDay() ([]*WeekDayData, error) {
	startTimeStr := "2017-01-01"
	endTimeStr := "2020-09-29"
	code := "sh.000300"
	startTime, err := time.Parse("2006-01-02", startTimeStr)
	if err != nil {
		return nil, utils.Errorf(err, "time.Parse fail")
	}
	endTime, err := time.Parse("2006-01-02", endTimeStr)
	if err != nil {
		return nil, utils.Errorf(err, "time.Parse fail")
	}
	result := make([]*WeekDayData, 5)
	for {
		if startTime.After(endTime) {
			break
		}

		file := filepath.Join(exportdata.DataPath, fmt.Sprintf(exportdata.StockPath, startTime.Year(), startTime.Month()), fmt.Sprintf(exportdata.StockFileName, code, startTime.Format("2006-01"), "d", "3"))
		fileData, err := utils.ReadCommonCSVFile(file)
		if err != nil {
			return nil, utils.Errorf(err, "utils.ReadCommonCSVFile fail")
		}
		columnIndexMap := map[string]int{}
		for i := 0; i < len(fileData.Column); i++ {
			columnIndexMap[fileData.Column[i]] = i
		}
		dateIdx := columnIndexMap["date"]
		openIdx := columnIndexMap["open"]
		closeIdx := columnIndexMap["close"]
		for i := 0; i < len(fileData.Data); i++ {
			if len(fileData.Data[i]) < dateIdx+1 || len(fileData.Data[i]) < openIdx+1 || len(fileData.Data[i]) < closeIdx+1 {
				return nil, utils.Errorf(nil, "data error %+v", fileData.Data[i])
			}
			dateTime, err := time.Parse("2006-01-02", fileData.Data[i][dateIdx])
			if err != nil {
				return nil, utils.Errorf(err, "time.Parse fail")
			}
			weekDay := dateTime.Weekday()
			if weekDay == time.Sunday || weekDay == time.Saturday {
				continue
			}
			open, err := strconv.ParseFloat(fileData.Data[i][openIdx], 64)
			if err != nil {
				return nil, utils.Errorf(err, "strconv.ParseFloat fail")
			}
			close, err := strconv.ParseFloat(fileData.Data[i][closeIdx], 64)
			if err != nil {
				return nil, utils.Errorf(err, "strconv.ParseFloat fail")
			}
			if result[int64(weekDay)-1] == nil {
				result[int64(weekDay)-1] = &WeekDayData{}
			}
			result[int64(weekDay)-1].Total++
			if close > open {
				result[int64(weekDay)-1].Win++
			}
		}
		utils.Log("log date " + startTime.Format("2006-01-02"))

		startTime = startTime.AddDate(0, 1, 0)
	}
	for i := 0; i < len(result); i++ {
		if result[i] == nil || result[i].Total <= 0 {
			continue
		}
	}
	return result, nil
}
