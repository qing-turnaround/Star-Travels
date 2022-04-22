GOPATH:=$(shell go env GOPATH)
.PHONY: build
build:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web_app main.go

.PHONY: dockerBuild
dockerBuild:
	docker build -t zhugeqing/star-travels:latest dockerfile

.PHONY: dockerRun
dockerRun:
	docker run -d --net host -v ./conf/config.yaml:/app/conf/config.yaml zhugeqing/star-travels:latest

