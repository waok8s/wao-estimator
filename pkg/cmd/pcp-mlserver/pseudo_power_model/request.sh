#!/bin/bash

curl http://localhost:8080/v2/models/model1/versions/v0.1.0/infer -H "Content-Type: application/json" -d '{"inputs":[{"name":"predict-prob",  "shape":[1,3], "datatype":"FP32", "data":[[10.0,22.0,0.2]]}]}'

