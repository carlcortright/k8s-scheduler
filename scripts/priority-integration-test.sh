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

echo "Deploying high-priority pod (should preempt one basic-test)..."
kubectl apply -f infra/test/priority-test-container.yaml

echo "Waiting for preemption and scheduling..."
sleep 15

# High-priority pod must be scheduled
PRIO_NODE=$(kubectl get pod priority-test -n custom-scheduler-namespace -o jsonpath='{.spec.nodeName}' 2>/dev/null || true)
if [ -z "$PRIO_NODE" ]; then
  echo "FAIL: priority-test was not assigned a node"
  kubectl get pods -n custom-scheduler-namespace
  exit 1
fi
echo "OK: priority-test is scheduled on node $PRIO_NODE"

# Exactly two of basic-test-1/2/3 should still be running (one was evicted)
RUNNING_BASIC=0
for i in 1 2 3; do
  if kubectl get pod "basic-test-$i" -n custom-scheduler-namespace &>/dev/null; then
    RUNNING_BASIC=$((RUNNING_BASIC + 1))
  fi
done
if [ "$RUNNING_BASIC" -ne 2 ]; then
  echo "FAIL: expected exactly 2 basic-test pods remaining (one evicted), got $RUNNING_BASIC"
  kubectl get pods -n custom-scheduler-namespace
  exit 1
fi
echo "OK: exactly 2 basic-test pods remaining (1 evicted for priority-test)"

kubectl delete pod basic-test-1 basic-test-2 basic-test-3 priority-test -n custom-scheduler-namespace --ignore-not-found=true
echo "Done. Check scheduler logs: kubectl logs -l app=custom-scheduler -n custom-scheduler -f"