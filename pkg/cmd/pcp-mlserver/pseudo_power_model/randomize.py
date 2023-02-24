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

df = pd.read_csv('data.csv', index_col=0)
dfNew = df.sample(frac=1, ignore_index=True)
print(dfNew)
dfNew.to_csv('randomizedData.csv')
