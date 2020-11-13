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

	ds := &DefaultStrategy{}
	err := Run(ds)
	if err != nil || ds.Result == nil {
		return utils.Errorf(err, "Run fail")
	}
	defaultData := ds.Result

	ws := &WeekDayStrategy{}
	err = Run(ws)
	if err != nil || ws.Result == nil {
		return utils.Errorf(err, "Run fail")
	}
	weekDayData := ws.Result

	chartData = append(chartData, []interface{}{"Date", ds.Code, ws.Code + "_weekday"})
	wsIdx := 0
	for i := 0; i < len(defaultData.LineData); i++ {
		for j := wsIdx; j < len(weekDayData.LineData); j++ {
			if weekDayData.LineData[j].Time == defaultData.LineData[i].Time {
				wsIdx = j
				break
			}
		}
		chartData = append(chartData, []interface{}{defaultData.LineData[i].Time, defaultData.LineData[i].Value, weekDayData.LineData[wsIdx].Value})
	}

	body, err := json.Marshal(chartData)
	if err != nil {
		return utils.Errorf(err, "json.Marshal fail")
	}
	utils.Log(string(body))
	return nil
}
