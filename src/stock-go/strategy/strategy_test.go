package strategy

import (
	"testing"
)

func TestStrategy(t *testing.T) {
	err := compareWeekDayAndDefault()
	if err != nil {
		t.Fatal(err)
	}
}
