package strategy

import (
	"stock-go/utils"
	"strconv"
	"time"
)

var _ Strategy = &WeekDayStrategy{}

type WeekDayStrategy struct {
	DefaultStrategy
	DayCount      int
	lastValue     float64
	lastCost      float64
	lastValueList []float64
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
	var avg20 float64
	if s.stepIndex-s.DayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-s.DayCount+1; i-- {
			avg20 += s.lastValueList[i]
		}
		avg20 = avg20 / float64(s.DayCount)
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
	if avg20 > 0 && len(s.baostockLocalData.StockDateList) > s.stepIndex+1 {
		nextTradeDateWeekDay := s.baostockLocalData.StockDateList[s.stepIndex+1].Time.Weekday()
		if avg20 < close {
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
		} else if avg20 > close {
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
