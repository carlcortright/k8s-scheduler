run: 
	go run ./cmd/main.go

start-minikube:
	minikube start
	kubectl create namespace custom-scheduler-namespace
	echo "Minikube started, exposing kubernetes API on port 8080 locally for development"
	kubectl proxy --port=8080

deploy-scheduler:
	bash ./scripts/deploy-scheduler.sh

integration-test:
	bash ./scripts/integration-test.sh