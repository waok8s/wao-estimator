1. センサーAPIの起動
./app.py
2. センサーへのアクセス例。センサーID 101037Bのデータを取得
curl localhost:5000/api/sensor/101037B | jq
3. 取得結果
{
  "code": 200,
  "sensor": [
    {
      "pressure": 0.01,
      "sensorid": "101037B",
      "temperature": 26.02
    }
  ]
}
