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

	// 导数据
	err := exportdata.ExportBaostockData()
	if err != nil {
		utils.LogErr(fmt.Sprint("exportdata.ExportBaostockData fail", err))
		return
	}

	return
}
