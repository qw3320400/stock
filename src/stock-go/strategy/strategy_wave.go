package strategy

import (
	"stock-go/collectdata"
	"stock-go/utils"
	"strconv"
)

var _ Strategy = &WaveStrategy{}

type WaveStrategy struct {
	DefaultStrategy
	// internal
	dayCount      int64
	lastValue     float64
	lastCost      float64
	lastValueList []float64
	lastAvgValue  float64
}

func (s *WaveStrategy) Init() error {
	s.dayCount = 5
	s.Tag += "_1"
	return s.DefaultStrategy.Init()
}

func (s *WaveStrategy) LoadData() error {
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
		Frequency:  "15",
	})
	if err != nil {
		return utils.Errorf(err, "collectdata.LoadData fail")
	}
	return nil
}

func (s *WaveStrategy) Step() (bool, error) {
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
	if s.stepIndex == 0 {
		s.lastValue = 1
	}
	var (
		curClose float64
		opt      string
	)
	if s.lastAvgValue > 0 {
		if s.baostockLocalData.StockKDateList[s.stepIndex].TimeCST.Format("15:04:05") == "10:30:00" {
			closeStr := s.baostockLocalData.StockKDateList[s.stepIndex].Close
			close, err := strconv.ParseFloat(closeStr, 64)
			if err != nil {
				return false, utils.Errorf(err, "strconv.ParseFloat fail")
			}
			curClose = close
			if close > s.lastAvgValue {
				opt = "sell"
			} else {
				opt = "buy"
			}
		} else if s.baostockLocalData.StockKDateList[s.stepIndex].TimeCST.Format("15:04:05") == "14:00:00" {
			closeStr := s.baostockLocalData.StockKDateList[s.stepIndex].Close
			close, err := strconv.ParseFloat(closeStr, 64)
			if err != nil {
				return false, utils.Errorf(err, "strconv.ParseFloat fail")
			}
			curClose = close
			if close > s.lastAvgValue {
				opt = "buy"
			} else {
				opt = "sell"
			}
		}
	}
	if s.baostockLocalData.StockKDateList[s.stepIndex].TimeCST.Format("15:04:05") == "15:00:00" {
		closeStr := s.baostockLocalData.StockKDateList[s.stepIndex].Close
		close, err := strconv.ParseFloat(closeStr, 64)
		if err != nil {
			return false, utils.Errorf(err, "strconv.ParseFloat fail")
		}
		s.lastValueList = append(s.lastValueList, close)
		// 均值
		var avg float64
		if len(s.lastValueList) >= int(s.dayCount) {
			for i := len(s.lastValueList) - 1; i >= len(s.lastValueList)-int(s.dayCount); i-- {
				avg += s.lastValueList[i]
			}
			avg = avg / float64(s.dayCount)
		}
		s.lastAvgValue = avg

	}
	if opt == "buy" {
		if s.lastCost == 0 {
			s.lastCost = curClose
			s.Result.LineData = append(s.Result.LineData, &PointData{
				Time:  s.baostockLocalData.StockKDateList[s.stepIndex].TimeCST,
				Value: s.lastValue,
			})
		}
	} else if opt == "sell" {
		if s.lastCost != 0 {
			s.lastValue = s.lastValue * curClose / s.lastCost
			s.lastCost = 0
			s.Result.LineData = append(s.Result.LineData, &PointData{
				Time:  s.baostockLocalData.StockKDateList[s.stepIndex].TimeCST,
				Value: s.lastValue,
			})
		}
	}
	s.stepIndex++
	return true, nil
}
