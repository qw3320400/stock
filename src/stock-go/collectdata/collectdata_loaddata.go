package collectdata

/*
import (
	"stock-go/common"
	"stock-go/mysql"
	"stock-go/utils"
	"time"
)

type LoadDataRequest struct {
	StartTime  time.Time
	EndTime    time.Time
	Code       string
	Frequency  string
	AdjustFlag string
}

type LoadDataRespose struct {
	StockDateList []*common.StockKData
}

func LoadData(request *LoadDataRequest) (*LoadDataRespose, error) {
	if request == nil || request.StartTime.Unix() <= 0 || request.EndTime.Unix() <= 0 ||
		request.Code == "" || request.Frequency == "" || request.AdjustFlag == "" {
		return nil, utils.Errorf(nil, "参数错误 %+v", request)
	}
	var (
		err error
	)
	response := LoadDataRespose{
		StockDateList: []*common.StockKData{},
	}
	response.StockDateList, err = mysql.GetStockKData(&mysql.GetStockKData{
		Code:       request.Code,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		Frequency:  request.Frequency,
		AdjustFlag: request.AdjustFlag,
	})
	if err != nil {
		return nil, utils.Errorf(err, "mysql.GetStockKData fail")
	}
	for tmpTime := request.StartTime; !request.EndTime.Before(tmpTime); tmpTime.AddDate(0, 1, 0) {
		dbCount := 0
		for _, data := range response.StockDateList {
			if data.TimeCST.Year() == tmpTime.Year() && data.TimeCST.Month() == tmpTime.Month() {
				dbCount++
				break
			}
		}
		if dbCount > 0 {
			continue
		}

	}
}
*/
