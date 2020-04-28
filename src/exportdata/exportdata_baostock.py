import baostock
import pandas
import datetime
import os
import string
import sys
sys.path.append("..")

from utils.utils import (
    Log,
)

class DataFrequency:
    def __init__(self, frequency, fields):
        self.frequency = frequency
        self.fields = fields

dataFrequencyList = [
    DataFrequency("d", "date,code,open,high,low,close,preclose,volume,amount,adjustflag,turn,tradestatus,pctChg,peTTM,psTTM,pcfNcfTTM,pbMRQ,isST"),
    DataFrequency("w", "date,code,open,high,low,close,volume,amount,adjustflag,turn,pctChg"),
    DataFrequency("m", "date,code,open,high,low,close,volume,amount,adjustflag,turn,pctChg"),
    DataFrequency("5", "date,code,open,high,low,close,volume,amount,adjustflag"),
    DataFrequency("15", "date,code,open,high,low,close,volume,amount,adjustflag"),
    DataFrequency("30", "date,code,open,high,low,close,volume,amount,adjustflag"),
    DataFrequency("60", "date,code,open,high,low,close,volume,amount,adjustflag"),
]

baostockDataStartDate = datetime.datetime.strptime("2006-01-01", '%Y-%m-%d')

def BaostockLogin():
    #### 登陆系统 ####
    lg = baostock.login()
    Log('login respond error_code:' + lg.error_code + ' error_msg:' + lg.error_msg)

def BaostockLogout():
    #### 登出系统 ####
    baostock.logout()

def ExportBaostockDataByMonth(code, date):
    filePath = os.path.dirname(__file__)
    savePath = os.path.join(filePath, "../../data/baostock/{:d}/{:d}".format(date.year, date.month))
    os.makedirs(savePath, exist_ok=True)
    
    # frequency
    for dataFrequency in dataFrequencyList:
        startTime = datetime.date(date.year, date.month, 1)
        endTime = startTime
        if date.month == 12:
            endTime = datetime.date(date.year + 1, 1, 1) - datetime.timedelta(days = 1)
        else:
            endTime = datetime.date(date.year, date.month + 1, 1) - datetime.timedelta(days = 1)

        # 后复权
        saveName = "{:s}:{:s}:{:s}:1.csv".format(code, date.strftime('%Y-%m'), dataFrequency.frequency)
        saveFile = os.path.join(savePath, saveName)
        Log('check file {:s}'.format(saveName))
        if os.path.exists(saveFile) == False:
            rs = baostock.query_history_k_data_plus(code, dataFrequency.fields,
                start_date = startTime.strftime('%Y-%m-%d') , end_date = endTime.strftime('%Y-%m-%d'), frequency = dataFrequency.frequency, adjustflag = "1")
            Log('query_history_k_data_plus error_code:{:s} error_msg:{:s} {:s}'.format(rs.error_code, rs.error_msg, saveName))
            #### 打印结果集 ####
            data_list = []
            while (rs.error_code == '0') & rs.next():
                # 获取一条记录，将记录合并在一起
                data_list.append(rs.get_row_data())
            #### 结果集输出到csv文件 ####
            result = pandas.DataFrame(data_list, columns = rs.fields)
            result.to_csv(saveFile, encoding = "utf-8", index = False)

        # 不复权
        saveName = "{:s}:{:s}:{:s}:3.csv".format(code, date.strftime('%Y-%m'), dataFrequency.frequency)
        saveFile = os.path.join(savePath, saveName)
        Log('check file {:s}'.format(saveName))
        if os.path.exists(saveFile) == False:
            rs = baostock.query_history_k_data_plus(code, dataFrequency.fields,
                start_date = startTime.strftime('%Y-%m-%d') , end_date = endTime.strftime('%Y-%m-%d'), frequency = dataFrequency.frequency, adjustflag = "3")
            Log('query_history_k_data_plus error_code:{:s} error_msg:{:s} {:s}'.format(rs.error_code, rs.error_msg, saveName))
            #### 打印结果集 ####
            data_list = []
            while (rs.error_code == '0') & rs.next():
                # 获取一条记录，将记录合并在一起
                data_list.append(rs.get_row_data())
            #### 结果集输出到csv文件 ####
            result = pandas.DataFrame(data_list, columns = rs.fields)
            result.to_csv(saveFile, encoding = "utf-8", index = False)

def ExportBaostockTradeDate():
    filePath = os.path.dirname(__file__)
    savePath = os.path.join(filePath, "../../data/baostock/tradedate")
    os.makedirs(savePath, exist_ok=True)

    nowDate = datetime.datetime.now()
    saveName = "tradedate:{:s}.csv".format(nowDate.strftime('%Y-%m-%d'))
    saveFile = os.path.join(savePath, saveName)
    rs = baostock.query_trade_dates(start_date = baostockDataStartDate.strftime('%Y-%m-%d'), end_date = nowDate.strftime('%Y-%m-%d'))
    Log('query_trade_dates respond error_code:' + rs.error_code + ' error_msg:' + rs.error_msg + ' filename:' + saveName)
    #### 打印结果集 ####
    data_list = []
    while (rs.error_code == '0') & rs.next():
        # 获取一条记录，将记录合并在一起
        data_list.append(rs.get_row_data())
    result = pandas.DataFrame(data_list, columns = rs.fields)
    #### 结果集输出到csv文件 ####
    result.to_csv(saveFile, encoding = "utf-8", index = False)

    return data_list

def ExportBaostockAllStock(date):
    filePath = os.path.dirname(__file__)
    savePath = os.path.join(filePath, "../../data/baostock/allstock")
    os.makedirs(savePath, exist_ok=True)
    
    saveName = "allstock:{:s}.csv".format(date.strftime('%Y-%m-%d'))
    saveFile = os.path.join(savePath, saveName)
    rs = baostock.query_all_stock(day = date.strftime('%Y-%m-%d'))
    Log('query_all_stock respond error_code:' + rs.error_code + ' error_msg:' + rs.error_msg + ' filename:' + saveName)
    #### 打印结果集 ####
    data_list = []
    while (rs.error_code == '0') & rs.next():
        # 获取一条记录，将记录合并在一起
        data_list.append(rs.get_row_data())
    result = pandas.DataFrame(data_list, columns = rs.fields)
    #### 结果集输出到csv文件 ####
    result.to_csv(saveFile, encoding = "utf-8", index = False)

    return data_list