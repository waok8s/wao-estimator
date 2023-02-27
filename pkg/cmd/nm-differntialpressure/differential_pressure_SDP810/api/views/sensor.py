from flask import Blueprint, request, make_response, jsonify
from api.models import Sensor, SensorSchema
import json

# Contents routing 
sensor_router = Blueprint('sensor_router', __name__)

@sensor_router.route('/sensors', methods=['GET'])
def getSensorList():

  sensors = Sensor.getSensorList()
  sensor_schema = SensorSchema(many=True)

  return make_response(jsonify({
    'code': 200,
    'sensors': sensor_schema.dump(sensors)
  }))

@sensor_router.route('/sensors', methods=['POST'])
def registSensor():

  # get data in json format
  jsonData = json.dumps(request.json)
  sensorData = json.loads(jsonData)

  sensor = Sensor.registSensor(sensorData)
  sensor_schema = SensorSchema(many=True)

  return make_response(jsonify({
    'code': 200,
    'sensor': sensor
  }))


@sensor_router.route('/sensor/<sensorid>', methods=['GET'])
def getSensor(sensorid):

  sensor = Sensor.getSensor(sensorid)
  sensor_schema = SensorSchema(many=True)

  return make_response(jsonify({
      'code': 200,
      'sensor': sensor_schema.dump(sensor)
  }))
