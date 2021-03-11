import tensorflow as tf
import tensorflow.keras as keras
import tensorflow.keras.layers as layers
import matplotlib.pyplot as plt
import matplotlib.dates as mdates
import mysql.connector as connector
import numpy as np
import pandas as pd

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
    x = np.zeros(shape=(len(data)-20, 15))
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
        # open high low close volume
        x[i-19][7] = data[i][1]
        x[i-19][8] = data[i][2]
        x[i-19][9] = data[i][3]
        x[i-19][10] = data[i][4]
        x[i-19][11] = data[i][5]
        # avg 5 10 20
        avg5, avg10, avg20 = 0, 0, 0
        for j in range(0, 20):
            if j < 5:
                avg5 += float(data[i-j][4])
            if j < 10:
                avg10 += float(data[i-j][4])
            if j < 20:
                avg20 += float(data[i-j][4])
        x[i-19][12] = avg5/5
        x[i-19][13] = avg10/10
        x[i-19][14] = avg20/20
        y[i-19] = float(data[i+1][1])
    # normal
    # x[:,6] = normal(x[:,6])
    # x[:,7] = normal(x[:,7])
    # x[:,8] = normal(x[:,8])
    # x[:,9] = normal(x[:,9])
    # x[:,10] = normal(x[:,10])
    x[:,11] = normal(x[:,11])
    # x[:,12] = normal(x[:,12])
    # x[:,13] = normal(x[:,13])
    # x[:,14] = normal(x[:,14])
    # y = normal(y)
    return x, y

def normal(data):
    _range = np.max(data) - np.min(data)
    return (data - np.min(data)) / _range


data = getData()
trainX, trainY = processData(data)

assert not np.any(np.isnan(trainX))
assert not np.any(np.isnan(trainY))

model = keras.Sequential([
    keras.layers.Dense(16, activation='relu'),
    keras.layers.Dense(16),
    keras.layers.Dense(1),
])
model.compile(
    optimizer=keras.optimizers.RMSprop(0.001), 
    loss='mse', 
    metrics=['mae', 'mse'])

history = model.fit(
    trainX, trainY, 
    epochs=100,
    validation_split=0.2,
    verbose=1,
)

hist = pd.DataFrame(history.history)
hist['epoch'] = history.epoch
plt.plot(hist['epoch'], hist['mse'], label = "train")
plt.plot(hist['epoch'], hist['val_mse'], label = "val")
plt.legend()
plt.show()