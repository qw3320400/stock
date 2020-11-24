package main

import (
	"flag"
	"fmt"
	"stock-go/cmd"
	"stock-go/utils"
)

// go run src/stock-go/main.go -command collectdata data_source baostock data_code all start_date 2015-04-01 end_date 2015-05-01

func main() {
	utils.Log("==start==")
	defer func() {
		utils.Log("==end==")
	}()

	var (
		command string
		param   map[string]string = map[string]string{}
	)
	flag.StringVar(&command, "command", "", "input command")
	flag.Parse()
	for i := 0; i < len(flag.Args())/2; i++ {
		param[flag.Args()[i*2]] = flag.Args()[i*2+1]
	}

	if cmd.CommandMap[command] == nil {
		utils.LogErr(fmt.Sprintf("command not found %s", command))
		return
	}
	err := cmd.CommandMap[command].Run(param)
	if err != nil {
		utils.LogErr(fmt.Sprintf("command process fail %s", err))
		return
	}

	return
}
