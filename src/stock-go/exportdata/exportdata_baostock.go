package exportdata

import (
	"errors"
	"fmt"
	"stock-go/thirdparty/baostock"
)

func ExportBaostockData() error {
	// 连接
	bc, err := baostock.NewBaostockConnection()
	if err != nil {
		return errors.New(fmt.Sprint("[ExportBaostockData] baostock.NewBaostockConnection fail\n\t", err))
	}
	defer func() {
		bc.CloseConnection()
	}()
	// 登陆
	err = bc.Login("", "", 0)
	if err != nil {
		return errors.New(fmt.Sprint("[ExportBaostockData] bc.Login fail\n\t", err))
	}
	defer func() {
		bc.Logout()
	}()
	// 测试查询数据
	_, err = bc.QueryHistoryKDataPlus("sh.600710",
		"date,code,open,high,low,close,volume,amount,adjustflag,turn,pctChg",
		"2016-07-01", "2016-07-31",
		"m", "3")
	if err != nil {
		return errors.New(fmt.Sprint("[ExportBaostockData] bc.QueryHistoryKDataPlus fail\n\t", err))
	}
	return nil
}
