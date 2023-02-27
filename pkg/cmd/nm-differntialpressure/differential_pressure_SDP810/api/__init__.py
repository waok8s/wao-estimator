from flask import Flask, make_response, jsonify
from .views.sensor import sensor_router
from flask_cors import CORS
from api.database import db
import config

def create_app():

  app = Flask(__name__)

  # CORS対応
  CORS(app)

  # DB設定を読み込む
  app.config.from_object('config.Config')
  db.init_app(app)

  app.register_blueprint(sensor_router, url_prefix='/api')

  return app

app = create_app()
