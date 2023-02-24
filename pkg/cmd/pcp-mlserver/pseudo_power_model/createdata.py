#!/usr/bin/python3

# Import Library 
import pandas as pd 
import numpy as np 
import warnings
warnings.filterwarnings('ignore')
import lightgbm as lgb #LightGBM
from sklearn import datasets
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_squared_error
from sklearn.metrics import r2_score

# power consumption model
a = 0.3
b = 0.1
c = -0.5
d = 90.0
df = pd.DataFrame({'power':90, 'cpu':0, 'amb': 10, 'dp':0}, index=[0])

# create data
for cpu in range(100):
    for amb in range(10, 40):
        for dp in range(100):
            power = a * cpu + b * (amb - 20.0) + c * (dp / 10.0)  + d
            df = df.append({'power':power, 'cpu':cpu, 'amb':(amb -20.0) , 'dp':(dp / 10.0)}, ignore_index=True)
df.to_csv('data.csv')
