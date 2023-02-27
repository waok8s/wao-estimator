#!/usr/bin/python
import serial
import time
import ctypes
from ast import literal_eval
import sqlite3

# Database
dbname = "../instance/SDP810.db"

# Open database connection
conn = sqlite3.connect(dbname)

cur = conn.cursor()
command = 'REPLACE INTO sensors values("0000001", 4.5, 21.2)'
cur.execute(command)
conn.commit()
