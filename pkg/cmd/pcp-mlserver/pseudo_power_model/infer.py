#!/usr/bin/python3

import requests
import numpy as np

x_0 = np.array([[10.0, 22.0, 0.2]]) 
inference_request = {
    "inputs": [
        {
          "name": "predict-prob",
          "shape": x_0.shape,
          "datatype": "FP32",
          "data": x_0.tolist()
        }
    ]
}

endpoint = "http://localhost:8080/v2/models/model1/versions/v0.1.0/infer"
response = requests.post(endpoint, json=inference_request)

print(response.json())
