package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"stock-go/cmd"
	"stock-go/data/mysql"
	"stock-go/thirdparty/baostock"
	"stock-go/utils"
	"syscall"
)

// go run src/stock-go/main.go -command collectdata data_source baostock data_code all start_date 2015-04-01 end_date 2015-04-30
// go run src/stock-go/main.go -command strategy tag default code sh.000300 start_date 2016-01-01 end_date 2020-10-31
// go run src/stock-go/main.go -command strategy tag weekday code sh.000300 start_date 2016-01-01 end_date 2020-10-31 day_count 5
// go run src/stock-go/main.go -command strategy tag rollreturn result_id 9
// go run src/stock-go/main.go -command strategy tag default code 000300.XSHG start_date 2006-01-01 end_date 2020-12-21 data_source jqdata
// go run src/stock-go/main.go -command strategy tag weekday code 510300.XSHG start_date 2006-01-01 end_date 2020-12-21 day_count 5 data_source jqdata

// grafana-server --config=/usr/local/etc/grafana/grafana.ini --homepath /usr/local/share/grafana cfg:default.paths.logs=/usr/local/var/log/grafana cfg:default.paths.data=/usr/local/var/lib/grafana cfg:default.paths.plugins=/usr/local/var/lib/grafana/plugins

func main() {
	utils.Log("==start==")
	defer func() {
		utils.Log("==end==")
	}()

	// 创建监听退出chan
	c := make(chan os.Signal)
	// 监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				utils.Log(fmt.Sprintf("Program Exit... %+v", s))
				onExit()
				os.Exit(1)
			default:
				utils.Log(fmt.Sprintf("other signal %+v", s))
			}
		}
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

func onExit() {
	err := baostock.CloseBaostockConnection()
	if err != nil {
		utils.LogErr(fmt.Sprintf("baostock.CloseBaostockConnection fail %s", err))
	}
	err = mysql.CloseConnection()
	if err != nil {
		utils.LogErr(fmt.Sprintf("mysql.CloseConnection fail %s", err))
	}
}
