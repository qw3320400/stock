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
        y[i] = float(data[i][4])/1000
    return x, y


def train():
    data = getData()
    trainX, trainY = processData(data)

    model = keras.Sequential([
        keras.layers.Dense(1),
    ])
    model.compile(
        optimizer='adam',
        loss='mse',
        metrics=['mae', 'mse'],
    )
    model.load_weights('checkpoints/price_line')
    model.fit(
        trainX, trainY,
        epochs=500,
    )
    model.save('model/price_line')
    model.save_weights('checkpoints/price_line')

def watch():
    data = getData()
    trainX, trainY = processData(data)

    model = keras.models.load_model('model/price_line')
    y = model.predict(trainX)
    plt.plot(trainX, trainY, label = "price_line")
    plt.plot(trainX, y[:,0])
    plt.plot(trainX, trainY/y[:,0])
    plt.plot(trainX, np.ones(shape=[len(trainX),]))
    plt.legend()
    plt.show()

def predict():
    data = getData()
    baseDate = datetime.datetime.strptime('2006-01-01', '%Y-%m-%d')
    volModel = keras.models.load_model('model/price_line')
    volBase = volModel.predict([[(data[3000][0] - baseDate).days]])
    print(volBase)


# train()
watch()