package strategy

import (
	"encoding/json"
	"fmt"
	"stock-go/utils"
	"time"
)

const (
	DefaultStartTimeStr = "2015-05-01"
	DefaultEndTimeStr   = "2020-10-31"
	DefaultCode         = "sh.000300"
)

type StrategyResult struct {
	AnualReturnRate float64
	DrawDown        float64
	LineData        []*PointData
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

	ws2 := &WeekDayStrategyV2{}
	err = Run(ws2)
	if err != nil || ws2.Result == nil {
		return utils.Errorf(err, "Run fail")
	}
	weekDayDataV2 := ws2.Result

	chartData = append(chartData, []interface{}{"Date", ds.Code, ws.Code + "_weekday", ws2.Code + "_weekday_v2"})
	wsIdx := 0
	for i := 0; i < len(defaultData.LineData); i++ {
		for j := wsIdx; j < len(weekDayData.LineData); j++ {
			if weekDayData.LineData[j].Time == defaultData.LineData[i].Time {
				wsIdx = j
				break
			}
		}
		for j := wsIdx; j < len(weekDayDataV2.LineData); j++ {
			if weekDayDataV2.LineData[j].Time == defaultData.LineData[i].Time {
				wsIdx = j
				break
			}
		}
		chartData = append(chartData, []interface{}{defaultData.LineData[i].Time, defaultData.LineData[i].Value, weekDayData.LineData[wsIdx].Value, weekDayDataV2.LineData[wsIdx].Value})
	}

	body, err := json.Marshal(chartData)
	if err != nil {
		return utils.Errorf(err, "json.Marshal fail")
	}
	utils.Log(string(body))

	utils.Log(fmt.Sprintf("%+v %+v %+v", ds.Code, defaultData.AnualReturnRate, defaultData.DrawDown))
	utils.Log(fmt.Sprintf("%+v %+v %+v", ws.Code+"_weekday", weekDayData.AnualReturnRate, weekDayData.DrawDown))
	utils.Log(fmt.Sprintf("%+v %+v %+v", ws2.Code+"_weekday_v2", weekDayDataV2.AnualReturnRate, weekDayDataV2.DrawDown))
	return nil
}
