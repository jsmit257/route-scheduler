SHELL := /bin/bash

.PHONY: test
test: 
	go test -v -cover ./...

.PHONY: run-local
run-local: 
	python3 ./bin/evaluateShared.py --cmd "go run ./cmd/..." --problemDir data

.PHONY: run-docker
run-docker: 
	docker-compose up --build --remove-orphans run-docker
