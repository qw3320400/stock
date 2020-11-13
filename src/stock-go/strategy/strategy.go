package strategy

import (
	"stock-go/utils"
	"time"
)

const (
	DefaultStartTimeStr = "2017-01-01"
	DefaultEndTimeStr   = "2020-10-31"
	DefaultCode         = "sh.000300"
)

type Strategy interface {
	Run() error
}

type StrategyResult struct {
	LineData []*PointData
}

type PointData struct {
	Time  time.Time
	Value float64
}

type DefaultStrategy struct {
	// input
	StartTimeStr string
	EndTimeStr   string
	Code         string
	// output
	Result *StrategyResult
	// internal
	startTime time.Time
	endTime   time.Time
}

func (s *DefaultStrategy) Run() error {
	if s == nil {
		return utils.Errorf(nil, "s is nil")
	}
	if s.StartTimeStr == "" {
		s.StartTimeStr = DefaultStartTimeStr
	}
	if s.EndTimeStr == "" {
		s.EndTimeStr = DefaultEndTimeStr
	}
	if s.Code == "" {
		s.Code = DefaultCode
	}
	var (
		err error
	)
	s.startTime, err = time.Parse("2006-01-02", s.StartTimeStr)
	if err != nil {
		return utils.Errorf(err, "time.Parse fail")
	}
	s.endTime, err = time.Parse("2006-01-02", s.EndTimeStr)
	if err != nil {
		return utils.Errorf(err, "time.Parse fail")
	}
	s.LoadData()

	return nil
}

func (s *DefaultStrategy) LoadData() error {
	if s == nil {
		return utils.Errorf(nil, "s is nil")
	}
	utils.Log("000001")

	return nil
}
