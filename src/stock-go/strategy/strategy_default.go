package strategy

import (
	"fmt"
	"stock-go/exportdata"
	"stock-go/utils"
	"strconv"
	"time"
)

var _ Strategy = &DefaultStrategy{}

type DefaultStrategy struct {
	// input
	StartTimeStr string
	EndTimeStr   string
	Code         string
	// output
	Result *StrategyResult
	// internal
	startTime         time.Time
	endTime           time.Time
	baostockLocalData *exportdata.LoadStockDataResponse
	stepIndex         int
	baseCost          float64
}

func (s *DefaultStrategy) Init() error {
	if s == nil {
		return utils.Errorf(nil, "s is nil")
	}
	if s.StartTimeStr == "" {
		s.StartTimeStr = DefaultStartTimeStr
	}
	if s.EndTimeStr == "" {
		s.EndTimeStr = DefaultEndTimeStr
	}
	if s.Code == "" {
		s.Code = DefaultCode
	}
	var (
		err error
	)
	s.startTime, err = time.Parse("2006-01-02", s.StartTimeStr)
	if err != nil {
		return utils.Errorf(err, "time.Parse fail")
	}
	s.endTime, err = time.Parse("2006-01-02", s.EndTimeStr)
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
	s.baostockLocalData, err = exportdata.LoadBaostockLocalData(&exportdata.LoadStockDataRequest{
		StartTime:  s.startTime,
		EndTime:    s.endTime,
		Code:       s.Code,
		Frequency:  "d",
		AdjustFlag: "3",
	})
	if err != nil {
		return utils.Errorf(err, "exportdata.LoadBaostockLocalData fail")
	}
	utils.Log(fmt.Sprintf("%+v", len(s.baostockLocalData.StockDateList)))
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
	if len(s.baostockLocalData.StockDateList) < s.stepIndex+1 {
		return false, nil
	}
	point := &PointData{
		Time: s.baostockLocalData.StockDateList[s.stepIndex].Time,
	}
	closeStr := s.baostockLocalData.StockDateList[s.stepIndex].Map["close"]
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
	return nil
}
