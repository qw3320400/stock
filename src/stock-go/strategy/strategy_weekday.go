package strategy

import (
	"stock-go/utils"
	"strconv"
	"time"
)

var _ Strategy = &WeekDayStrategy{}

type WeekDayStrategy struct {
	DefaultStrategy
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
	// 20日均值
	var avg20 float64
	if s.stepIndex-19 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-19; i-- {
			avg20 += s.lastValueList[i]
		}
		avg20 = avg20 / 20
	}
	if s.stepIndex == 0 {
		s.lastValue = 1
	}
	var hit bool
	if avg20 > 0 && len(s.baostockLocalData.StockDateList) > s.stepIndex+1 {
		nextTradeDateWeekDay := s.baostockLocalData.StockDateList[s.stepIndex+1].Time.Weekday()
		if avg20 < close {
			// 牛市收盘
			if nextTradeDateWeekDay == time.Monday || nextTradeDateWeekDay == time.Friday {
				// buy
				if s.lastCost == 0 {
					point.Value = s.lastValue
					s.lastCost = close
				}
			}
			if nextTradeDateWeekDay == time.Tuesday || nextTradeDateWeekDay == time.Thursday {
				// sell
				if s.lastCost != 0 {
					point.Value = s.lastValue * close / s.lastCost
					s.lastValue = point.Value
					s.lastCost = 0
				}
			}
		} else if avg20 > close {
			// 熊市收盘
			if nextTradeDateWeekDay == time.Tuesday || nextTradeDateWeekDay == time.Wednesday {
				// buy
				if s.lastCost == 0 {
					point.Value = s.lastValue
					s.lastCost = close
				}
			}
			if nextTradeDateWeekDay == time.Monday || nextTradeDateWeekDay == time.Thursday {
				// sell
				if s.lastCost != 0 {
					point.Value = s.lastValue * close / s.lastCost
					s.lastValue = point.Value
					s.lastCost = 0
				}
			}
		}
	}
	if !hit {
		if s.lastCost == 0 {
			// 没有持仓
			point.Value = s.lastValue
		} else {
			// 有持仓
			point.Value = s.lastValue * close / s.lastCost
			s.lastValue = point.Value
		}
	}
	s.Result.LineData = append(s.Result.LineData, point)
	s.stepIndex++
	return true, nil
}
