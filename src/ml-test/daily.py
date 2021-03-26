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
            where code = '000300.XSHG' and frequency = '5' and adjust_flag = 'no' \
            order by time_cst \
        " % i)
        result.extend(cursor.fetchall())
    db.close()
    return result

def processData(data):
    return None

def trainHigh():
    return None

def watch():
    data = getData()

    date = data[0][0].date()
    maxTime = data[0][0].strftime('%H-%M-%S')
    minTime = data[0][0].strftime('%H-%M-%S')
    max = float(0)
    min = float(99999)
    mMap = {}
    for i in range(len(data)) :
        date = data[i][0].date()
        if float(data[i][4]) > max :
            max = float(data[i][4])
            maxTime = data[i][0].strftime('%H-%M-%S')
        if float(data[i][4]) < min :
            min = float(data[i][4])
            minTime = data[i][0].strftime('%H-%M-%S')

        if i == len(data) - 1 or date != data[i+1][0].date() :
            if maxTime not in mMap:
                mMap[maxTime] = {'max':0,'min':0}
            mMap[maxTime]['max'] = mMap[maxTime]['max'] + 1
            if minTime not in mMap:
                mMap[minTime] = {'max':0,'min':0}
            mMap[minTime]['min'] = mMap[minTime]['min'] + 1        
            
            date = data[i][0].date()
            maxTime = data[0][0].strftime('%H-%M-%S')
            minTime = data[0][0].strftime('%H-%M-%S')
            max = float(0)
            min = float(99999)
    
    df = {
        'time' : [],
        'max' : [],
        'min' : [],
    }
    for key in mMap :
        df['time'].append(pd.to_datetime(key, format = '%H-%M-%S'))
        df['max'].append(mMap[key]['max'])
        df['min'].append(mMap[key]['min'])
    t = pd.DataFrame(df)
    t.plot(x='time', y='max', lw=0, marker='o', figsize=(8,4)) 
    t.plot(x='time', y='min', lw=0, marker='o', figsize=(8,4)) 
    # plt.legend()
    plt.show()

watch()