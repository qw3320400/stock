package analyse

import "testing"

func TestBuildRelativeLineChart(t *testing.T) {
	err := BuildRelativeLineChart()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWeekDay(t *testing.T) {
	result, err := WeekDay()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(result); i++ {
		t.Log(i+1, float64(result[i].Win)/float64(result[i].Total))
	}
}
