#!/bin/bash
curl -X POST -H "Content-Type: application/json" -d '{"sensorid":"0000001", "pressure": 5.1, "temperature": 22.2}' localhost:5000/api/sensors

curl -X POST -H "Content-Type: application/json" -d '{"sensorid":"0000002", "pressure": 4.1, "temperature": 27.2}' localhost:5000/api/sensors



