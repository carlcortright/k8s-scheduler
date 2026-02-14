#!/bin/bash

echo "Running basic integration test..."
echo "Deploying custom scheduler to local cluster..."


kubectl apply -f infra/test/basic-test-container.yaml 

sleep 10

kubectl delete -f infra/test/basic-test-container.yaml