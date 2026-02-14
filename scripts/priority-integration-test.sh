#!/bin/bash
set -e

echo "Running priority integration test..."
echo "Deploying three test pods (scheduler must already be running in cluster)..."

# Ensure custom namespace exists
kubectl create namespace custom-scheduler-namespace --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f infra/test/priority-class.yaml

for i in 1 2 3; do
  sed "s/name: basic-test/name: basic-test-$i/" infra/test/basic-test-container.yaml | kubectl apply -f -
done

echo "Waiting for pods to be scheduled..."
sleep 3
kubectl get pods -n custom-scheduler-namespace

sleep 10

FAIL=0
for i in 1 2 3; do
  NODE=$(kubectl get pod "basic-test-$i" -n custom-scheduler-namespace -o jsonpath='{.spec.nodeName}' 2>/dev/null || true)
  if [ -z "$NODE" ]; then
    echo "FAIL: basic-test-$i was not assigned a node"
    FAIL=1
  else
    echo "OK: basic-test-$i is scheduled on node $NODE"
  fi
done
if [ "$FAIL" -ne 0 ]; then
  kubectl get pods -n custom-scheduler-namespace
  exit 1
fi

echo "All three pods deployed and scheduled."
kubectl delete pod basic-test-1 basic-test-2 basic-test-3 -n custom-scheduler-namespace --ignore-not-found=true
echo "Done. Check scheduler logs: kubectl logs -l app=custom-scheduler -n custom-scheduler -f"