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
import pickle
import optuna

df = pd.read_csv('randomizedData.csv', index_col=0)

X = df.drop('power', axis=1).values
y = df['power'].values
print(X)
print(y)

X_train,X_test,y_train,y_test = train_test_split(X, y, test_size = 0.20, random_state = 2)
lgb_train = lgb.Dataset(X_train, y_train)
lgb_eval = lgb.Dataset(X_test, y_test, reference = lgb_train)

params = {
        'task': 'train',
        'boosting_type': 'gbdt',
        'objective': 'regression', # 目的 : 回帰  
        'metric': {'rmse'}, # 評価指標 : rsme(平均二乗誤差の平方根) 
}

# モデルの学習
model = lgb.train(params,
                  train_set=lgb_train, # トレーニングデータの指定
                  valid_sets=lgb_eval, # 検証データの指定
                  )

# テストデータの予測
y_pred = model.predict(X_test)
print(X_test)

# モデル評価
# rmse : 平均二乗誤差の平方根
mse = mean_squared_error(y_test, y_pred) # MSE(平均二乗誤差)の算出
rmse = np.sqrt(mse) # RSME = √MSEの算出
print('RMSE :',rmse)

#r2 : 決定係数
r2 = r2_score(y_test,y_pred)
print('R2 :',r2)

# モデルの保存
model.save_model('pseudo_power_model.bst', num_iteration=model.best_iteration)
