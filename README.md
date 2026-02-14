# k8s-scheduler

A basic scheduler implementation for k8s, hand crafted with :heart:

## Stack

- Golang for efficency and parallel execution 
- Zap for logging in development and production (fast logging)
- k8s API calls via net/http and a custom client
- viper for configuration management 
- make for easy local development

## Scheduler design 

The scheduler has two main components: 
- Core loop on the polling interval which handles all of the business logic (internal/scheduler/scheduler.go) and maintains an in-memory map of both nodes and pods
- Two background goroutines (internal/scheduler/nodes-listener.go and internal/scheduler/pods-listener.go) to poll the cluster to get the latest list of nodes and pods available to the scheduler and sync in-memory maps of pods we're tracking synced with any new pods an external service adds to the namespace while unblocking expensive network calls

In addition it also has logging and configuration implemented for easy deployment and debugging. 

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

The scheduler object uses backoff retry for binding and evicting pods (see internal/scheduler/scheduler.go). This ensures that pods are correctly bound or evicted while not overwhelming internal APIs in the case of failure. 

# Performance improvements on large clusters

One potential improvement we could make with more time would be to instead of blocking on cluster updates (bind, evict) implement some sort of queueing mechanism that unblocks the main scheduler thread and tracks which pods are being actively updated. This way we could potentially bind / evict multiple pods at the same time in the case of a cluster that has thousands or 10s of thousands of pods. 


# Useful Docs 

[Kube Scheduler](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/)