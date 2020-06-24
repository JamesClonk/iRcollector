.DEFAULT_GOAL := run
SHELL := /bin/bash
APP ?= $(shell basename $$(pwd) | tr '[:upper:]' '[:lower:]')
COMMIT_SHA = $(shell git rev-parse HEAD)

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: run
## run: runs the application
run: build
	scripts/run.sh

.PHONY: dev
## dev: sets up postgres and runs the application
dev:
	scripts/dev.sh

.PHONY: build
## build: builds the application
build: clean
	@echo "Building binary ..."
	go build -o ${APP}

.PHONY: clean
## clean: cleans up binary files
clean:
	@echo "Cleaning up ..."
	@go clean

.PHONY: test
## test: runs go test with the race detector
test:
	@source .env; GOARCH=amd64 GOOS=linux go test -v -race ./...

.PHONY: init
## init: sets up go modules
init:
	@echo "Setting up modules ..."
	@go mod init 2>/dev/null; go mod tidy && go mod vendor

.PHONY: push
## push: pushes the application onto CF
push: test build
	cf push

.PHONY: setup
## setup: downloads gin
setup:
	go get -v -u github.com/codegangsta/gin

.PHONY: db
## db: runs postgres backend on docker
db: stop-db start-db connect

.PHONY: start-db
start-db:
	docker run --name ircollector_db \
		-e POSTGRES_USER=dev-user \
		-e POSTGRES_PASSWORD=dev-secret \
		-e POSTGRES_DB=ircollector_db \
		-p "5432:5432" \
		-d postgres:9-alpine
	scripts/db.sh

.PHONY: stop-db
## stop-db: cleans up postgres backend
stop-db:
	docker kill ircollector_db || true
	docker rm -f ircollector_db || true

.PHONY: cleanup
## cleanup: cleans up local docker images and volumes
cleanup: docker-cleanup
.PHONY: docker-cleanup
docker-cleanup:
	docker system prune --volumes -a

.PHONY: connect
## connect: connects to postgres backend with CLI
connect:
	docker exec -it ircollector_db psql -U dev-user -d ircollector_db

########################################################################################################################
####### docker/kubernetes related stuff ################################################################################
########################################################################################################################
.PHONY: image-login
## image-login: login to docker hub
image-login:
	@export PATH="$$HOME/bin:$$PATH"
	@echo $$DOCKER_PASS | docker login -u $$DOCKER_USER --password-stdin

.PHONY: image-build
## image-build: build docker image
image-build: build
	@export PATH="$$HOME/bin:$$PATH"
	docker build -t jamesclonk/${APP}:${COMMIT_SHA} .

.PHONY: image-publish
## image-publish: build and publish docker image
image-publish:
	@export PATH="$$HOME/bin:$$PATH"
	docker push jamesclonk/${APP}:${COMMIT_SHA}
	docker tag jamesclonk/${APP}:${COMMIT_SHA} jamesclonk/${APP}:latest
	docker push jamesclonk/${APP}:latest

.PHONY: image-run
## image-run: run docker image
image-run:
	@export PATH="$$HOME/bin:$$PATH"
	docker run --rm --env-file .dockerenv -p 9090:9090 jamesclonk/${APP}:${COMMIT_SHA}
