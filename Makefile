# Untested, best attempt at setup deps for mac
setup:
	brew install minikube
	brew install kubectl
	brew install docker
	brew install go@1.24

# Minikube commands
start-minikube:
	minikube start --nodes 3
	kubectl create namespace custom-scheduler-namespace
	echo "Minikube started, exposing kubernetes API on port 8080 locally for development"

expose-locally:
	kubectl proxy --port=8080

list-pods:
	kubectl get pods 

# Scheduler commands
deploy-scheduler:
	echo "Deploying custom scheduler to local cluster..."
	docker build -t custom-scheduler .
	minikube image load custom-scheduler:latest
	kubectl apply -f infra/scheduler/scheduler.yaml
	echo "Custom scheduler deployed to local cluster!"

remove-scheduler:
	echo "Removing custom scheduler from local cluster..."
	kubectl delete -f infra/scheduler/scheduler.yaml
	echo "Custom scheduler removed from local cluster!"

scheduler-logs:
	echo "Showing logs for custom scheduler..."
	kubectl logs -l app=custom-scheduler -n custom-scheduler

scheduler-describe:
	echo "Describing custom scheduler..."
	kubectl describe pod -l app=custom-scheduler -n custom-scheduler

# Integration tests
basic-integration-test:
	bash ./scripts/basic-integration-test.sh