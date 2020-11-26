package common

import (
	"time"
)

type StockCode struct {
	ID                     int64
	Code                   string
	Name                   string
	Industry               string
	IndustryClassification string
}

type StockTradeDate struct {
	ID           int64
	DateCST      time.Time
	IsTradingDay bool
}

type StockKData struct {
	ID          int64
	Code        string
	TimeCST     time.Time
	Frequency   string
	AdjustFlag  string
	Open        string
	High        string
	Low         string
	Close       string
	Preclose    string
	Volume      string
	Amount      string
	Turn        string
	TradeStatus string
	PctChg      string
	PeTTM       string
	PsTTM       string
	PcfNcfTTM   string
	PbMRQ       string
	IsST        bool
}

type StockStrategyData struct {
	ID                    int64
	StockStrategyResultID int64
	Code                  string
	Tag                   string
	TimeCST               time.Time
	Value                 string
}

type StockStrategyResult struct {
	ID              int64
	Code            string
	Tag             string
	StartTimeCST    time.Time
	EndTimeCST      time.Time
	AnualReturnRate string
	DrawDown        string
}
