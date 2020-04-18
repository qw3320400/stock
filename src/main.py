import datetime

from exportdata.exportdata_baostock import *
from utils.utils import *

def main():
    Log("==start==")

    tradeDateData = ExportBaostockTradeDate()
    lastTradeDateStr = ""
    index = len(tradeDateData) - 1
    while (True):
        if tradeDateData[index][1] == "1":
            # this is a trade date
            lastTradeDateStr = tradeDateData[index][0]
            break
        index -= 1

    allStockData = ExportBaostockAllStock(datetime.datetime.strptime(lastTradeDateStr, '%Y-%m-%d'))
    codeList = []
    for stock in allStockData:
        codeList.append(stock[0])
    
    startTimeStr = '2006-01-01'
    # 2020-04-18
    endTimeStr = '2020-03-01'
    startTime = datetime.datetime.strptime(startTimeStr, '%Y-%m-%d')
    endTime =  datetime.datetime.strptime(endTimeStr, '%Y-%m-%d')
    ExportBaostockData(codeList, startTime, endTime)
    
    Log("==end==")

def ExportBaostockData(codeList, startTime, endTime):
    if startTime > endTime:
        return
    # each day
    while (True):
        # check if break
        if startTime > endTime:
            break
        # each code
        for code in codeList:
            ExportBaostockDataByMonth(code, endTime)
        # sub 1 month
        if endTime.month == 1:
            endTime = datetime.date(endTime.year - 1, 12, 1)
        else:
            endTime = datetime.date(endTime.year, endTime.month - 1, 1)

if __name__ == "__main__":
    main()