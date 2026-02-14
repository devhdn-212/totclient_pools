running:
	nodemon --exec go run main.go --signal SIGTERM
build:
	docker-compose up -d --build