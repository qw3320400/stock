package strategy

import (
	"encoding/json"
	"fmt"
	"stock-go/utils"
	"time"
)

type RunStrategyRequest struct {
	Code      string `json:"code"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	// internal
	startTime time.Time `json:"-"`
	endTime   time.Time `json:"-"`
}

func RunStrategy(request *RunStrategyRequest) error {

	return nil
}

func run(s Strategy) error {
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
	err := run(ds)
	if err != nil || ds.Result == nil {
		return utils.Errorf(err, "run fail")
	}
	defaultData := ds.Result

	wsAvg5 := &WeekDayStrategy{
		DayCount: 5,
	}
	err = run(wsAvg5)
	if err != nil || wsAvg5.Result == nil {
		return utils.Errorf(err, "run fail")
	}
	wsAvg5Data := wsAvg5.Result

	wsAvg20 := &WeekDayStrategy{
		DayCount: 20,
	}
	err = run(wsAvg20)
	if err != nil || wsAvg20.Result == nil {
		return utils.Errorf(err, "run fail")
	}
	wsAvg20Data := wsAvg20.Result

	chartData = append(chartData, []interface{}{"Date", ds.Code, wsAvg5.Code + "_weekday_5", wsAvg20.Code + "_weekday_20"})
	wsIdx := 0
	for i := 0; i < len(defaultData.LineData); i++ {
		for j := wsIdx; j < len(wsAvg5Data.LineData); j++ {
			if wsAvg5Data.LineData[j].Time == defaultData.LineData[i].Time {
				wsIdx = j
				break
			}
		}
		for j := wsIdx; j < len(wsAvg20Data.LineData); j++ {
			if wsAvg20Data.LineData[j].Time == defaultData.LineData[i].Time {
				wsIdx = j
				break
			}
		}
		chartData = append(chartData, []interface{}{defaultData.LineData[i].Time, defaultData.LineData[i].Value, wsAvg5Data.LineData[wsIdx].Value, wsAvg20Data.LineData[wsIdx].Value})
	}

	body, err := json.Marshal(chartData)
	if err != nil {
		return utils.Errorf(err, "json.Marshal fail")
	}
	utils.Log(string(body))

	utils.Log(fmt.Sprintf("%+v %+v %+v", ds.Code, defaultData.AnualReturnRate, defaultData.DrawDown))
	utils.Log(fmt.Sprintf("%+v %+v %+v", wsAvg5.Code+"_weekday_5", wsAvg5Data.AnualReturnRate, wsAvg5Data.DrawDown))
	utils.Log(fmt.Sprintf("%+v %+v %+v", wsAvg20.Code+"_weekday_20", wsAvg20Data.AnualReturnRate, wsAvg20Data.DrawDown))
	return nil
}

func compareAverage() error {
	chartData := [][]interface{}{}

	as5 := &AverageStrategy{
		DayCount: 5,
	}
	err := run(as5)
	if err != nil || as5.Result == nil {
		return utils.Errorf(err, "run fail")
	}
	as5Data := as5.Result

	as10 := &AverageStrategy{
		DayCount: 10,
	}
	err = run(as10)
	if err != nil || as10.Result == nil {
		return utils.Errorf(err, "run fail")
	}
	as10Data := as10.Result

	as20 := &AverageStrategy{
		DayCount: 20,
	}
	err = run(as20)
	if err != nil || as20.Result == nil {
		return utils.Errorf(err, "run fail")
	}
	as20Data := as20.Result

	chartData = append(chartData, []interface{}{"Date", as5.Code + "_avg5", as10.Code + "_avg10", as20.Code + "_avg20"})
	idx := 0
	for i := 0; i < len(as10Data.LineData); i++ {
		for j := idx; j < len(as10Data.LineData); j++ {
			if as10Data.LineData[j].Time == as5Data.LineData[i].Time {
				idx = j
				break
			}
		}
		for j := idx; j < len(as20Data.LineData); j++ {
			if as20Data.LineData[j].Time == as5Data.LineData[i].Time {
				idx = j
				break
			}
		}
		chartData = append(chartData, []interface{}{as5Data.LineData[i].Time, as5Data.LineData[i].Value, as10Data.LineData[idx].Value, as20Data.LineData[idx].Value})
	}

	body, err := json.Marshal(chartData)
	if err != nil {
		return utils.Errorf(err, "json.Marshal fail")
	}
	utils.Log(string(body))

	return nil
}
