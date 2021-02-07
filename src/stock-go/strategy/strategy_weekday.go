package strategy

import (
	"stock-go/utils"
	"strconv"
	"time"
)

var _ Strategy = &WeekDayStrategy{}

type WeekDayStrategy struct {
	DefaultStrategy
	DayCountStr string `json:"day_count"`
	// internal
	lastValue     float64
	lastCost      float64
	lastValueList []float64
}

func (s *WeekDayStrategy) Init() error {
	s.Tag += "_test2"
	return s.DefaultStrategy.Init()
}

func (s *WeekDayStrategy) Step() (bool, error) {
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
	var (
		avg5, avg20 float64
		dayCount    int64 = 5
	)
	if s.stepIndex-dayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-dayCount+1; i-- {
			avg5 += s.lastValueList[i]
		}
		avg5 = avg5 / float64(dayCount)
	}
	dayCount = 20
	if s.stepIndex-dayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-dayCount+1; i-- {
			avg20 += s.lastValueList[i]
		}
		avg20 = avg20 / float64(dayCount)
	}
	if s.stepIndex == 0 {
		s.lastValue = 1
	}
	if s.lastCost != 0 {
		// 有持仓
		point.Value = s.lastValue * close / s.lastCost
	} else {
		point.Value = s.lastValue
	}
	// 策略
	var opt string = "-"
	if avg5 > 0 && avg20 > 0 && int64(len(s.baostockLocalData.StockKDateList)) > s.stepIndex+1 {
		nextTradeDateWeekDay := s.baostockLocalData.StockKDateList[s.stepIndex+1].TimeCST.Weekday()
		if (avg20 > close && nextTradeDateWeekDay == time.Tuesday) ||
			(avg5 < close && nextTradeDateWeekDay == time.Monday || nextTradeDateWeekDay == time.Friday) {
			// 买入
			opt = "buy"
		} else {
			// 卖出
			opt = "sell"
		}
	}
	if opt == "buy" {
		if s.lastCost == 0 {
			s.lastCost = close
		}
	} else if opt == "sell" {
		if s.lastCost != 0 {
			s.lastCost = 0
			s.lastValue = point.Value
		}
	}
	s.Result.LineData = append(s.Result.LineData, point)
	s.stepIndex++
	return true, nil
}
