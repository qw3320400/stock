package cmd

import (
	"stock-go/collectdata"
	"stock-go/utils"

	"encoding/json"
)

func init() {
	setCommand("collectdata", &CMDCollectData{})
}

type CMDCollectData struct {
}

func (*CMDCollectData) Run(param map[string]string) error {
	paramBody, err := json.Marshal(param)
	if err != nil {
		return utils.Errorf(err, "json.Marshal fail")
	}
	request := &collectdata.CollectDataRequest{}
	err = json.Unmarshal(paramBody, request)
	if err != nil {
		return utils.Errorf(err, "json.Unmarshal fail")
	}
	return collectdata.CollectData(request)
}
