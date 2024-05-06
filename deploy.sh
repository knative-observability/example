#!/bin/bash

export REGISTRY=192.168.103.220/library

kubectl get namespace example || kubectl create namespace example
kubectl get broker default -n example || kn broker create default -n example

func deploy -n example -r $REGISTRY --path ./func/receive-order
func deploy -n example -r $REGISTRY --path ./func/restock
func deploy -n example -r $REGISTRY --path ./func/update-stock
func deploy -n example -r $REGISTRY --path ./func/update-stock # Revision 2
func deploy -n example -r $REGISTRY --path ./func/notify-merchant
func deploy -n example -r $REGISTRY --path ./func/payment
func deploy -n example -r $REGISTRY --path ./func/validate-order
func deploy -n example -r $REGISTRY --path ./func/generate-invoice
func deploy -n example -r $REGISTRY --path ./func/notify-user

kubectl create -f resources/triggers.yaml
kubectl create -f resources/sequence.yaml
kubectl create -f resources/redis.yaml
kubectl create -f resources/client.yaml

kn service update update-stock --traffic update-stock-00001=70 --traffic update-stock-00002=30 -n example