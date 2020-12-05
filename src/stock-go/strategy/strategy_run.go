package strategy

import (
	"encoding/json"
	"stock-go/utils"
)

func RunStrategy(request map[string]string) error {
	if request == nil || request["tag"] == "" {
		return utils.Errorf(nil, "param error %+v", request)
	}
	var strategy Strategy
	switch request["tag"] {
	case "weekday":
		strategy = &WeekDayStrategy{}
	case "rollreturn":
		strategy = &RollingReturn{}
	case "default":
		strategy = &DefaultStrategy{}
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return utils.Errorf(err, "json.Marshal fail")
	}
	err = json.Unmarshal(requestBody, strategy)
	if err != nil {
		return utils.Errorf(err, "json.Unmarshal fail")
	}
	return run(strategy)
}

func run(s Strategy) error {
	if s == nil {
		return utils.Errorf(nil, "s is nil")
	}
	err := s.Init()
	if err != nil {
		return utils.Errorf(err, "s.Init fail")
	}
	err = s.LoadData()
	if err != nil {
		return utils.Errorf(err, "s.LoadData fail")
	}
	for {
		ok, err := s.Step()
		if err != nil {
			return utils.Errorf(err, "s.Step fail")
		}
		if !ok {
			break
		}
	}
	err = s.Final()
	if err != nil {
		return utils.Errorf(err, "s.Final fail")
	}
	return nil
}
