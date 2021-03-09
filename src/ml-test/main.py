import tensorflow as tf
import tensorflow.keras as keras
import tensorflow.keras.layers as layers
import matplotlib.pyplot as plt
import mysql.connector as connector

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
        " % i)
        result.extend(cursor.fetchall())
    db.close()
    return result


data = getData()
print(data)