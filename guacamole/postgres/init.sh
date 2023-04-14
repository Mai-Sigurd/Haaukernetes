#!/bin/bash

# Exit if any command fails
set -e

DB_PASSWORD=""

read -p "Enter password for db: " -s DB_PASSWORD
echo ""

echo "Creating guacamole namespace"
kubectl apply -f guacamole-namespace.yaml

echo "Creating K8 db secret"
kubectl create secret generic postgres \
    --from-literal postgres-user=guacamole \
    --from-literal postgres-password="$DB_PASSWORD" \
    --from-literal postgres-hostname=localhost \
    --from-literal postgres-database=guacamole \
    --from-literal postgres-port=5432 \
    --namespace=guacamole

echo "Create PVC and PV"
kubectl apply -f postgres-pvc-pv.yaml

echo "Create postgres deployment"
kubectl apply -f postgres-deployment.yaml

echo "Create postgres service"
kubectl apply -f postgres-service.yaml