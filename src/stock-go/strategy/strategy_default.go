package strategy

import (
	"math"
	"stock-go/collectdata"
	"stock-go/common"
	"stock-go/data/mysql/stock"
	"stock-go/utils"
	"strconv"
	"time"
)

var _ Strategy = &DefaultStrategy{}

type DefaultStrategy struct {
	// input
	Tag        string `json:"tag"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Code       string `json:"code"`
	DataSource string `json:"data_source"`
	// output
	Result *StrategyResult
	// internal
	startTime         time.Time
	endTime           time.Time
	baostockLocalData *collectdata.LoadDataRespose
	stepIndex         int64
	baseCost          float64
}

func (s *DefaultStrategy) Init() error {
	if s == nil {
		return utils.Errorf(nil, "s is nil")
	}
	if s.StartDate == "" {
		s.StartDate = DefaultStartTimeStr
	}
	if s.EndDate == "" {
		s.EndDate = DefaultEndTimeStr
	}
	if s.Code == "" {
		s.Code = DefaultCode
	}
	if s.DataSource == "" {
		s.DataSource = collectdata.DataSourceBaostock
	}
	var (
		err error
	)
	s.startTime, err = time.Parse("2006-01-02", s.StartDate)
	if err != nil {
		return utils.Errorf(err, "time.Parse fail")
	}
	s.endTime, err = time.Parse("2006-01-02", s.EndDate)
	if err != nil {
		return utils.Errorf(err, "time.Parse fail")
	}
	return nil
}

func (s *DefaultStrategy) LoadData() error {
	if s == nil || s.Code == "" {
		return utils.Errorf(nil, "param error %+v", s)
	}
	var (
		err error
	)
	s.baostockLocalData, err = collectdata.LoadData(&collectdata.LoadDataRequest{
		StartTime:  s.startTime,
		EndTime:    s.endTime,
		Code:       s.Code,
		DataSource: s.DataSource,
	})
	if err != nil {
		return utils.Errorf(err, "collectdata.LoadData fail")
	}
	return nil
}

func (s *DefaultStrategy) Step() (bool, error) {
	if s == nil || s.baostockLocalData == nil || s.stepIndex < 0 {
		return false, utils.Errorf(nil, "param error %+v", s)
	}
	if s.Result == nil {
		s.Result = &StrategyResult{
			LineData: []*PointData{},
		}
	}
	if int64(len(s.baostockLocalData.StockKDateList)) < s.stepIndex+1 {
		return false, nil
	}
	point := &PointData{
		Time: s.baostockLocalData.StockKDateList[s.stepIndex].TimeCST,
	}
	closeStr := s.baostockLocalData.StockKDateList[s.stepIndex].Close
	close, err := strconv.ParseFloat(closeStr, 64)
	if err != nil {
		return false, utils.Errorf(err, "trconv.ParseFloat fail")
	}
	if s.stepIndex == 0 {
		point.Value = 1
		s.baseCost = close
	} else {
		if s.baseCost == 0 {
			return false, utils.Errorf(err, "数据错误 %+v", point)
		}
		point.Value = close / s.baseCost
	}
	s.Result.LineData = append(s.Result.LineData, point)
	s.stepIndex++
	return true, nil
}

func (s *DefaultStrategy) Final() error {
	if s.Result == nil || s.Result.LineData == nil || len(s.Result.LineData) <= 1 {
		return nil
	}
	startTime := s.Result.LineData[0].Time
	endTime := s.Result.LineData[len(s.Result.LineData)-1].Time
	year := float64(endTime.Unix()-startTime.Unix()) / float64(86400*365)
	if year <= 0 {
		return utils.Errorf(nil, "数据错误 %+v %+v", startTime, endTime)
	}
	totalReturnRate := s.Result.LineData[len(s.Result.LineData)-1].Value
	s.Result.AnualReturnRate = math.Pow(totalReturnRate, 1/year) - 1
	var (
		maxDrawDown  float64
		maxValue     float64
		maxDownValue float64
	)
	for i := 0; i < len(s.Result.LineData); i++ {
		if s.Result.LineData[i].Value > maxValue {
			maxValue = s.Result.LineData[i].Value
			maxDownValue = s.Result.LineData[i].Value
		} else {
			if s.Result.LineData[i].Value < maxDownValue {
				maxDownValue = s.Result.LineData[i].Value
				rate := 1 - maxDownValue/maxValue
				if rate > maxDrawDown {
					maxDrawDown = rate
				}
			}
		}
	}
	s.Result.DrawDown = maxDrawDown
	// 保存数据
	resultResp, err := stock.InsertStockStrategyResult(&stock.InsertStockStrategyResultRequest{
		StockStrategyResult: &common.StockStrategyResult{
			Code:            s.Code,
			Tag:             s.Tag,
			StartTimeCST:    s.startTime,
			EndTimeCST:      s.endTime,
			AnualReturnRate: strconv.FormatFloat(s.Result.AnualReturnRate, 'f', -1, 64),
			DrawDown:        strconv.FormatFloat(s.Result.DrawDown, 'f', -1, 64),
		},
	})
	if err != nil {
		return utils.Errorf(err, "stock.InsertStockStrategyResult fail")
	}
	strategyDataRequest := &stock.InsertStockStrategyDataRequest{
		StockStrategyDataList: []*common.StockStrategyData{},
	}
	for i := 0; i < len(s.Result.LineData); i++ {
		strategyDataRequest.StockStrategyDataList = append(strategyDataRequest.StockStrategyDataList, &common.StockStrategyData{
			StockStrategyResultID: resultResp.StockStrategyResult.ID,
			Code:                  s.Code,
			Tag:                   s.Tag,
			TimeCST:               s.Result.LineData[i].Time,
			Value:                 strconv.FormatFloat(s.Result.LineData[i].Value, 'f', -1, 64),
		})
	}
	err = stock.InsertStockStrategyData(strategyDataRequest)
	if err != nil {
		return utils.Errorf(err, "stock.InsertStockStrategyData fail")
	}

	return nil
}
