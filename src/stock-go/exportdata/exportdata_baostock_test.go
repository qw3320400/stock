package exportdata

import (
	"stock-go/thirdparty/baostock"
	"testing"
	"time"
)

func TestErrorData(t *testing.T) {
	// 连接
	bc, err := baostock.NewBaostockConnection()
	if err != nil {
		t.Fatalf("[ExportBaostockData] baostock.NewBaostockConnection fail\n\t%s", err)
	}
	defer func() {
		bc.CloseConnection()
	}()
	// 登陆
	err = bc.Login("", "", 0)
	if err != nil {
		t.Fatalf("[ExportBaostockData] bc.Login fail\n\t%s", err)
	}
	defer func() {
		bc.Logout()
	}()

	startTime := time.Date(2020, 4, 1, int(0), int(0), int(0), int(0), time.UTC)
	endTime := time.Date(2020, 4, 30, int(0), int(0), int(0), int(0), time.UTC)
	queryAndSaveBaostockKData(bc, "sz.002188", startTime, endTime, "/Users/k/Desktop/code/stock/data/baostock/error",
		[]string{"5", "date,code,open,high,low,close,volume,amount,adjustflag"}, "3")
}
