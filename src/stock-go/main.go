package main

import (
	"fmt"
	"stock-go/exportdata"
	"stock-go/utils"
)

func main() {
	utils.Log("==start==")
	defer func() {
		utils.Log("==end==")
	}()

	// 测试导数据
	err := exportdata.ExportBaostockData()
	if err != nil {
		utils.LogErr(fmt.Sprint("[main] exportdata.ExportBaostockData fail", err))
		return
	}

	return
}
