gen:
	go generate ./...

docker-build:
	docker buildx build . --build-arg APP_NAME=go-template -f docker/Dockerfile -t go-template

docker-run:
	docker run --name go-template -d -p 8080:8080 go-template

tests:
	go test -v  ./usecase/.