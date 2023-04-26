#!/bin/bash

# kubectl logs guacamole-7f86c55467-24442 -c guacamole -n guacamole
# Kig lige lidt på versions i de andre filer
# sammenlign med https://oopflow.medium.com/how-to-install-guacamole-on-kubernetes-7d747438c141

# Exit if any command fails
set -e

DB_PASSWORD=""

read -p "Enter password for db: " -s DB_PASSWORD
echo ""

echo "Creating guacamole namespace"
kubectl apply -f guacamole-namespace.yaml

# FIXX SÆT hostname til inside IP på postgres
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


# FIXX POD NAME ETC
kubectl exec -it postgres-7996cd45c5-7qpsw -n guacamole -- psql -h localhost -d guacamole -U guacamole -p 5432 < initdb.sql


echo "Creating guacamole deployment"
kubectl apply -f guacamole-deployment.yaml

echo "Creating guacamole service"
kubectl apply -f guacamole-service.yaml