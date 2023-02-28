#!/bin/bash

kubectl apply -f vnc_deploy.yaml
kubectl apply -f vnc_service.yaml
kubectl apply -f vnc_expose.yaml
