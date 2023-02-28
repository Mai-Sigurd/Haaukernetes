#!/bin/bash

kubectl delete -f vnc_deploy.yaml
kubectl delete -f vnc_service.yaml
kubectl delete -f vnc_expose.yaml
