# k8s-scheduler

A basic scheduler implementation for k8s, hand crafted with :heart: in Golang

## Scheduler design 

The design of our scheduler uses 

## Requirements

- Local docker installation 
- Go 1.24 for running locally
- Make (for shortccuts)

## Setup and Installation 

### Installing minikube (macos apple silicon)

`brew install minikube`

### Start Minikube

`make start-minikube`

### Deploy scheduler to your local cluster [todo]

`make deploy-local`

## Development 

The command `make start-minikube` exposes the minikube cluster locally. Use `make run` to run the scheduler to talk to this cluster on localhost:8080

# Useful Docs 

[Kube Scheduler](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/)