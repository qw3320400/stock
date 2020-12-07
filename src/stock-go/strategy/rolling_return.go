package strategy

import (
	"fmt"
	"stock-go/data/mysql/stock"
	"stock-go/utils"
	"strconv"
)

var _ Strategy = &RollingReturn{}

type RollingReturn struct {
	DefaultStrategy
	ResultIDStr string `json:"result_id"`
	// internal
	resultID   int64
	resultData *stock.GetStockStrategyResultDataResponse
}

func (r *RollingReturn) Init() error {
	if r == nil || r.ResultIDStr == "" {
		return utils.Errorf(nil, "param error %+v", r)
	}
	var (
		err error
	)
	r.resultID, err = strconv.ParseInt(r.ResultIDStr, 10, 64)
	if err != nil {
		return utils.Errorf(err, "strconv.ParseInt fail")
	}
	if r.resultID <= 0 {
		return utils.Errorf(nil, "param error %+v", r)
	}
	r.Tag += fmt.Sprintf("_1y_%d", r.resultID)
	return nil
}

func (r *RollingReturn) LoadData() error {
	if r == nil || r.resultID <= 0 {
		return utils.Errorf(nil, "param error %+v", r)
	}
	var (
		err error
	)
	r.resultData, err = stock.GetStockStrategyResultData(&stock.GetStockStrategyResultDataRequest{
		StrategyResultID: r.resultID,
	})
	if err != nil || r.resultData == nil {
		return utils.Errorf(err, "stock.GetStockStrategyResultData fail")
	}
	return nil
}

func (r *RollingReturn) Step() (bool, error) {
	if r == nil || r.resultData == nil || r.stepIndex < 0 {
		return false, utils.Errorf(nil, "param error %+v", r)
	}
	if r.Result == nil {
		r.Result = &StrategyResult{
			LineData: []*PointData{},
		}
	}
	if int64(len(r.resultData.StockStrategyDataList)) < r.stepIndex+1 {
		return false, nil
	}
	point := &PointData{
		Time: r.resultData.StockStrategyDataList[r.stepIndex].TimeCST,
	}
	// 均值
	for i := r.stepIndex; i >= 0; i-- {
		if !r.resultData.StockStrategyDataList[r.stepIndex].TimeCST.AddDate(-1, 0, 0).Before(r.resultData.StockStrategyDataList[i].TimeCST) {
			curValueStr := r.resultData.StockStrategyDataList[r.stepIndex].Value
			curValue, err := strconv.ParseFloat(curValueStr, 64)
			if err != nil {
				return false, utils.Errorf(err, "strconv.ParseFloat fail")
			}
			oldValueStr := r.resultData.StockStrategyDataList[i].Value
			oldValue, err := strconv.ParseFloat(oldValueStr, 64)
			if err != nil {
				return false, utils.Errorf(err, "strconv.ParseFloat fail")
			}
			point.Value = (curValue - oldValue) / oldValue
			break
		}
	}
	r.Result.LineData = append(r.Result.LineData, point)
	r.stepIndex++
	return true, nil
}
