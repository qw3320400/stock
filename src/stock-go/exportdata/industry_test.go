package exportdata

import "testing"

func TestAllCodeIndustry(t *testing.T) {
	err := ExportBaostockAllCodeIndustry()
	if err != nil {
		t.Fatal(err)
	}
}
