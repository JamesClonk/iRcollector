.PHONY: run dev build push setup glide test db start-db stop-db cleanup connect
SHELL := /bin/bash

all: run

run: build
	scripts/run.sh

dev: db
	scripts/dev.sh

build:
	rm -f iRcollector || true
	GOARCH=amd64 GOOS=linux go build -i -o iRcollector

push: test build
	cf push

setup:
	go get -v -u github.com/codegangsta/gin
	go get -v -u github.com/Masterminds/glide

glide:
	glide install --force

test:
	GOARCH=amd64 GOOS=linux go test $$(go list ./... | grep -v /vendor/)

db: stop-db start-db

start-db:
	docker run --name ircollector_db \
		-e POSTGRES_USER=dev-user \
		-e POSTGRES_PASSWORD=dev-secret \
		-e POSTGRES_DB=ircollector_db \
		-p "5432:5432" \
		-d postgres:9-alpine
	scripts/db.sh

stop-db:
	docker kill ircollector_db || true
	docker rm -f ircollector_db || true

cleanup:
	docker system prune --volumes -a

connect:
	docker exec -it ircollector_db psql -U dev-user -d ircollector_db
