package stock

import (
	"fmt"
	"stock-go/common"
	"stock-go/data/mysql"
	"stock-go/utils"
	"time"
)

type GetAllStockCodeResponse struct {
	StockCodeList []*common.StockCode
}

func GetAllStockCode() (*GetAllStockCodeResponse, error) {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return nil, utils.Errorf(err, "mysql.GetConnection fail")
	}
	str := `
select
	id,
	ifnull(code,''),
	ifnull(name,''),
	ifnull(industry,''),
	ifnull(industry_classification,'')
from
	stock_all_code
	`
	rows, err := db.Query(str)
	if err != nil {
		return nil, utils.Errorf(err, "db.Query fail")
	}
	defer rows.Close()
	response := &GetAllStockCodeResponse{
		StockCodeList: []*common.StockCode{},
	}
	for rows.Next() {
		tmp := &common.StockCode{}
		err = rows.Scan(
			&tmp.ID,
			&tmp.Code,
			&tmp.Name,
			&tmp.Industry,
			&tmp.IndustryClassification,
		)
		if err != nil {
			return nil, utils.Errorf(err, "rows.Scan fail")
		}
		response.StockCodeList = append(response.StockCodeList, tmp)
	}
	return response, nil
}

type InsertStockCodeRequest struct {
	StockCodeList []*common.StockCode
}

func InsertStockCode(request *InsertStockCodeRequest) error {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || len(request.StockCodeList) <= 0 {
		return utils.Errorf(nil, "param error %+v", request)
	}
	str := `
insert into
	stock_all_code
(code,name,industry,industry_classification)
values
	`
	paramList := []interface{}{}
	for idx, code := range request.StockCodeList {
		if idx == len(request.StockCodeList)-1 {
			str += "(?,?,?,?)"
		} else {
			str += "(?,?,?,?),"
		}
		paramList = append(paramList, code.Code, code.Name, code.Industry, code.IndustryClassification)
	}
	str += `on duplicate key update update_time_utc = CURRENT_TIMESTAMP`
	_, err = db.Exec(str, paramList...)
	if err != nil {
		return utils.Errorf(err, "db.Exec fail")
	}
	return nil
}

type GetStockTradeDateResponse struct {
	StockTradeDateList []*common.StockTradeDate
}

func GetAllStockTradeDate() (*GetStockTradeDateResponse, error) {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return nil, utils.Errorf(err, "mysql.GetConnection fail")
	}
	str := `
select
	id,
	ifnull(date_cst,''),
	ifnull(is_trading_day,'')
from
	stock_trade_date
	`
	rows, err := db.Query(str)
	if err != nil {
		return nil, utils.Errorf(err, "db.Query fail")
	}
	defer rows.Close()
	result := &GetStockTradeDateResponse{
		StockTradeDateList: []*common.StockTradeDate{},
	}
	for rows.Next() {
		var (
			tmp                      = &common.StockTradeDate{}
			dateStr, isTradingDayStr string
		)
		err = rows.Scan(
			&tmp.ID,
			&dateStr,
			&isTradingDayStr,
		)
		if err != nil {
			return nil, utils.Errorf(err, "rows.Scan fail")
		}
		tmp.DateCST, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
		if isTradingDayStr == "yes" {
			tmp.IsTradingDay = true
		}
		result.StockTradeDateList = append(result.StockTradeDateList, tmp)
	}
	return result, nil
}

type InsertStockTradeDateRequest struct {
	StockTradeDateList []*common.StockTradeDate
}

func InsertStockTradeDate(request *InsertStockTradeDateRequest) error {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || len(request.StockTradeDateList) <= 0 {
		return utils.Errorf(nil, "param error %+v", request)
	}
	str := `
insert into
	stock_trade_date
(date_cst,is_trading_day)
values
	`
	paramList := []interface{}{}
	for idx, tradeDate := range request.StockTradeDateList {
		if idx == len(request.StockTradeDateList)-1 {
			str += "(?,?)"
		} else {
			str += "(?,?),"
		}
		var (
			isTradingDayStr string = "no"
		)
		if tradeDate.IsTradingDay {
			isTradingDayStr = "yes"
		}
		paramList = append(paramList, tradeDate.DateCST.Format("2006-01-02 15:04:05"), isTradingDayStr)
	}
	str += `on duplicate key update update_time_utc = CURRENT_TIMESTAMP`
	_, err = db.Exec(str, paramList...)
	if err != nil {
		return utils.Errorf(err, "db.Exec fail")
	}
	return nil
}

type GetStockKDataCountRequest struct {
	Code       string
	StartTime  time.Time
	EndTime    time.Time
	Frequency  string
	AdjustFlag string
}

type GetStockKDataCountResponse struct {
	Count int64
}

func GetStockKDataCount(request *GetStockKDataCountRequest) (*GetStockKDataCountResponse, error) {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return nil, utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || request.Code == "" || request.StartTime.Unix() <= 0 || request.EndTime.Unix() <= 0 || request.Frequency == "" || request.AdjustFlag == "" {
		return nil, utils.Errorf(nil, "param error %+v", request)
	}
	response := &GetStockKDataCountResponse{}
	for i := request.StartTime.Year(); i <= request.EndTime.Year(); i++ {
		yResponse, err := GetStockKDataCountYear(request, int64(i))
		if err != nil {
			return nil, utils.Errorf(err, "GetStockKDataCountYear fail")
		}
		response.Count += yResponse.Count
	}
	return response, nil
}

func GetStockKDataCountYear(request *GetStockKDataCountRequest, year int64) (*GetStockKDataCountResponse, error) {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return nil, utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || request.Code == "" || request.StartTime.Unix() <= 0 || request.EndTime.Unix() <= 0 || request.Frequency == "" || request.AdjustFlag == "" {
		return nil, utils.Errorf(nil, "param error %+v", request)
	}
	str := fmt.Sprintf(`
select
	count(id)
from
	stock_k_data_%d
where
	code = ? and time_cst >= ? and time_cst < ? and frequency = ? and adjust_flag = ?
	`, year)
	response := &GetStockKDataCountResponse{}
	err = db.QueryRow(str,
		request.Code,
		request.StartTime.Format("2006-01-02 15:04:05"),
		request.EndTime.Format("2006-01-02 15:04:05"),
		request.Frequency,
		request.AdjustFlag,
	).Scan(&response.Count)
	if err != nil {
		return nil, utils.Errorf(err, "db.QueryRow fail")
	}
	return response, nil
}

type GetStockKDataRequest struct {
	Code       string
	StartTime  time.Time
	EndTime    time.Time
	Frequency  string
	AdjustFlag string
}

type GetStockKDataResponse struct {
	StockKDataList []*common.StockKData
}

func GetStockKData(request *GetStockKDataRequest) (*GetStockKDataResponse, error) {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return nil, utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || request.Code == "" || request.StartTime.Unix() <= 0 || request.EndTime.Unix() <= 0 || request.Frequency == "" || request.AdjustFlag == "" {
		return nil, utils.Errorf(nil, "param error %+v", request)
	}
	response := &GetStockKDataResponse{
		StockKDataList: []*common.StockKData{},
	}
	for i := request.StartTime.Year(); i <= request.EndTime.Year(); i++ {
		yResponse, err := GetStockKDataYear(request, int64(i))
		if err != nil {
			return nil, utils.Errorf(err, "GetStockKDataYear fail")
		}
		response.StockKDataList = append(response.StockKDataList, yResponse.StockKDataList...)
	}
	return response, nil
}

func GetStockKDataYear(request *GetStockKDataRequest, year int64) (*GetStockKDataResponse, error) {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return nil, utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || request.Code == "" || request.StartTime.Unix() <= 0 || request.EndTime.Unix() <= 0 || request.Frequency == "" || request.AdjustFlag == "" {
		return nil, utils.Errorf(nil, "param error %+v", request)
	}
	str := fmt.Sprintf(`
select
	ifnull(code,''),
	ifnull(time_cst,''),
	ifnull(frequency,''),
	ifnull(adjust_flag,''),
	ifnull(open,''),
	ifnull(high,''),
	ifnull(low,''),
	ifnull(close,''),
	ifnull(preclose,''),
	ifnull(volume,''),
	ifnull(amount,''),
	ifnull(turn,''),
	ifnull(trade_status,''),
	ifnull(pct_chg,''),
	ifnull(pe_ttm,''),
	ifnull(pb_mrq,''),
	ifnull(ps_ttm,''),
	ifnull(pcf_ncf_ttm,''),
	ifnull(is_st,'')
from
	stock_k_data_%d
where
	code = ? and time_cst >= ? and time_cst < ? and frequency = ? and adjust_flag = ?
	`, year)
	rows, err := db.Query(str,
		request.Code,
		request.StartTime.Format("2006-01-02 15:04:05"),
		request.EndTime.Format("2006-01-02 15:04:05"),
		request.Frequency,
		request.AdjustFlag,
	)
	if err != nil {
		return nil, utils.Errorf(err, "db.Query fail")
	}
	defer rows.Close()
	response := &GetStockKDataResponse{
		StockKDataList: []*common.StockKData{},
	}
	for rows.Next() {
		var (
			timeStr, isSTStr string
		)
		tmp := &common.StockKData{}
		err = rows.Scan(
			&tmp.Code,
			&timeStr,
			&tmp.Frequency,
			&tmp.AdjustFlag,
			&tmp.Open,
			&tmp.High,
			&tmp.Low,
			&tmp.Close,
			&tmp.Preclose,
			&tmp.Volume,
			&tmp.Amount,
			&tmp.Turn,
			&tmp.TradeStatus,
			&tmp.PctChg,
			&tmp.PeTTM,
			&tmp.PbMRQ,
			&tmp.PsTTM,
			&tmp.PcfNcfTTM,
			&isSTStr,
		)
		if err != nil {
			return nil, utils.Errorf(err, "rows.Scan fail")
		}
		if isSTStr == "yes" {
			tmp.IsST = true
		}
		tmp.TimeCST, err = time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
		response.StockKDataList = append(response.StockKDataList, tmp)
	}
	return response, nil
}

type InsertStockKDataRequest struct {
	StockKDataList []*common.StockKData
}

func InsertStockKData(request *InsertStockKDataRequest) error {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || len(request.StockKDataList) <= 0 {
		return utils.Errorf(nil, "param error %+v", request)
	}
	dataMap := map[int64][]*common.StockKData{}
	for _, data := range request.StockKDataList {
		if dataMap[int64(data.TimeCST.Year())] == nil {
			dataMap[int64(data.TimeCST.Year())] = []*common.StockKData{}
		}
		dataMap[int64(data.TimeCST.Year())] = append(dataMap[int64(data.TimeCST.Year())], data)
	}
	for year, dataList := range dataMap {
		err := InsertStockKDataYear(&InsertStockKDataRequest{
			StockKDataList: dataList,
		}, year)
		if err != nil {
			return utils.Errorf(err, "InsertStockKDataYear fail")
		}
	}
	return nil
}

func InsertStockKDataYear(request *InsertStockKDataRequest, year int64) error {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || len(request.StockKDataList) <= 0 {
		return utils.Errorf(nil, "param error %+v", request)
	}
	str := fmt.Sprintf(`
insert into
	stock_k_data_%d
(
	code,
	time_cst,
	frequency,
	adjust_flag,
	open,
	high,
	low,
	close,
	preclose,
	volume,
	amount,
	turn,
	trade_status,
	pct_chg,
	pe_ttm,
	pb_mrq,
	ps_ttm,
	pcf_ncf_ttm,
	is_st
)
values
	`, year)
	paramList := []interface{}{}
	for idx, kData := range request.StockKDataList {
		if idx == len(request.StockKDataList)-1 {
			str += "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		} else {
			str += "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?),"
		}
		var (
			isST string = "no"
		)
		if kData.IsST {
			isST = "yes"
		}
		paramList = append(paramList,
			kData.Code,
			kData.TimeCST.Format("2006-01-02 15:04:05"),
			kData.Frequency,
			kData.AdjustFlag,
			kData.Open,
			kData.High,
			kData.Low,
			kData.Close,
			kData.Preclose,
			kData.Volume,
			kData.Amount,
			kData.Turn,
			kData.TradeStatus,
			kData.PctChg,
			kData.PeTTM,
			kData.PbMRQ,
			kData.PsTTM,
			kData.PcfNcfTTM,
			isST,
		)
	}
	str += `on duplicate key update update_time_utc = CURRENT_TIMESTAMP`
	_, err = db.Exec(str, paramList...)
	if err != nil {
		return utils.Errorf(err, "db.Exec fail")
	}
	return nil
}

type InsertStockStrategyDataRequest struct {
	StockStrategyDataList []*common.StockStrategyData
}

func InsertStockStrategyData(request *InsertStockStrategyDataRequest) error {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || len(request.StockStrategyDataList) <= 0 {
		return utils.Errorf(nil, "param error %+v", request)
	}
	str := `
insert into
	stock_strategy_result_data
(
	stock_strategy_result_id,
	code,
	tag,
	time_cst,
	value
)
values
	`
	paramList := []interface{}{}
	for idx, data := range request.StockStrategyDataList {
		if idx == len(request.StockStrategyDataList)-1 {
			str += "(?,?,?,?,?)"
		} else {
			str += "(?,?,?,?,?),"
		}
		paramList = append(paramList,
			data.StockStrategyResultID,
			data.Code,
			data.Tag,
			data.TimeCST.Format("2006-01-02 15:04:05"),
			data.Value,
		)
	}
	_, err = db.Exec(str, paramList...)
	if err != nil {
		return utils.Errorf(err, "db.Exec fail")
	}
	return nil
}

type InsertStockStrategyResultRequest struct {
	StockStrategyResult *common.StockStrategyResult
}

type InsertStockStrategyResultResponse struct {
	StockStrategyResult *common.StockStrategyResult
}

func InsertStockStrategyResult(request *InsertStockStrategyResultRequest) (*InsertStockStrategyResultResponse, error) {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return nil, utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || request.StockStrategyResult == nil {
		return nil, utils.Errorf(nil, "param error %+v", request)
	}
	str := `
insert into
	stock_strategy_result
(
	code,
	tag,
	start_time_cst,
	end_time_cst,
	anual_return_rate,
	draw_down
)
values
	(?,?,?,?,?,?)
	`
	paramList := []interface{}{}
	paramList = append(paramList,
		request.StockStrategyResult.Code,
		request.StockStrategyResult.Tag,
		request.StockStrategyResult.StartTimeCST.Format("2006-01-02 15:04:05"),
		request.StockStrategyResult.EndTimeCST.Format("2006-01-02 15:04:05"),
		request.StockStrategyResult.AnualReturnRate,
		request.StockStrategyResult.DrawDown,
	)
	result, err := db.Exec(str, paramList...)
	if err != nil {
		return nil, utils.Errorf(err, "db.Exec fail")
	}
	request.StockStrategyResult.ID, err = result.LastInsertId()
	if err != nil {
		return nil, utils.Errorf(err, "result.LastInsertId fail")
	}
	return &InsertStockStrategyResultResponse{
		StockStrategyResult: request.StockStrategyResult,
	}, nil
}

type GetStockStrategyResultDataRequest struct {
	StrategyResultID int64
}

type GetStockStrategyResultDataResponse struct {
	StockStrategyDataList []*common.StockStrategyData
}

func GetStockStrategyResultData(request *GetStockStrategyResultDataRequest) (*GetStockStrategyResultDataResponse, error) {
	db, err := mysql.GetConnection()
	if err != nil || db == nil {
		return nil, utils.Errorf(err, "mysql.GetConnection fail")
	}
	if request == nil || request.StrategyResultID <= 0 {
		return nil, utils.Errorf(nil, "param error %+v", request)
	}
	str := `
select
	ifnull(id,''),
	ifnull(stock_strategy_result_id,0),
	ifnull(code,''),
	ifnull(tag,''),
	ifnull(time_cst,''),
	ifnull(value,'')
from
	stock_strategy_result_data
where
	stock_strategy_result_id = ?
order by
	time_cst
	`
	rows, err := db.Query(str, request.StrategyResultID)
	if err != nil {
		return nil, utils.Errorf(err, "db.Query fail")
	}
	defer rows.Close()
	response := &GetStockStrategyResultDataResponse{
		StockStrategyDataList: []*common.StockStrategyData{},
	}
	for rows.Next() {
		tmp := &common.StockStrategyData{}
		var timeStr string
		err = rows.Scan(
			&tmp.ID,
			&tmp.StockStrategyResultID,
			&tmp.Code,
			&tmp.Tag,
			&timeStr,
			&tmp.Value,
		)
		if err != nil {
			return nil, utils.Errorf(err, "rows.Scan fail")
		}
		tmp.TimeCST, err = time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
		response.StockStrategyDataList = append(response.StockStrategyDataList, tmp)
	}
	return response, nil
}
