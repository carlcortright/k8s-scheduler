#!/bin/bash
set -e

echo "Running basic integration test..."
echo "Deploying test pod (scheduler must already be running in cluster)..."

# Ensure custom namespace exists
kubectl create namespace custom-scheduler-namespace --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f infra/test/priority-class.yaml

kubectl apply -f infra/test/basic-test-container.yaml

echo "Waiting for pod to appear in cluster..."
sleep 3
kubectl get pods -n custom-scheduler-namespace

sleep 10

NODE=$(kubectl get pod basic-test -n custom-scheduler-namespace -o jsonpath='{.spec.nodeName}' 2>/dev/null || true)
if [ -z "$NODE" ]; then
  echo "FAIL: basic-test pod was not assigned a node"
  kubectl get pods -n custom-scheduler-namespace
  exit 1
fi
echo "OK: basic-test is scheduled on node $NODE"

kubectl delete -f infra/test/basic-test-container.yaml
echo "Done. Check scheduler logs: kubectl logs -l app=custom-scheduler -n custom-scheduler -f"