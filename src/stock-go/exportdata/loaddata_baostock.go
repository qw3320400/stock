package exportdata

import (
	"fmt"
	"path/filepath"
	"stock-go/utils"
	"time"
)

type LoadStockDataRequest struct {
	StartTime  time.Time
	EndTime    time.Time
	Code       string
	Frequency  string
	AdjustFlag string
}

type LoadStockDataResponse struct {
	StockDateList []*StockDate
}

type StockDate struct {
	Time time.Time
	Map  map[string]string
}

func LoadBaostockLocalData(request *LoadStockDataRequest) (*LoadStockDataResponse, error) {
	if request == nil || request.Code == "" || request.Frequency == "" || request.AdjustFlag == "" {
		return nil, utils.Errorf(nil, "参数错误 %+v", request)
	}
	response := &LoadStockDataResponse{
		StockDateList: []*StockDate{},
	}
	for {
		if request.StartTime.After(request.EndTime) {
			break
		}

		file := filepath.Join(DataPath,
			fmt.Sprintf(StockPath, request.StartTime.Year(), request.StartTime.Month()),
			fmt.Sprintf(StockFileName, request.Code, request.StartTime.Format("2006-01"), request.Frequency, request.AdjustFlag))
		fileData, err := utils.ReadCommonCSVFile(file)
		if err != nil {
			return nil, utils.Errorf(err, "utils.ReadCommonCSVFile fail")
		}
		for i := 0; i < len(fileData.Data); i++ {
			var (
				tmp = &StockDate{
					Map: map[string]string{},
				}
				timeExist bool
			)
			for j := 0; j < len(fileData.Data[i]); j++ {
				if len(fileData.Column) < j+1 {
					return nil, utils.Errorf(nil, "数据错误 %s %+v", file, fileData.Data[i])
				}
				column := fileData.Column[j]
				if column == "date" {
					tmp.Time, err = time.Parse("2006-01-02", fileData.Data[i][j])
					if err != nil {
						return nil, utils.Errorf(err, "数据错误 %s %+v", file, fileData.Data[i])
					}
					timeExist = true
				}
				tmp.Map[column] = fileData.Data[i][j]
			}
			if !timeExist {
				return nil, utils.Errorf(err, "数据错误 %s %+v", file, fileData.Data[i])
			}
			response.StockDateList = append(response.StockDateList, tmp)
		}

		request.StartTime = request.StartTime.AddDate(0, 1, 0)
	}
	return response, nil
}
