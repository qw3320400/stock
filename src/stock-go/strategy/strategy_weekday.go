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
	dayCount      int64
	lastValue     float64
	lastCost      float64
	lastValueList []float64
}

func (s *WeekDayStrategy) Init() error {
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
	var avg float64
	if s.stepIndex-s.dayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-s.dayCount+1; i-- {
			avg += s.lastValueList[i]
		}
		avg = avg / float64(s.dayCount)
	}
	if s.stepIndex == 0 {
		s.lastValue = 1
	} else {
		if s.lastCost != 0 {
			// 有持仓
			s.lastValue = s.lastValue * close / s.lastCost
		}
	}
	point.Value = s.lastValue
	// 策略
	var opt string = "-"
	if avg > 0 && int64(len(s.baostockLocalData.StockKDateList)) > s.stepIndex+1 {
		nextTradeDateWeekDay := s.baostockLocalData.StockKDateList[s.stepIndex+1].TimeCST.Weekday()
		if avg < close {
			// 牛市收盘
			if nextTradeDateWeekDay == time.Friday || nextTradeDateWeekDay == time.Monday || nextTradeDateWeekDay == time.Tuesday {
				// 买入
				if s.lastCost == 0 {
					opt = "buy"
				}
			} else if nextTradeDateWeekDay == time.Wednesday || nextTradeDateWeekDay == time.Thursday {
				// 卖出
				if s.lastCost != 0 {
					opt = "sell"
				}
			}
		} else if avg > close {
			// 熊市收盘
			// 卖出
			if s.lastCost != 0 {
				opt = "sell"
			}
		}
	}
	if opt == "buy" {
		s.lastCost = close
	} else if opt == "sell" {
		s.lastCost = 0
	}
	s.Result.LineData = append(s.Result.LineData, point)
	s.stepIndex++
	return true, nil
}
