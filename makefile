gen:
	go generate ./...

docker-build:
	docker buildx build . --build-arg APP_NAME=tonfura-exercise -f docker/Dockerfile -t tonfura-exercise

docker-run:
	docker run --name tonfura-exercise -d -p 8080:8080 tonfura-exercise

tests:
	go test -v  ./usecase/.