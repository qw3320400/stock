import tensorflow as tf
import tensorflow.keras as keras
import tensorflow.keras.layers as layers
import matplotlib.pyplot as plt
import matplotlib.dates as mdates
import mysql.connector as connector
import numpy as np
import pandas as pd
import datetime

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
    baseDay = datetime.datetime.strptime('2006-01-01', '%Y-%m-%d')
    x = np.zeros(shape=(len(data),1))
    y = np.zeros(shape=(len(data),))
    for i in range(len(data)):
        x[i][0] = (data[i][0] - baseDay).days
        y[i] = float(data[i][5])/10000000000
    return x, y


data = getData()
trainX, trainY = processData(data)

model = keras.Sequential([
    keras.layers.Dense(1),
])
model.compile(
    optimizer=keras.optimizers.RMSprop(0.0001),
    loss='mae',
    metrics=['mae', 'mse'],
)
model.fit(
    trainX, trainY,
    epochs=500,
)
model.save('model/my_volume')

y = model.predict(trainX)
plt.plot(trainX, trainY, label = "volume")
plt.plot(trainX, y)
plt.legend()
plt.show()

# data = getData()
# baseDate = datetime.datetime.strptime('2006-01-01', '%Y-%m-%d')
# volModel = keras.models.load_model('model/my_volume')
# volBase = volModel.predict([[(data[3000][0] - baseDate).days]])
# print(volBase)