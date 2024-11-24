include .env

build:
	go build -o main ./cmd/api/

start:
	./main

seed:
	go run ./cmd/ape/seed.go

restart: build start