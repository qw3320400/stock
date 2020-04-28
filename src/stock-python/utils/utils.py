import datetime

def Log(values):
    print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S') + " " + values)