#!/usr/bin/python
import serial
import time
import ctypes
from ast import literal_eval
import sqlite3

# Serial interface to TWELITE
ser=serial.Serial('/dev/serial0', 115200)

# Database
dbname = "/home/kazuhiro/SDP810/instance/SDP810.db"

# Open database connection
conn = sqlite3.connect(dbname)

while True:
    # Receive data from serial I/F 
    try:
        readText = ser.readline()
        sensorid = readText[8:15].decode('utf-8')
        scale = ctypes.c_int16(literal_eval("0x" + readText[45:49].decode('utf-8'))).value
        pressure = ctypes.c_int16(literal_eval("0x" + readText[37:41].decode('utf-8'))).value / scale
        temperature = ctypes.c_int16(literal_eval("0x" + readText[41:45].decode('utf-8'))).value / scale
        print("sensorid: ", sensorid, "pressure: ", pressure, " temperature: ", temperature) 
        time.sleep(1)

    except serial.SerialException:
        print("/dev/serial0 not found")

    # Insert data into database
    try:
        cur = conn.cursor()
        sql = 'REPLACE INTO sensors (sensorid, pressure, temperature) values (?, ?, ?);'
        data = (sensorid, pressure, temperature)
        cur.execute(sql, data)
        conn.commit()

    except KeyboardInterrupt:
        conn.close()
        exit()
