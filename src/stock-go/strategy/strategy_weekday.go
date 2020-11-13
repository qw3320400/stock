package strategy

var _ Strategy = &WeekDayStrategy{}

type WeekDayStrategy struct {
	DefaultStrategy
}
