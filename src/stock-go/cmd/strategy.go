package cmd

import "stock-go/collectdata"

/*
func init() {
	setCommand("strategy", &CMDStrategy{})
}

type CMDStrategy struct {
}

func (*CMDStrategy) Run(param map[string]string) error {
	paramBody, err := json.Marshal(param)
	if err != nil {
		return utils.Errorf(err, "json.Marshal fail")
	}
	request := &strategy.RunStrategyRequest{}
	err = json.Unmarshal(paramBody, request)
	if err != nil {
		return utils.Errorf(err, "json.Unmarshal fail")
	}
	return strategy.RunStrategy(request)
}

*/

func init() {
	CommandMap["test"] = &CMDStrategy{}
}

type CMDStrategy struct {
}

func (*CMDStrategy) Run(param map[string]string) error {

	return collectdata.FileToMysql()
}
