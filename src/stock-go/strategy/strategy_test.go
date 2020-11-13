package strategy

import "testing"

func TestStrategy(t *testing.T) {
	s := &DefaultStrategy{
		StartTimeStr: "1111",
	}
	err := s.Run()
	if err != nil {
		t.Fatal(err)
	}
}
