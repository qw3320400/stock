package strategy

import (
	"fmt"
	"stock-go/utils"
	"strconv"
)

var _ Strategy = &WaveStrategy{}

type WaveStrategy struct {
	DefaultStrategy
	// internal
	lastValue     float64
	testData      map[string]*data
	lastValueList []float64
}

type data struct {
	Up        int
	Total     int
	TotalRate float64
}

func (s *WaveStrategy) Init() error {
	s.Tag += ("_" + "3")
	s.lastValue = 1
	s.testData = map[string]*data{}
	return s.DefaultStrategy.Init()
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
		avg5, avg10, avg20 float64
	)
	var dayCount int64 = 5
	if s.stepIndex-dayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-dayCount+1; i-- {
			avg5 += s.lastValueList[i]
		}
		avg5 = avg5 / float64(dayCount)
	}
	dayCount = 10
	if s.stepIndex-dayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-dayCount+1; i-- {
			avg10 += s.lastValueList[i]
		}
		avg10 = avg10 / float64(dayCount)
	}
	dayCount = 20
	if s.stepIndex-dayCount+1 >= 0 {
		for i := s.stepIndex; i >= s.stepIndex-dayCount+1; i-- {
			avg20 += s.lastValueList[i]
		}
		avg20 = avg20 / float64(dayCount)
	}

	if len(s.lastValueList) > 1 &&
		avg5 > 0 && avg10 > 0 && avg20 > 0 {
		key := "other"
		if avg5 > avg10 && avg10 > avg20 {
			key = "5-10-20"
		} else if avg5 > avg20 && avg20 > avg10 {
			key = "5-20-10"
		} else if avg10 > avg5 && avg5 > avg20 {
			key = "10-5-20"
		} else if avg10 > avg20 && avg20 > avg5 {
			key = "10-20-5"
		} else if avg20 > avg5 && avg5 > avg20 {
			key = "20-5-10"
		} else if avg20 > avg10 && avg10 > avg5 {
			key = "20-10-5"
		}
		if s.testData[key] == nil {
			s.testData[key] = &data{}
		}
		s.testData[key].Total++
		s.testData[key].TotalRate += (s.lastValueList[len(s.lastValueList)-1]/s.lastValueList[len(s.lastValueList)-2] - 1)
		if s.lastValueList[len(s.lastValueList)-1] > s.lastValueList[len(s.lastValueList)-2] {
			s.testData[key].Up++
		}
	}

	s.Result.LineData = append(s.Result.LineData, point)
	s.stepIndex++
	return true, nil
}

func (s *WaveStrategy) Final() error {
	for k, v := range s.testData {
		utils.Log(fmt.Sprintf("%s - %d %d %+v %+v", k, v.Up, v.Total, float64(v.Up)/float64(v.Total), v.TotalRate/float64(v.Total)))
	}
	return nil
}
