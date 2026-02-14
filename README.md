# k8s-scheduler

A basic scheduler implementation for k8s, hand crafted with :heart: in Golang

## Scheduler design 

The design of our scheduler uses two background goroutines (internal/scheduler/nodes-listener.go and internal/scheduler/pods-listener.go) to poll the cluster to get the latest list of nodes and pods available to the scheduler. These maintain an accessible in-memory representation of the current state of the cluster. The scheduler acts as the main thread which implements the scheduling logic using the bind endpoint in kubernetes to bind new pods to their respective nodes while obeying the priority consideration

## Requirements

- Local docker installation 
- Go 1.24 for running locally
- Make (for shortcuts)

## Setup and Installation 

### Installing dependencies (untested, mac)

Install requirements on mac:

```bash
make setup
```

### Start Minikube

```bash
make start-minikube
```

### Deploy scheduler to your local cluster 

```bash
make deploy-scheduler
```

## Development 

This is for running the scheduler on your local machine and letting it talk to the local cluster via an exposed port. The command `make start-minikube` then `make expose-locally` exposes the minikube cluster locally. Use `make run` to run the scheduler to talk to this cluster on localhost:8080. 

## Run integration tests 

Before running integration tests deploy the scheduler to the cluster:

```bash
make deploy-scheduler
```

When complete with integration tests run the following to tear down the scheduler pod:

```bash
make remove-scheduler
```

### Basic integration test (scheduling a pod)

The following command will run an integration test which runs the scheduler, schedules a pod, confirms it's scheduled then tears everything down:

```bash
make basic-integration-test
```

### Priority integration test (removing lower-priorty pods)

The following command will schedule and confirm 3 basic pods with a lower priority, then wait, then schedule a pod with a higher priorty and confirm one of the basic pods was evicted:

```bash
make priority-integration-test
```

### Gang grouping integration test 

The following command will test gang grouping by first successfully scheduling a gang group, then creating a setup where a gang group shouldn't be scheduled and confirming it doesnt:

```bash
make group-integration-test
```

# Scheduling Retry



# Performance improvements



# Future improvements




# Useful Docs 

[Kube Scheduler](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/)