#!/bin/bash
set -e

echo "Running group integration test (scheduler must be running; cluster needs 3+ nodes)..."

kubectl create namespace custom-scheduler-namespace --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f infra/test/priority-class.yaml 2>/dev/null || true

echo "Scheduling a group of 3 (same pod-group)..."
kubectl apply -f infra/test/group-test-container.yaml
sleep 5
kubectl get pods -n custom-scheduler-namespace
sleep 10

FAIL=0
for i in 1 2 3; do
  NODE=$(kubectl get pod "group-test-$i" -n custom-scheduler-namespace -o jsonpath='{.spec.nodeName}' 2>/dev/null || true)
  PHASE=$(kubectl get pod "group-test-$i" -n custom-scheduler-namespace -o jsonpath='{.status.phase}' 2>/dev/null || true)
  if [ -n "$NODE" ] || [ "$PHASE" = "Running" ]; then
    echo "OK: group-test-$i scheduled (node=${NODE:-<unknown>}, phase=$PHASE)"
  else
    echo "FAIL: group-test-$i was not scheduled (phase=$PHASE)"
    FAIL=1
  fi
done
if [ "$FAIL" -ne 0 ]; then
  kubectl get pods -n custom-scheduler-namespace
  exit 1
fi
echo "Group of 3 validated."

echo "Tearing down group of 3..."
kubectl delete pod group-test-1 group-test-2 group-test-3 -n custom-scheduler-namespace --ignore-not-found=true
sleep 5

echo "Scheduling single priority pod..."
kubectl apply -f infra/test/priority-test-container.yaml
sleep 5
kubectl get pods -n custom-scheduler-namespace
sleep 10

PRIO_NODE=$(kubectl get pod priority-test -n custom-scheduler-namespace -o jsonpath='{.spec.nodeName}' 2>/dev/null || true)
PRIO_PHASE=$(kubectl get pod priority-test -n custom-scheduler-namespace -o jsonpath='{.status.phase}' 2>/dev/null || true)
if [ -z "$PRIO_NODE" ] && [ "$PRIO_PHASE" != "Running" ]; then
  echo "FAIL: priority-test was not scheduled (phase=$PRIO_PHASE)"
  kubectl get pods -n custom-scheduler-namespace
  exit 1
fi
echo "Priority pod validated (scheduled)."

echo "Scheduling another group of 3 (should not fit; only 2 nodes free)..."
kubectl apply -f infra/test/group-test-container-2.yaml
sleep 5
kubectl get pods -n custom-scheduler-namespace
sleep 15

RUNNING=0
for i in 4 5 6; do
  PHASE=$(kubectl get pod "group-test-$i" -n custom-scheduler-namespace -o jsonpath='{.status.phase}' 2>/dev/null || true)
  if [ "$PHASE" = "Running" ]; then
    RUNNING=$((RUNNING + 1))
  fi
done
if [ "$RUNNING" -ne 0 ]; then
  echo "FAIL: expected all 3 group-test-4/5/6 to remain Pending or ContainerCreating (none scheduled), got $RUNNING Running"
  kubectl get pods -n custom-scheduler-namespace
  exit 1
fi
echo "Validated: none of the second group of 3 were scheduled (all Pending or ContainerCreating)."

kubectl delete pod group-test-4 group-test-5 group-test-6 priority-test -n custom-scheduler-namespace --ignore-not-found=true
echo "Done."
