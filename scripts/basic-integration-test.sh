#!/bin/bash
set -e

echo "Running basic integration test..."
echo "Deploying test pod (scheduler must already be running in cluster)..."

# Ensure custom namespace exists
kubectl create namespace custom-scheduler-namespace --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f infra/test/basic-test-container.yaml

echo "Waiting for pod to appear in cluster..."
sleep 3
kubectl get pods -n custom-scheduler-namespace

sleep 10

kubectl delete -f infra/test/basic-test-container.yaml
echo "Done. Check scheduler logs: kubectl logs -l app=custom-scheduler -n custom-scheduler -f"