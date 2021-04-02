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
    for i in range(2010, 2031):
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
    x = {
        'weekday': np.empty(shape=(len(data)-20), dtype='int'),
        'next_day_count': np.zeros(shape=(len(data)-20)),
        'avg5': np.zeros(shape=(len(data)-20)),
        'avg10': np.zeros(shape=(len(data)-20)),
        'avg20': np.zeros(shape=(len(data)-20)),
    }
    for j in range(10):
        x['open'+str(j)] = np.zeros(shape=(len(data)-20))
        x['high'+str(j)] = np.zeros(shape=(len(data)-20))
        x['low'+str(j)] = np.zeros(shape=(len(data)-20))
        x['close'+str(j)] = np.zeros(shape=(len(data)-20))
        x['volume'+str(j)] = np.zeros(shape=(len(data)-20))
    y = np.zeros(shape=(len(data)-20, 2))
    for i in range(len(data)):
        if i-19 < 0 or i+1 >= len(data):
            continue
        # week day
        x['weekday'][i-19] = data[i][0].weekday()+1
        # next day count
        timeDelta = data[i+1][0] - data[i][0]
        x['next_day_count'][i-19] = timeDelta.days
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
        x['avg5'][i-19] = avg5/5/base
        x['avg10'][i-19] = avg10/10/base
        x['avg20'][i-19] = avg20/20/base
        # open high low close volume
        for j in range(10): 
            x['open'+str(j)][i-19] = float(data[i-j][1])/base
            x['high'+str(j)][i-19] = float(data[i-j][2])/base
            x['low'+str(j)][i-19] = float(data[i-j][3])/base
            x['close'+str(j)][i-19] = float(data[i-j][4])/base
            x['volume'+str(j)][i-19] = float(data[i-j][5])/10000000000/volY[i-j][0]
        y[i-19][0] = float(data[i+1][2])/base
        y[i-19][1] = float(data[i+1][3])/base
    return x, y

def make_input_fn(data_df, label_df, num_epochs=5000, shuffle=True, batch_size=128):
    def input_function():
        ds = tf.data.Dataset.from_tensor_slices((dict(data_df), label_df))
        if shuffle:
            ds = ds.shuffle(10000)
        ds = ds.batch(batch_size).repeat(num_epochs)
        return ds
    return input_function

def trainHigh():
    data = getData()
    x, y = processData(data)

    feature_columns = []
    for key in x.keys():
        if key == 'weekday':
            vocabulary = np.unique(x[key])
            weekday_column = tf.feature_column.categorical_column_with_vocabulary_list(key, vocabulary)
            feature_columns.append(tf.feature_column.indicator_column(weekday_column))
        else:
            feature_columns.append(tf.feature_column.numeric_column(key, dtype=tf.float32))
    
    model = tf.estimator.BoostedTreesRegressor(
        feature_columns=feature_columns,
        n_batches_per_layer=10,
        model_dir='model/dnn_test',
    )
    tf.compat.v1.logging.set_verbosity(tf.compat.v1.logging.INFO)
    train_input_fn = make_input_fn(x, y[:,0])
    model.train(train_input_fn,
        # hooks=[],
    )

def predict():
    data = getData()
    x, y = processData(data)

    feature_columns = []
    for key in x.keys():
        if key == 'weekday':
            vocabulary = np.unique(x[key])
            weekday_column = tf.feature_column.categorical_column_with_vocabulary_list(key, vocabulary)
            feature_columns.append(tf.feature_column.indicator_column(weekday_column))
        else:
            feature_columns.append(tf.feature_column.numeric_column(key, dtype=tf.float32))
    
    model = tf.estimator.BoostedTreesRegressor(
        feature_columns=feature_columns,
        n_batches_per_layer=10,
        model_dir='model/dnn_test',
    )
    predictions = model.predict(make_input_fn(x, y[:,0], num_epochs=1))
    y_pred = np.array([pred['predictions'][0] for pred in predictions])
    mae = tf.keras.metrics.mean_absolute_error(y[:,0], y_pred)
    print(mae)


trainHigh()
# predict()
