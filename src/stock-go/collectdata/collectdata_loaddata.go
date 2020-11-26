package collectdata

import (
	"sort"
	"stock-go/common"
	"stock-go/data/mysql/stock"
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
	StockKDateList []*common.StockKData
}

func (r *LoadDataRespose) Len() int {
	if r == nil || r.StockKDateList == nil {
		return 0
	}
	return len(r.StockKDateList)
}

func (r *LoadDataRespose) Less(i, j int) bool {
	if r == nil || r.StockKDateList == nil || i+1 > len(r.StockKDateList) || j+1 > len(r.StockKDateList) {
		return false
	}
	return r.StockKDateList[i].TimeCST.Before(r.StockKDateList[j].TimeCST)
}

func (r *LoadDataRespose) Swap(i, j int) {
	if r == nil || r.StockKDateList == nil || i+1 > len(r.StockKDateList) || j+1 > len(r.StockKDateList) {
		return
	}
	tmp := r.StockKDateList[i]
	r.StockKDateList[i] = r.StockKDateList[j]
	r.StockKDateList[j] = tmp
}

func LoadData(request *LoadDataRequest) (*LoadDataRespose, error) {
	if request == nil || request.StartTime.Unix() <= 0 || request.EndTime.Unix() <= 0 ||
		request.Code == "" {
		return nil, utils.Errorf(nil, "参数错误 %+v", request)
	}
	if request.Frequency == "" {
		request.Frequency = "d"
	}
	if request.AdjustFlag == "" {
		request.AdjustFlag = "no"
	}
	response := &LoadDataRespose{
		StockKDateList: []*common.StockKData{},
	}
	dbResponse, err := stock.GetStockKData(&stock.GetStockKDataRequest{
		Code:       request.Code,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		Frequency:  request.Frequency,
		AdjustFlag: request.AdjustFlag,
	})
	if err != nil {
		return nil, utils.Errorf(err, "stock.GetStockKData fail")
	}
	response.StockKDateList = dbResponse.StockKDataList
	tmpTime := request.StartTime
	for {
		if request.EndTime.Before(tmpTime) {
			break
		}
		startTime := tmpTime
		endTime := time.Date(tmpTime.Year(), tmpTime.Month()+1, int(1), int(0), int(0), int(0), int(0), time.UTC).Add(time.Hour * -24)
		if endTime.After(request.EndTime) {
			endTime = request.EndTime
		}

		dbCount := 0
		for _, data := range response.StockKDateList {
			if data.TimeCST.Year() == tmpTime.Year() && data.TimeCST.Month() == tmpTime.Month() {
				dbCount++
				break
			}
		}
		if dbCount > 0 {
			tmpTime = time.Date(tmpTime.Year(), tmpTime.Month()+1, int(1), int(0), int(0), int(0), int(0), time.UTC)
			continue
		}
		err = loadStockKDataAndSave(request.Code, startTime, endTime, request.Frequency, request.AdjustFlag)
		if err != nil {
			return nil, utils.Errorf(err, "loadStockKDataAndSave fail")
		}
		resp, err := stock.GetStockKData(&stock.GetStockKDataRequest{
			Code:       request.Code,
			StartTime:  startTime,
			EndTime:    endTime,
			Frequency:  request.Frequency,
			AdjustFlag: request.AdjustFlag,
		})
		if err != nil {
			return nil, utils.Errorf(err, "stock.GetStockKData fail")
		}
		response.StockKDateList = append(response.StockKDateList, resp.StockKDataList...)

		tmpTime = time.Date(tmpTime.Year(), tmpTime.Month()+1, int(1), int(0), int(0), int(0), int(0), time.UTC)
	}
	sort.Sort(response)
	return response, nil
}
