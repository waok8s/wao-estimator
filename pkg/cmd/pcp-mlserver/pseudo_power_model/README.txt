1. 疑似データを作成（出力はdata.csv）
./createdata.py
2. 学習用にデータの並びをランダム化（出力はrandomizedData.csv）
./randomize.py
3. lightGBMによる学習モデルを生成（出力はpseudo_power_model.bst）
./learning.py
4. モデルを格納するディレクトリを作成
mkdir -p models/model1
5. モデルを移動
mv pseudo_power_model.bst models/model1
5. MLServer用のモデル設定ファイルを作成（models/model1/model-settings.json）
6. MLServerの設定ファイルを作成（settings.json）
7. MLServerを起動
mlserver start .
8. 推論のサンプルプログラムを起動（Pythonクライアント）
./infer.py
9. 推論のサンプルスクリプトを起動（curlコマンド）
./request.sh
