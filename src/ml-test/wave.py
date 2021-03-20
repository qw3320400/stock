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
    baseDate = datetime.datetime.strptime('2006-01-01', '%Y-%m-%d')
    volModel = keras.models.load_model('model/my_volume')
    volX = np.zeros(shape=(len(data), 1))
    volY = np.zeros(shape=(len(data),))
    for i in range(len(data)):
        volX[i][0] = (data[i][0] - baseDate).days
    volY = volModel.predict(volX)
    x = np.zeros(shape=(len(data)-20, 64))
    y = np.zeros(shape=(len(data)-20, 2))
    for i in range(len(data)):
        if i-19 < 0 or i+1 >= len(data):
            continue
        # week day
        x[i-19][data[i][0].weekday()] = 1
        # next day count
        timeDelta = data[i+1][0] - data[i][0]
        x[i-19][5] = timeDelta.days
        # base
        base = float(data[i][4])
        # avg 5 10 20
        avg5, avg10, avg20 = 0, 0, 0
        for j in range(0, 20):
            if j < 5:
                avg5 += float(data[i-j][4])
            if j < 10:
                avg10 += float(data[i-j][4])
            if j < 20:
                avg20 += float(data[i-j][4])
        x[i-19][6] = avg5/5/base
        x[i-19][7] = avg10/10/base
        x[i-19][8] = avg20/20/base
        # open high low close volume
        for j in range(10): 
            x[i-19][9+j*5] = float(data[i-j][1])/base
            x[i-19][10+j*5] = float(data[i-j][2])/base
            x[i-19][11+j*5] = float(data[i-j][3])/base
            x[i-19][12+j*5] = float(data[i-j][4])/base
            x[i-19][13+j*5] = float(data[i-j][5])/10000000000/volY[i-j][0]

        y[i-19][0] = float(data[i+1][2])/base
        y[i-19][1] = float(data[i+1][3])/base
    return x, y

## train high ##
def trainHigh():
    data = getData()
    trainX, trainY = processData(data)

    assert not np.any(np.isnan(trainX))
    assert not np.any(np.isnan(trainY))

    model = keras.Sequential([
        keras.layers.Dense(64),
        keras.layers.Dense(64),
        keras.layers.Dense(32),
        keras.layers.Dense(1),
    ])
    model.compile(
        optimizer='adam', 
        loss='mse', 
        metrics=['mae', 'mse'])
    model.load_weights('checkpoints/my_high_64')
    # tensorboard_callback = keras.callbacks.TensorBoard(log_dir='logs') # tensorboard --logdir logs
    history = model.fit(
        trainX, trainY[:,0], 
        epochs=5000,
        validation_split=0.2,
        verbose=1,
        # callbacks=[tensorboard_callback],
    )
    model.save('model/my_high_64')
    model.save_weights('checkpoints/my_high_64')

    hist = pd.DataFrame(history.history)
    hist['epoch'] = history.epoch
    plt.plot(hist['epoch'], hist['val_mae'], label = "val")
    plt.plot(hist['epoch'], hist['mae'], label = "train")
    plt.legend()
    plt.show()

## train low ##
def trainLow():
    data = getData()
    trainX, trainY = processData(data)

    assert not np.any(np.isnan(trainX))
    assert not np.any(np.isnan(trainY))

    model = keras.Sequential([
        keras.layers.Dense(64),
        keras.layers.Dense(64),
        keras.layers.Dense(32),
        keras.layers.Dense(1),
    ])
    model.compile(
        optimizer='adam', 
        loss='mse', 
        metrics=['mae', 'mse'])
    model.load_weights('checkpoints/my_low_64')
    # tensorboard_callback = keras.callbacks.TensorBoard(log_dir='logs') # tensorboard --logdir logs
    history = model.fit(
        trainX, trainY[:,1], 
        epochs=5000,
        validation_split=0.2,
        verbose=1,
        # callbacks=[tensorboard_callback],
    )
    model.save('model/my_low_64')
    model.save_weights('checkpoints/my_low_64')

    hist = pd.DataFrame(history.history)
    hist['epoch'] = history.epoch
    plt.plot(hist['epoch'], hist['val_mae'], label = "val")
    plt.plot(hist['epoch'], hist['mae'], label = "train")
    plt.legend()
    plt.show()

def predictData(data):
    baseDate = datetime.datetime.strptime('2006-01-01', '%Y-%m-%d')
    volModel = keras.models.load_model('model/my_volume')
    volX = np.zeros(shape=(len(data), 1))
    volY = np.zeros(shape=(len(data),))
    for i in range(len(data)):
        volX[i][0] = (data[i][0] - baseDate).days
    volY = volModel.predict(volX)
    x = np.zeros(shape=(len(data)-19, 64))
    baseData = [[]] * (len(data)-19)
    for i in range(len(data)):
        if i-19 < 0 :
            continue
        # week day
        x[i-19][data[i][0].weekday()] = 1
        # next day count
        if i+1 >= len(data): # todo
            x[i-19][5] = 1
        else:
            timeDelta = data[i+1][0] - data[i][0]
            x[i-19][5] = timeDelta.days
        # base
        base = float(data[i][4])
        # avg 5 10 20
        avg5, avg10, avg20 = 0, 0, 0
        for j in range(0, 20):
            if j < 5:
                avg5 += float(data[i-j][4])
            if j < 10:
                avg10 += float(data[i-j][4])
            if j < 20:
                avg20 += float(data[i-j][4])
        x[i-19][6] = avg5/5/base
        x[i-19][7] = avg10/10/base
        x[i-19][8] = avg20/20/base
        # open high low close volume
        for j in range(10): 
            x[i-19][9+j*5] = float(data[i-j][1])/base
            x[i-19][10+j*5] = float(data[i-j][2])/base
            x[i-19][11+j*5] = float(data[i-j][3])/base
            x[i-19][12+j*5] = float(data[i-j][4])/base
            x[i-19][13+j*5] = float(data[i-j][5])/10000000000/volY[i-j][0]

        baseData[i-19] = [data[i][0], base]
    return x, baseData

## predict ##
def predict():
    data = getData()
    x, baseData = predictData(data)
    x[len(x)-1][5] = 3 # todo
    hgihModel = keras.models.load_model('model/my_high_64')
    lowModel = keras.models.load_model('model/my_low_64')
    highY = hgihModel.predict(x)
    lowY = lowModel.predict(x)
    time = datetime.datetime.strptime('2021-03-20', '%Y-%m-%d')
    for i in range(len(baseData)):
        if (time - baseData[i][0]).days <= 10 :
            base = baseData[i][1]
            print(baseData[i][0], " -> base: ", base, " high: ", highY[i]*base, " low: ", lowY[i]*base)

## watch ##
def watch():
    data = getData()
    trainX, trainY = processData(data)
    hgihModel = keras.models.load_model('model/my_high_64')
    lowModel = keras.models.load_model('model/my_low_64')
    highY = hgihModel.predict(trainX)
    lowY = lowModel.predict(trainX)

    plt.scatter(range(len(trainX)), highY, label = 'pre_h', s = 1)
    plt.scatter(range(len(trainX)), lowY, label = 'pre_l', s = 1)
    # plt.plot(range(len(trainX)), trainY[:,0], label = 'act_h')
    # plt.plot(range(len(trainX)), trainY[:,1], label = 'act_l')
    plt.legend()
    plt.show()

# trainHigh()
# trainLow()
# watch()
predict()