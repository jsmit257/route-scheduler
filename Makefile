SHELL := /bin/bash

.PHONY: test
test: 
	go test -v -cover ./... 

.PHONY: build
build:
	go build -v ./cmd/main.go
	
.PHONY: run-local
run-local: build
	#python3 ./bin/evaluateShared.py --cmd "go run ./cmd/..." --problemDir data
	python3 ./bin/evaluateShared.py --cmd "./main" --problemDir data

.PHONY: count-total
count-total:
	./bin/shift-cost.sh

.PHONY: run-docker
run-docker: 
	docker-compose up --build --remove-orphans run-docker
