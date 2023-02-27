#!/bin/bash
curl -X GET -H "Content-Type: application/json" -d '{"sensorid":"010137B"}' localhost:5000/api/sensor
