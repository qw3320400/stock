package strategy

import (
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
