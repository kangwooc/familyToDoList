#!/bin/bash
GOOS=linux go build
docker build -t kangwooc/final .
go clean