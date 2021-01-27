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
	weekDay       map[int]*dayData
	lastValueList []float64
}

type dayData struct {
	Up    int
	Total int
}

func (s *WaveStrategy) Init() error {
	s.Tag += ("_" + "3")
	s.weekDay = map[int]*dayData{}
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

	if len(s.lastValueList) > 1 {
		weekDay := int(s.baostockLocalData.StockKDateList[s.stepIndex].TimeCST.Weekday())
		if s.weekDay[weekDay] == nil {
			s.weekDay[weekDay] = &dayData{}
		}
		s.weekDay[weekDay].Total++
		length := len(s.lastValueList)
		if s.lastValueList[length] > s.lastValueList[length-1] {
			s.weekDay[weekDay].Up++
		}
	}

	s.Result.LineData = append(s.Result.LineData, point)
	s.stepIndex++
	return true, nil
}

func (s *WaveStrategy) Final() error {
	for k, v := range s.weekDay {
		utils.Log(fmt.Sprintf("%d %+v", k, float64(v.Up)/float64(v.Total)))
	}
	return nil
}
