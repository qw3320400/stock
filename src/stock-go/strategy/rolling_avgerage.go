package strategy

import (
	"stock-go/utils"
	"strconv"
)

var _ Strategy = &RollingAverage{}

type RollingAverage struct {
	DefaultStrategy
	DayCountStr string `json:"day_count"`
	// internal
	dayCount      int64
	lastValueList []float64
}

func (s *RollingAverage) Init() error {
	var err error
	s.dayCount, err = strconv.ParseInt(s.DayCountStr, 10, 64)
	if err != nil {
		return utils.Errorf(err, "strconv.ParseInt fail")
	}
	if s.dayCount <= 0 || s.dayCount > 100 {
		return utils.Errorf(nil, "param error %+v", s)
	}
	s.Tag += ("_" + s.DayCountStr)
	return s.DefaultStrategy.Init()
}

func (s *RollingAverage) Step() (bool, error) {
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
	if int64(len(s.baostockLocalData.StockKDateList)) < s.stepIndex+1 {
		return false, nil
	}
	point := &PointData{
		Time: s.baostockLocalData.StockKDateList[s.stepIndex].TimeCST,
	}
	closeStr := s.baostockLocalData.StockKDateList[s.stepIndex].Close
	close, err := strconv.ParseFloat(closeStr, 64)
	if err != nil {
		return false, utils.Errorf(err, "strconv.ParseFloat fail")
	}
	s.lastValueList = append(s.lastValueList, close)
	// 均值
	var avg float64
	if s.stepIndex-s.dayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-s.dayCount+1; i-- {
			avg += s.lastValueList[i]
		}
		avg = avg / float64(s.dayCount)
	}
	point.Value = avg
	s.Result.LineData = append(s.Result.LineData, point)
	s.stepIndex++
	return true, nil
}
