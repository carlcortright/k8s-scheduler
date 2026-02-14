run: 
	go run ./cmd/main.go

start-minikube:
	minikube start
	echo "Minikube started, exposing kubernetes API on port 8080 locally for development"
	kubectl proxy --port=8080