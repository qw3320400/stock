package analyse

import "testing"

func TestRelativeIndustry(t *testing.T) {
	err := RelativeIndustry()
	if err != nil {
		t.Fatal(err)
	}
}
