#!/bin/bash

REGISTRY='192.168.103.220/library'
go build -o ./build/
docker build -t $REGISTRY/example-client:latest .
docker push $REGISTRY/example-client:latest
