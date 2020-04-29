package baostock

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net"
	"stock-go/utils"
	"strconv"
	"strings"
	"time"
)

type BaostockConnection struct {
	connection net.Conn
	user       string
}

func NewBaostockConnection() (*BaostockConnection, error) {
	utils.Log("[NewBaostockConnection] connecting to baostock ...")
	bc := &BaostockConnection{}
	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", BAOSTOCK_SERVER_IP, BAOSTOCK_SERVER_PORT))
	if err != nil {
		return nil, fmt.Errorf("[NewBaostockConnection] net.Dial fail\n\t%s", err)
	}
	bc.connection = connection
	utils.Log("[NewBaostockConnection] connected to baostock")
	return bc, nil
}

func (bc *BaostockConnection) CloseConnection() error {
	if bc == nil || bc.connection == nil {
		return fmt.Errorf("[CloseConnection] bc is nil or connection is nil")
	}
	utils.Log("[CloseConnection] closing baostock connection ...")
	err := bc.connection.Close()
	if err != nil {
		return fmt.Errorf("[CloseConnection] Close fail\n\t%s", err)
	}
	utils.Log("[CloseConnection] baostock connection closed")
	return nil
}

func (bc *BaostockConnection) Login(user, password string, options int64) error {
	if bc == nil || bc.connection == nil {
		return fmt.Errorf("[Login] bc is nil or connection is nil")
	}
	if user == "" {
		user = "anonymous"
	}
	if password == "" {
		password = "123456"
	}
	utils.Log("[Login] login to baostock ...")
	msgBody := "login" + MESSAGE_SPLIT +
		user + MESSAGE_SPLIT +
		password + MESSAGE_SPLIT +
		strconv.FormatInt(options, 10)
	msgHeader := messageHeader(MESSAGE_TYPE_LOGIN_REQUEST, int64(len(msgBody)))
	utils.Log("[Login] login to baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.Write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return fmt.Errorf("[Login] bc.Write fail\n\t%s", err)
	}
	recieveBody, err := bc.Read()
	if err != nil {
		return fmt.Errorf("[Login] bc.Read fail\n\t%s", err)
	}
	recieveData, err := bc.DecodeRecieve(recieveBody)
	if err != nil {
		return fmt.Errorf("[Login] bc.DecodeRecieve fail\n\t%s", err)
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return fmt.Errorf("[Login] return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetLoginResponse()
	if err != nil {
		return fmt.Errorf("[Login] recieveData.GetLoginResponse fail\n\t%s", err)
	}
	bc.user = response.User
	utils.Log("[Login] login to baostock success")
	return nil
}

func (bc *BaostockConnection) Logout() error {
	if bc == nil || bc.connection == nil {
		return fmt.Errorf("[Logout] bc is nil or connection is nil")
	}
	utils.Log("[Logout] logout from baostock ...")
	if bc.user == "" {
		return fmt.Errorf("[Logout] bc.user is nil, not login")
	}
	msgBody := "logout" + MESSAGE_SPLIT +
		bc.user + MESSAGE_SPLIT +
		time.Now().Format("20060102150405")
	msgHeader := messageHeader(MESSAGE_TYPE_LOGOUT_REQUEST, int64(len(msgBody)))
	utils.Log("[Logout] logout from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.Write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return fmt.Errorf("[Logout] bc.Write fail\n\t%s", err)
	}
	recieveBody, err := bc.Read()
	if err != nil {
		return fmt.Errorf("[Logout] bc.Read fail\n\t%s", err)
	}
	recieveData, err := bc.DecodeRecieve(recieveBody)
	if err != nil {
		return fmt.Errorf("[Logout] bc.DecodeRecieve fail\n\t%s", err)
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return fmt.Errorf("[Logout] return error code is not 0 %+v", recieveData)
	}
	utils.Log("[Logout] logout to baostock success")
	return nil
}

func (bc *BaostockConnection) QueryHistoryKDataPlus(code, fields, startDate, endDate, frequency, adjustFlag string) (*QueryHistoryKDataResponse, error) {
	return bc.QueryHistoryKDataPage(code, fields, startDate, endDate, frequency, adjustFlag, 1, BAOSTOCK_PER_PAGE_COUNT)
}

func (bc *BaostockConnection) QueryHistoryKDataPage(code, fields, startDate, endDate, frequency, adjustFlag string, curPageNum, perPageCount int64) (*QueryHistoryKDataResponse, error) {
	if bc == nil || bc.connection == nil {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] bc is nil or connection is nil")
	}
	if code == "" || len(code) != STOCK_CODE_LENGTH {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] code len error %s", code)
	}
	code = strings.ToLower(code)
	if strings.HasSuffix(code, "sh") || strings.HasSuffix(code, "sz") {
		code = code[7:9] + "." + code[0:6]
	}
	if fields == "" {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] fields error %s", fields)
	}
	if startDate == "" {
		startDate = DEFAULT_START_DATE
	} else {
		_, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, fmt.Errorf("[QueryHistoryKDataPage] time.Parse fail %s", err)
		}
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, fmt.Errorf("[QueryHistoryKDataPage] time.Parse fail %s", err)
		}
	}
	if frequency == "" {
		frequency = "d"
	}
	if adjustFlag == "" {
		adjustFlag = "3"
	}
	utils.Log("[QueryHistoryKDataPage] querying history k data from baostock ...")
	if bc.user == "" {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] bc.user is nil, not login")
	}
	msgBody := "query_history_k_data" + MESSAGE_SPLIT + bc.user + MESSAGE_SPLIT +
		strconv.FormatInt(curPageNum, 10) + MESSAGE_SPLIT + strconv.FormatInt(perPageCount, 10) + MESSAGE_SPLIT +
		code + MESSAGE_SPLIT + fields + MESSAGE_SPLIT + startDate + MESSAGE_SPLIT + endDate + MESSAGE_SPLIT + frequency + MESSAGE_SPLIT + adjustFlag
	msgHeader := messageHeader(MESSAGE_TYPE_GETKDATA_REQUEST, int64(len(msgBody)))
	utils.Log("[QueryHistoryKDataPage] query history k data from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.Write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] bc.Write fail\n\t%s", err)
	}
	recieveBody, err := bc.Read()
	if err != nil {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] bc.Read fail\n\t%s", err)
	}
	recieveData, err := bc.DecodeRecieve(recieveBody)
	if err != nil {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] bc.DecodeRecieve fail\n\t%s", err)
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetQueryHistoryKDataResponse()
	if err != nil {
		return nil, fmt.Errorf("[QueryHistoryKDataPage] recieveData.GetQueryHistoryKDataResponse fail\n\t%s", err)
	}
	utils.Log("[QueryHistoryKDataPage] query history k data baostock success")
	return response, nil
}

func (bc *BaostockConnection) QueryTradeDates(startDate, endDate string) (*QueryTradeDatesResponse, error) {
	if bc == nil || bc.connection == nil {
		return nil, fmt.Errorf("[QueryTradeDates] bc is nil or connection is nil")
	}
	if startDate == "" {
		startDate = DEFAULT_START_DATE
	} else {
		_, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, fmt.Errorf("[QueryTradeDates] time.Parse fail %s", err)
		}
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, fmt.Errorf("[QueryTradeDates] time.Parse fail %s", err)
		}
	}
	utils.Log("[QueryTradeDates] querying trade dates from baostock ...")
	if bc.user == "" {
		return nil, fmt.Errorf("[QueryTradeDates] bc.user is nil, not login")
	}
	msgBody := "query_trade_dates" + MESSAGE_SPLIT + bc.user + MESSAGE_SPLIT +
		"1" + MESSAGE_SPLIT + strconv.FormatInt(BAOSTOCK_PER_PAGE_COUNT, 10) + MESSAGE_SPLIT +
		startDate + MESSAGE_SPLIT + endDate + MESSAGE_SPLIT
	msgHeader := messageHeader(MESSAGE_TYPE_QUERYTRADEDATES_REQUEST, int64(len(msgBody)))
	utils.Log("[QueryTradeDates] query trade dates from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.Write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return nil, fmt.Errorf("[QueryTradeDates] bc.Write fail\n\t%s", err)
	}
	recieveBody, err := bc.Read()
	if err != nil {
		return nil, fmt.Errorf("[QueryTradeDates] bc.Read fail\n\t%s", err)
	}
	recieveData, err := bc.DecodeRecieve(recieveBody)
	if err != nil {
		return nil, fmt.Errorf("[QueryTradeDates] bc.DecodeRecieve fail\n\t%s", err)
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return nil, fmt.Errorf("[QueryTradeDates] return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetQueryTradeDatesResponse()
	if err != nil {
		return nil, fmt.Errorf("[QueryTradeDates] recieveData.GetQueryTradeDatesResponse fail\n\t%s", err)
	}
	utils.Log("[QueryTradeDates] query trade dates from baostock success")
	return response, nil
}

func (bc *BaostockConnection) QueryAllStock(date string) (*QueryAllStockResponse, error) {
	if bc == nil || bc.connection == nil {
		return nil, fmt.Errorf("[QueryAllStock] bc is nil or connection is nil")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("[QueryAllStock] time.Parse fail %s", err)
		}
	}
	utils.Log("[QueryAllStock] querying all stock from baostock ...")
	if bc.user == "" {
		return nil, fmt.Errorf("[QueryAllStock] bc.user is nil, not login")
	}
	msgBody := "query_all_stock" + MESSAGE_SPLIT + bc.user + MESSAGE_SPLIT +
		"1" + MESSAGE_SPLIT + strconv.FormatInt(BAOSTOCK_PER_PAGE_COUNT, 10) + MESSAGE_SPLIT +
		date + MESSAGE_SPLIT
	msgHeader := messageHeader(MESSAGE_TYPE_QUERYALLSTOCK_REQUEST, int64(len(msgBody)))
	utils.Log("[QueryAllStock] query all stock from baostock param " + msgBody + " " + msgHeader)
	headBody := msgHeader + msgBody
	crcSum := crc32.ChecksumIEEE([]byte(headBody))
	_, err := bc.Write([]byte(headBody + MESSAGE_SPLIT + strconv.FormatInt(int64(crcSum), 10)))
	if err != nil {
		return nil, fmt.Errorf("[QueryAllStock] bc.Write fail\n\t%s", err)
	}
	recieveBody, err := bc.Read()
	if err != nil {
		return nil, fmt.Errorf("[QueryAllStock] bc.Read fail\n\t%s", err)
	}
	recieveData, err := bc.DecodeRecieve(recieveBody)
	if err != nil {
		return nil, fmt.Errorf("[QueryAllStock] bc.DecodeRecieve fail\n\t%s", err)
	}
	if recieveData.ErrorCode != BSERR_SUCCESS {
		return nil, fmt.Errorf("[QueryAllStock] return error code is not 0 %+v", recieveData)
	}
	response, err := recieveData.GetQueryAllStockResponse()
	if err != nil {
		return nil, fmt.Errorf("[QueryAllStock] recieveData.GetQueryAllStockResponse fail\n\t%s", err)
	}
	utils.Log("[QueryAllStock] query all stock from baostock success")
	return response, nil
}

func (bc *BaostockConnection) Write(data []byte) (int64, error) {
	if bc == nil || bc.connection == nil {
		return 0, fmt.Errorf("[Write] bc is nil or connection is nil")
	}
	if data == nil || len(data) <= 0 {
		return 0, fmt.Errorf("[Write] data is nil")
	}
	data = append(data, []byte("\n")...)
	n, err := bc.connection.Write(data)
	return int64(n), err
}

func (bc *BaostockConnection) Read() ([]byte, error) {
	if bc == nil || bc.connection == nil {
		return nil, fmt.Errorf("[Read] bc is nil or connection is nil")
	}
	data := []byte{}
	buf := make([]byte, 8192)
	for {
		n, err := bc.connection.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("[Read] bc.connection.Read fail\n\t%s", err)
		}
		// utils.Log("[Read] recive data " + string(buf[:n]))
		data = append(data, buf[:n]...)
		if len(data) >= 13 && bytes.Compare(data[len(data)-13:], []byte("<![CDATA[]]>\n")) == 0 {
			break
		}
	}
	return data, nil
}

func (bc *BaostockConnection) DecodeRecieve(data []byte) (*RecieveData, error) {
	if bc == nil || bc.connection == nil {
		return nil, fmt.Errorf("[DecodeRecieve] bc is nil or connection is nil")
	}
	if data == nil || len(data) <= 0 {
		return nil, fmt.Errorf("[DecodeRecieve] data is nil")
	}
	if len(data) < MESSAGE_HEADER_LENGTH {
		return nil, fmt.Errorf("[DecodeRecieve] data is error %s", string(data))
	}
	ret := &RecieveData{}
	headStr := string(data[:MESSAGE_HEADER_LENGTH])
	// utils.Log("[DecodeRecieve] recive data head " + string(headStr))
	headArr := strings.Split(headStr, MESSAGE_SPLIT)
	// utils.Log(fmt.Sprint("[DecodeRecieve] recive data head ", headArr))
	if len(headArr) < 3 {
		return nil, fmt.Errorf("[DecodeRecieve] data is error %s", string(data))
	}
	ret.MessageType = headArr[1]
	headInnerLength, err := strconv.ParseInt(headArr[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[DecodeRecieve] strconv.ParseInt fail\n\t%s", err)
	}
	ret.MessageBodyLength = headInnerLength
	bodyStr := ""
	if utils.StringIsIn(headArr[1], COMPRESSED_MESSAGE_TYPE_TUPLE...) {
		body, err := readSegment(data[MESSAGE_HEADER_LENGTH : MESSAGE_HEADER_LENGTH+headInnerLength])
		if err != nil {
			return nil, fmt.Errorf("[DecodeRecieve] readSegment fail\n\t%s", err)
		}
		bodyStr = string(body)
	} else {
		bodyStr = string(data[MESSAGE_HEADER_LENGTH:])
	}
	bodyArr := strings.Split(bodyStr, MESSAGE_SPLIT)
	// utils.Log(fmt.Sprint("[DecodeRecieve] recive data body ", bodyArr))
	if len(bodyArr) < 2 {
		return nil, fmt.Errorf("[DecodeRecieve] data is error %s", string(data))
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
		return nil, fmt.Errorf("[DecodeRecieve] ret.GetResponse fail\n\t%s", err)
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
