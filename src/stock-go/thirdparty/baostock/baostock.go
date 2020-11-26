package baostock

import (
	"bytes"
	"compress/zlib"
	"context"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net"
	"runtime/debug"
	"stock-go/utils"
	"strconv"
	"strings"
	"time"
)

var (
	baostockConnection *BaostockConnection
)

type BaostockConnection struct {
	connection net.Conn
	user       string
}

func GetBaostockConnection() (*BaostockConnection, error) {
	if baostockConnection != nil {
		return baostockConnection, nil
	}
	tmpConnection, err := newBaostockConnection()
	if err != nil {
		return nil, utils.Errorf(err, "NewBaostockConnection fail")
	}
	err = tmpConnection.login("", "", 0)
	if err != nil {
		return nil, utils.Errorf(err, "tmpConnection.login fail")
	}
	baostockConnection = tmpConnection
	return baostockConnection, nil
}

func CloseBaostockConnection() error {
	if baostockConnection == nil {
		return nil
	}
	tmpConnection := baostockConnection
	baostockConnection = nil
	err := tmpConnection.logout()
	if err != nil {
		utils.LogErr(fmt.Sprintf("bc.logout fail %s", err))
	}
	err = tmpConnection.closeConnection()
	if err != nil {
		utils.LogErr(fmt.Sprintf("bc.CloseConnection fail %s", err))
	}
	return nil
}

func CloseBaostockConnectionWithoutLogout() error {
	if baostockConnection == nil {
		return nil
	}
	tmpConnection := baostockConnection
	baostockConnection = nil
	err := tmpConnection.closeConnection()
	if err != nil {
		utils.LogErr(fmt.Sprintf("bc.CloseConnection fail %s", err))
	}
	return nil
}

func ReconnectBaostock() (*BaostockConnection, error) {
	utils.Log("reconnectiong to baostock ...")
	err := CloseBaostockConnectionWithoutLogout()
	if err != nil {
		return nil, utils.Errorf(err, "CloseBaostockConnectionWithoutLogout fail")
	}
	bc, err := GetBaostockConnection()
	if err != nil {
		return nil, utils.Errorf(err, "GetBaostockConnection fail")
	}
	utils.Log("reconnected to baostock")
	return bc, nil
}

func newBaostockConnection() (*BaostockConnection, error) {
	bc := &BaostockConnection{}
	err := bc.connect()
	if err != nil {
		return nil, utils.Errorf(err, "bc.Connect fail")
	}
	return bc, nil
}

func (bc *BaostockConnection) connect() error {
	if bc == nil {
		return utils.Errorf(nil, "bc is nil")
	}
	if bc.connection != nil {
		return nil
	}
	utils.Log("connecting to baostock ...")
	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", BAOSTOCK_SERVER_IP, BAOSTOCK_SERVER_PORT))
	if err != nil {
		return utils.Errorf(err, "net.Dial fail")
	}
	bc.connection = connection
	utils.Log("connected to baostock")
	return nil
}

func (bc *BaostockConnection) closeConnection() error {
	if bc == nil || bc.connection == nil {
		return utils.Errorf(nil, "bc is nil or connection is nil")
	}
	utils.Log("closing baostock connection ...")
	err := bc.connection.Close()
	if err != nil {
		return utils.Errorf(err, "Close fail")
	}
	bc.connection = nil
	utils.Log("baostock connection closed")
	return nil
}

func (bc *BaostockConnection) login(user, password string, options int64) error {
	if bc == nil || bc.connection == nil {
		return utils.Errorf(nil, "bc is nil or connection is nil")
	}
	if user == "" {
		user = "anonymous"
	}
	if password == "" {
		password = "123456"
	}
	utils.Log("login to baostock ...")
	msgBody := "login" + MESSAGE_SPLIT +
		user + MESSAGE_SPLIT +
		password + MESSAGE_SPLIT +
		strconv.FormatInt(options, 10)
	msgHeader := messageHeader(MESSAGE_TYPE_LOGIN_REQUEST, int64(len(msgBody)))
	utils.Log("login to baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return utils.Errorf(err, "bc.Write fail")
	}
	recieveBody, err := bc.read()
	if err != nil {
		return utils.Errorf(err, "bc.Read fail")
	}
	recieveData, err := bc.decodeRecieve(recieveBody)
	if err != nil {
		return utils.Errorf(err, "bc.DecodeRecieve fail")
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return utils.Errorf(nil, "return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetLoginResponse()
	if err != nil {
		return utils.Errorf(err, "recieveData.GetLoginResponse fail")
	}
	bc.user = response.User
	utils.Log("login to baostock success")
	return nil
}

func (bc *BaostockConnection) logout() error {
	if bc == nil || bc.connection == nil {
		return utils.Errorf(nil, "bc is nil or connection is nil")
	}
	utils.Log("logout from baostock ...")
	if bc.user == "" {
		return utils.Errorf(nil, "bc.user is nil, not login")
	}
	msgBody := "logout" + MESSAGE_SPLIT +
		bc.user + MESSAGE_SPLIT +
		time.Now().Format("20060102150405")
	msgHeader := messageHeader(MESSAGE_TYPE_LOGOUT_REQUEST, int64(len(msgBody)))
	utils.Log("logout from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return utils.Errorf(err, "bc.Write fail")
	}
	recieveBody, err := bc.read()
	if err != nil {
		return utils.Errorf(err, "bc.Read fail")
	}
	recieveData, err := bc.decodeRecieve(recieveBody)
	if err != nil {
		return utils.Errorf(err, "bc.DecodeRecieve fail")
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return utils.Errorf(nil, "return error code is not 0 %+v", recieveData)
	}
	utils.Log("logout to baostock success")
	return nil
}

func (bc *BaostockConnection) QueryHistoryKDataPlusWithTimeOut(code, fields, startDate, endDate, frequency, adjustFlag string, timeoutSecond int64) (*QueryHistoryKDataResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
	var (
		ret       *QueryHistoryKDataResponse
		err       error
		startTime = time.Now()
	)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				utils.Log(fmt.Sprintf("panic:\n\t%s", err))
				stack := strings.Join(strings.Split(string(debug.Stack()), "\n")[2:], "\n")
				utils.Log(fmt.Sprintf("stack:\n\t%s", stack))
			}
		}()
		ret, err = bc.QueryHistoryKDataPlus(code, fields, startDate, endDate, frequency, adjustFlag)
		cancel()
	}()
	select {
	case <-ctx.Done():
		switch ctxErr := ctx.Err(); ctxErr {
		case context.DeadlineExceeded:
			utils.LogErr(fmt.Sprintf("bc.QueryHistoryKDataPlus timeout cost %d err %s", time.Since(startTime), err))
			err = QueryTimeoutErr
		default:
			// 没有超时
		}
	}
	return ret, err
}

func (bc *BaostockConnection) QueryHistoryKDataPlus(code, fields, startDate, endDate, frequency, adjustFlag string) (*QueryHistoryKDataResponse, error) {
	return bc.QueryHistoryKDataPage(code, fields, startDate, endDate, frequency, adjustFlag, 1, BAOSTOCK_PER_PAGE_COUNT)
}

func (bc *BaostockConnection) QueryHistoryKDataPage(code, fields, startDate, endDate, frequency, adjustFlag string, curPageNum, perPageCount int64) (*QueryHistoryKDataResponse, error) {
	if bc == nil || bc.connection == nil {
		return nil, utils.Errorf(nil, "bc is nil or connection is nil")
	}
	if code == "" || len(code) != STOCK_CODE_LENGTH {
		return nil, utils.Errorf(nil, "code len error %s", code)
	}
	code = strings.ToLower(code)
	if strings.HasSuffix(code, "sh") || strings.HasSuffix(code, "sz") {
		code = code[7:9] + "." + code[0:6]
	}
	if fields == "" {
		return nil, utils.Errorf(nil, "fields error %s", fields)
	}
	if startDate == "" {
		startDate = DEFAULT_START_DATE
	} else {
		_, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
	}
	if frequency == "" {
		frequency = "d"
	}
	if adjustFlag == "" {
		adjustFlag = "3"
	}
	utils.Log("querying history k data from baostock ...")
	if bc.user == "" {
		return nil, utils.Errorf(nil, "bc.user is nil, not login")
	}
	msgBody := "query_history_k_data" + MESSAGE_SPLIT + bc.user + MESSAGE_SPLIT +
		strconv.FormatInt(curPageNum, 10) + MESSAGE_SPLIT + strconv.FormatInt(perPageCount, 10) + MESSAGE_SPLIT +
		code + MESSAGE_SPLIT + fields + MESSAGE_SPLIT + startDate + MESSAGE_SPLIT + endDate + MESSAGE_SPLIT + frequency + MESSAGE_SPLIT + adjustFlag
	msgHeader := messageHeader(MESSAGE_TYPE_GETKDATA_REQUEST, int64(len(msgBody)))
	utils.Log("query history k data from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return nil, utils.Errorf(err, "bc.Write fail")
	}
	recieveBody, err := bc.read()
	if err != nil {
		return nil, utils.Errorf(err, "bc.Read fail")
	}
	recieveData, err := bc.decodeRecieve(recieveBody)
	if err != nil {
		return nil, utils.Errorf(err, "bc.DecodeRecieve fail")
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return nil, utils.Errorf(nil, "return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetQueryHistoryKDataResponse()
	if err != nil {
		return nil, utils.Errorf(err, "recieveData.GetQueryHistoryKDataResponse fail")
	}
	utils.Log("query history k data baostock success")
	return response, nil
}

func (bc *BaostockConnection) QueryTradeDates(startDate, endDate string) (*QueryTradeDatesResponse, error) {
	if bc == nil || bc.connection == nil {
		return nil, utils.Errorf(nil, "bc is nil or connection is nil")
	}
	if startDate == "" {
		startDate = DEFAULT_START_DATE
	} else {
		_, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
	}
	utils.Log("querying trade dates from baostock ...")
	if bc.user == "" {
		return nil, utils.Errorf(nil, "bc.user is nil, not login")
	}
	msgBody := "query_trade_dates" + MESSAGE_SPLIT + bc.user + MESSAGE_SPLIT +
		"1" + MESSAGE_SPLIT + strconv.FormatInt(BAOSTOCK_PER_PAGE_COUNT, 10) + MESSAGE_SPLIT +
		startDate + MESSAGE_SPLIT + endDate + MESSAGE_SPLIT
	msgHeader := messageHeader(MESSAGE_TYPE_QUERYTRADEDATES_REQUEST, int64(len(msgBody)))
	utils.Log("query trade dates from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return nil, utils.Errorf(err, "bc.Write fail")
	}
	recieveBody, err := bc.read()
	if err != nil {
		return nil, utils.Errorf(err, " bc.Read fail")
	}
	recieveData, err := bc.decodeRecieve(recieveBody)
	if err != nil {
		return nil, utils.Errorf(err, "bc.DecodeRecieve fail")
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return nil, utils.Errorf(nil, "return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetQueryTradeDatesResponse()
	if err != nil {
		return nil, utils.Errorf(err, "recieveData.GetQueryTradeDatesResponse fail")
	}
	utils.Log("query trade dates from baostock success")
	return response, nil
}

func (bc *BaostockConnection) QueryAllStock(date string) (*QueryAllStockResponse, error) {
	if bc == nil || bc.connection == nil {
		return nil, utils.Errorf(nil, "bc is nil or connection is nil")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
	}
	utils.Log("querying all stock from baostock ...")
	if bc.user == "" {
		return nil, utils.Errorf(nil, "bc.user is nil, not login")
	}
	msgBody := "query_all_stock" + MESSAGE_SPLIT + bc.user + MESSAGE_SPLIT +
		"1" + MESSAGE_SPLIT + strconv.FormatInt(BAOSTOCK_PER_PAGE_COUNT, 10) + MESSAGE_SPLIT +
		date + MESSAGE_SPLIT
	msgHeader := messageHeader(MESSAGE_TYPE_QUERYALLSTOCK_REQUEST, int64(len(msgBody)))
	utils.Log("query all stock from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return nil, utils.Errorf(err, "bc.Write fail")
	}
	recieveBody, err := bc.read()
	if err != nil {
		return nil, utils.Errorf(err, "bc.Read fail")
	}
	recieveData, err := bc.decodeRecieve(recieveBody)
	if err != nil {
		return nil, utils.Errorf(err, "bc.DecodeRecieve fail")
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return nil, utils.Errorf(nil, "return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetQueryAllStockResponse()
	if err != nil {
		return nil, utils.Errorf(err, "recieveData.GetQueryAllStockResponse fail")
	}
	utils.Log("query all stock from baostock success")
	return response, nil
}

func (bc *BaostockConnection) QueryStockIndustry(code, date string) (*QueryStockIndustryResponse, error) {
	if bc == nil || bc.connection == nil {
		return nil, utils.Errorf(nil, "bc is nil or connection is nil")
	}

	if code == "" {
		return nil, utils.Errorf(nil, "code is nil")
	}
	if date != "" {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, utils.Errorf(err, "time.Parse fail")
		}
	}
	utils.Log("querying stock industry from baostock ...")
	if bc.user == "" {
		return nil, utils.Errorf(nil, "bc.user is nil, not login")
	}
	msgBody := "query_stock_industry" + MESSAGE_SPLIT + bc.user + MESSAGE_SPLIT +
		"1" + MESSAGE_SPLIT + strconv.FormatInt(BAOSTOCK_PER_PAGE_COUNT, 10) + MESSAGE_SPLIT +
		code + MESSAGE_SPLIT + date
	msgHeader := messageHeader(MESSAGE_TYPE_QUERYSTOCKINDUSTRY_REQUEST, int64(len(msgBody)))
	utils.Log("[QueryStockIndustry] query stock industry from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return nil, utils.Errorf(err, "bc.Write fail")
	}
	recieveBody, err := bc.read()
	if err != nil {
		return nil, utils.Errorf(err, "bc.Read fail")
	}
	recieveData, err := bc.decodeRecieve(recieveBody)
	if err != nil {
		return nil, utils.Errorf(err, "bc.DecodeRecieve fail")
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return nil, utils.Errorf(nil, "return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetQueryStockIndustryResponse()
	if err != nil {
		return nil, utils.Errorf(err, "recieveData.GetQueryStockIndustryResponse fail")
	}
	utils.Log("query stock industry from baostock success")
	return response, nil
}

func (bc *BaostockConnection) write(data []byte) (int64, error) {
	if bc == nil || bc.connection == nil {
		return 0, utils.Errorf(nil, "bc is nil or connection is nil")
	}
	if data == nil || len(data) <= 0 {
		return 0, utils.Errorf(nil, "data is nil")
	}
	data = append(data, []byte("\n")...)
	n, err := bc.connection.Write(data)
	return int64(n), err
}

func (bc *BaostockConnection) read() ([]byte, error) {
	if bc == nil || bc.connection == nil {
		return nil, utils.Errorf(nil, "bc is nil or connection is nil")
	}
	data := []byte{}
	buf := make([]byte, 8192)
	for {
		n, err := bc.connection.Read(buf)
		if err != nil {
			return nil, utils.Errorf(err, "bc.connection.Read fail")
		}
		// utils.Log("[Read] recive data " + string(buf[:n]))
		data = append(data, buf[:n]...)
		if len(data) >= 13 && bytes.Compare(data[len(data)-13:], []byte("<![CDATA[]]>\n")) == 0 {
			break
		}
	}
	return data, nil
}

func (bc *BaostockConnection) decodeRecieve(data []byte) (*RecieveData, error) {
	if bc == nil || bc.connection == nil {
		return nil, utils.Errorf(nil, "bc is nil or connection is nil")
	}
	if data == nil || len(data) <= 0 {
		return nil, utils.Errorf(nil, "data is nil")
	}
	if len(data) < MESSAGE_HEADER_LENGTH {
		return nil, utils.Errorf(nil, "data is error %s", string(data))
	}
	ret := &RecieveData{}
	headStr := string(data[:MESSAGE_HEADER_LENGTH])
	// utils.Log("[DecodeRecieve] recive data head " + string(headStr))
	headArr := strings.Split(headStr, MESSAGE_SPLIT)
	// utils.Log(fmt.Sprint("[DecodeRecieve] recive data head ", headArr))
	if len(headArr) < 3 {
		return nil, utils.Errorf(nil, "data is error %s", string(data))
	}
	ret.MessageType = headArr[1]
	headInnerLength, err := strconv.ParseInt(headArr[2], 10, 64)
	if err != nil {
		return nil, utils.Errorf(err, "strconv.ParseInt fail")
	}
	ret.MessageBodyLength = headInnerLength
	bodyStr := ""
	if utils.StringIsIn(headArr[1], COMPRESSED_MESSAGE_TYPE_TUPLE...) {
		body, err := readSegment(data[MESSAGE_HEADER_LENGTH : MESSAGE_HEADER_LENGTH+headInnerLength])
		if err != nil {
			return nil, utils.Errorf(err, "readSegment fail")
		}
		bodyStr = string(body)
	} else {
		bodyStr = string(data[MESSAGE_HEADER_LENGTH:])
	}
	bodyArr := strings.Split(bodyStr, MESSAGE_SPLIT)
	// utils.Log(fmt.Sprint("[DecodeRecieve] recive data body ", bodyArr))
	if len(bodyArr) < 2 {
		return nil, utils.Errorf(nil, "data is error %s", string(data))
	}
	ret.ErrorCode = bodyArr[0]
	ret.ErrorMessage = bodyArr[1]
	if len(bodyArr) > 2 {
		ret.DataList = []string{}
		for i := 2; i < len(bodyArr); i++ {
			ret.DataList = append(ret.DataList, bodyArr[i])
		}
	}
	err = ret.GetResponse()
	if err != nil {
		return nil, utils.Errorf(err, "ret.GetResponse fail")
	}
	return ret, nil
}

func messageHeader(msgType string, msgLength int64) string {
	return BAOSTOCK_CLIENT_VERSION + MESSAGE_SPLIT +
		msgType + MESSAGE_SPLIT +
		utils.AddZeroForString(msgLength, 10, true)
}

func readSegment(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	p, err := ioutil.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return p, nil
}
