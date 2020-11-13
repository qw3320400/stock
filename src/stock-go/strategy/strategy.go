package strategy

import (
	"encoding/json"
	"stock-go/utils"
	"time"
)

const (
	DefaultStartTimeStr = "2017-01-01"
	DefaultEndTimeStr   = "2020-10-31"
	DefaultCode         = "sh.000300"
)

type StrategyResult struct {
	LineData []*PointData
}

type PointData struct {
	Time  time.Time
	Value float64
}

type Strategy interface {
	Init() error
	LoadData() error
	Step() (bool, error)
	Final() error
}

func Run(s Strategy) error {
	if s == nil {
		return utils.Errorf(nil, "s is nil")
	}
	err := s.Init()
	if err != nil {
		return utils.Errorf(err, "s.Init fail")
	}
	err = s.LoadData()
	if err != nil {
		return utils.Errorf(err, "s.LoadData fail")
	}
	for {
		ok, err := s.Step()
		if err != nil {
			return utils.Errorf(err, "s.Step fail")
		}
		if !ok {
			break
		}
	}
	err = s.Final()
	if err != nil {
		return utils.Errorf(err, "s.Final fail")
	}
	return nil
}

func compareWeekDayAndDefault() error {
	chartData := [][]interface{}{}

	s := &DefaultStrategy{}
	err := Run(s)
	if err != nil || s.Result == nil {
		return utils.Errorf(err, "Run fail")
	}
	defaultData := s.Result

	chartData = append(chartData, []interface{}{"Date", s.Code})
	for i := 0; i < len(defaultData.LineData); i++ {
		chartData = append(chartData, []interface{}{defaultData.LineData[i].Time, defaultData.LineData[i].Value})
	}

	body, err := json.Marshal(chartData)
	if err != nil {
		return utils.Errorf(err, "json.Marshal fail")
	}
	utils.Log(string(body))
	return nil
}
