package analyse

import "testing"

func TestBuildRelativeLineChart(t *testing.T) {
	err := BuildRelativeLineChart()
	if err != nil {
		t.Fatal(err)
	}
}
