include .env

SHELL := $(shell command -v bash;)
PROJECTNAME := $(shell basename "$(PWD)")

all: help

## build: Build the application
.PHONY: build
build:
	go build -o=./bin/${PROJECTNAME} ${PROJECTNAME}.go

## run: Build and run the application
.PHONY: run
run:
	./bin/${PROJECTNAME} --log-level ${LOG_LEVEL} --path ${PATH} --log-file ${LOG_FILE} \
		--s3.access-key ${AWS_ACCESS_KEY} --s3.secret-key ${AWS_SECRET_KEY} --s3.region ${AWS_REGION} \
		--s3.bucket ${AWS_S3_BUCKET} --domain ${DOMAIN_NAME} --workers-count ${WORKERS} --dsn ${DSN}

## help: Help about any command
.PHONY: help
help: Makefile
	@echo "Available Commands in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
