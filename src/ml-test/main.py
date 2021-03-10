# import tensorflow as tf
# import tensorflow.keras as keras
# import tensorflow.keras.layers as layers
import matplotlib.pyplot as plt
import matplotlib.dates as mdates
import mysql.connector as connector
import numpy as np

def getData():
    result = []
    db = connector.connect(
        host="localhost",
        user="root",
        passwd="123456",
        database="stock",
    )
    cursor = db.cursor()
    for i in range(2006, 2031):
        cursor.execute(" \
            select time_cst, open, high, low, close, volume, turn \
            from stock_k_data_%d \
            where code = 'sh.000300' and frequency = 'd' and adjust_flag = 'no' \
            order by time_cst \
        " % i)
        result.extend(cursor.fetchall())
    db.close()
    return result

def processData(data):
    x = np.zeros(shape=(len(data)-20, 16))
    y = np.zeros(shape=(len(data)-20,))
    for i in range(len(data)):
        if i-19 < 0 or i+1 >= len(data):
            continue
        # week day
        x[i-19][data[i][0].weekday()] = 1
        # next day count
        timeDelta = data[i+1][0] - data[i][0]
        x[i-19][5] = timeDelta.days
        # last close
        x[i-19][6] = float(data[i-1][4])
        # open high low close volume turn
        x[i-19][7] = data[i][1]
        x[i-19][8] = data[i][2]
        x[i-19][9] = data[i][3]
        x[i-19][10] = data[i][4]
        x[i-19][11] = data[i][5]
        x[i-19][12] = data[i][6]
        # avg 5 10 20
        avg5, avg10, avg20 = 0, 0, 0
        for j in range(0, 20):
            if j < 5:
                avg5 += float(data[i-j][4])
            if j < 10:
                avg10 += float(data[i-j][4])
            if j < 20:
                avg20 += float(data[i-j][4])
        x[i-19][13] = avg5/5
        x[i-19][14] = avg10/10
        x[i-19][15] = avg20/20
        if float(data[i+1][4]) > float(data[i][4]):
            y[i-19] = 1
    return x, y


data = getData()
trainX, trainY = processData(data)

x1 = []
for i in range(19, len(data)-1):
    x1.append(data[i][0])
y1 = trainX[:,10]
plt.plot(x1, y1, label = "close")
plt.gca().xaxis.set_major_formatter(mdates.DateFormatter('%Y-%m-%d'))
plt.legend()
plt.show()