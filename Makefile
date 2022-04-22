GOPATH:=$(shell go env GOPATH)
.PHONY: build
build:
	go build -o web_app main.go

.PHONY: dockerBuild
dockerBuild:
	docker build -t zhugeqing/star-travels:latest .

.PHONY: dockerRun
dockerRun:
	docker run -d --name web_app --net host -v ${PWD}:/app/ zhugeqing/star-travels:latest

