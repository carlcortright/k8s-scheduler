#!/bin/bash

kubectl apply -f infra/test/basic-test-container.yaml 

sleep 10

kubectl delete -f infra/test/basic-test-container.yaml