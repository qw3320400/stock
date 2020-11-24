package strategy

import (
	"stock-go/utils"
	"strconv"
)

var _ Strategy = &AverageStrategy{}

type AverageStrategy struct {
	DefaultStrategy
	DayCount      int
	lastValueList []float64
}

func (s *AverageStrategy) Step() (bool, error) {
	if s == nil || s.baostockLocalData == nil || s.stepIndex < 0 {
		return false, utils.Errorf(nil, "param error %+v", s)
	}
	if s.Result == nil {
		s.Result = &StrategyResult{
			LineData: []*PointData{},
		}
	}
	if s.lastValueList == nil {
		s.lastValueList = []float64{}
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
	s.lastValueList = append(s.lastValueList, close)
	// 均值
	if s.DayCount <= 0 {
		s.DayCount = 5
	}
	var avg float64
	if s.stepIndex-s.DayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-s.DayCount+1; i-- {
			avg += s.lastValueList[i]
		}
		avg = avg / float64(s.DayCount)
		point.Value = avg
	}
	s.Result.LineData = append(s.Result.LineData, point)
	s.stepIndex++
	return true, nil
}
