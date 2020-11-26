package cmd

import (
	"stock-go/strategy"
)

func init() {
	setCommand("strategy", &CMDStrategy{})
}

type CMDStrategy struct {
}

func (*CMDStrategy) Run(param map[string]string) error {

	return strategy.RunStrategy(param)
}
