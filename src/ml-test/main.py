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
    x = np.zeros(shape=(len(data)-20, 24))
    y = np.zeros(shape=(len(data)-20,))
    for i in range(len(data)):
        if i-19 < 0 or i+1 >= len(data):
            continue
        # week day
        x[i-19][data[i][0].weekday()] = 1
        # next day count
        timeDelta = data[i+1][0] - data[i][0]
        x[i-19][5] = timeDelta.days
        # open high low close volume
        base = float(data[i-2][1])
        x[i-19][6] = float(data[i-2][1])/base
        x[i-19][7] = float(data[i-2][2])/base
        x[i-19][8] = float(data[i-2][3])/base
        x[i-19][9] = float(data[i-2][4])/base
        x[i-19][10] = float(data[i-2][5])/10000000000/volY[i-2][0]
        x[i-19][11] = float(data[i-1][1])/base
        x[i-19][12] = float(data[i-1][2])/base
        x[i-19][13] = float(data[i-1][3])/base
        x[i-19][14] = float(data[i-1][4])/base
        x[i-19][15] = float(data[i-1][5])/10000000000/volY[i-1][0]
        x[i-19][16] = float(data[i][1])/base
        x[i-19][17] = float(data[i][2])/base
        x[i-19][18] = float(data[i][3])/base
        x[i-19][19] = float(data[i][4])/base
        x[i-19][20] = float(data[i][5])/10000000000/volY[i][0]
        # avg 5 10 20
        avg5, avg10, avg20 = 0, 0, 0
        for j in range(0, 20):
            if j < 5:
                avg5 += float(data[i-j][4])
            if j < 10:
                avg10 += float(data[i-j][4])
            if j < 20:
                avg20 += float(data[i-j][4])
        x[i-19][21] = avg5/5/base
        x[i-19][22] = avg10/10/base
        x[i-19][23] = avg20/20/base
        y[i-19] = float(data[i+1][2])/base
    return x, y


# data = getData()
# trainX, trainY = processData(data)

# assert not np.any(np.isnan(trainX))
# assert not np.any(np.isnan(trainY))

# model = keras.Sequential([
#     keras.layers.Dense(32, activation='relu'),
#     keras.layers.Dense(32),
#     keras.layers.Dense(1),
# ])
# model.compile(
#     optimizer=keras.optimizers.RMSprop(0.00001), 
#     loss='mae', 
#     metrics=['mae', 'mse'])
# model.load_weights('checkpoints/my_high')
# history = model.fit(
#     trainX, trainY, 
#     epochs=5000,
#     validation_split=0.2,
#     verbose=1,
# )
# model.save('model/my_high')
# model.save_weights('checkpoints/my_high')

# hist = pd.DataFrame(history.history)
# hist['epoch'] = history.epoch
# plt.plot(hist['epoch'], hist['mae'], label = "train")
# plt.plot(hist['epoch'], hist['val_mae'], label = "val")
# plt.legend()
# plt.show()

testData = np.array([[0,0,0,0,1,3,
    5047.06,5055.28,4981.62,5003.60,1.70,
    5024.65,5138.41,5020.58,5128.22,1.90,
    5153.67,5153.67,5086.82,5146.38,2.01,
    5065.84,5209.33,5406.00]])
baseDate = datetime.datetime.strptime('2006-01-01', '%Y-%m-%d')
volX = np.zeros(shape=(3, 1))
volY = np.zeros(shape=(3,))
volX[0][0] = (datetime.datetime.strptime('2021-03-10', '%Y-%m-%d') - baseDate).days
volX[1][0] = (datetime.datetime.strptime('2021-03-11', '%Y-%m-%d') - baseDate).days
volX[2][0] = (datetime.datetime.strptime('2021-03-12', '%Y-%m-%d') - baseDate).days
volModel = keras.models.load_model('model/my_volume')
volY = volModel.predict(volX)
base = testData[0,6]
testData[0,6:10] = testData[0,6:10]/base
testData[0,10] = testData[0,10]/volY[0]
testData[0,11:15] = testData[0,11:15]/base
testData[0,15] = testData[0,15]/volY[1]
testData[0,16:20] = testData[0,16:20]/base
testData[0,20] = testData[0,20]/volY[2]
testData[0,21:24] = testData[0,21:24]/base
hgihModel = keras.models.load_model('model/my_high')
lowModel = keras.models.load_model('model/my_low')
high = hgihModel.predict(testData)
low = lowModel.predict(testData)
print(testData)
print(high*base, low*base)