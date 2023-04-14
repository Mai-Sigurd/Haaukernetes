#!/bin/bash

# Exit if any command fails
set -e

DB_PASSWORD=""

read -p "Enter password for mariadb: " -s DB_PASSWORD
echo ""

echo "Creating guacamole namespace"
kubectl apply -f guacamole-namespace.yaml

echo "Creating K8 db secret"
kubectl create secret generic mariadb \
    --from-literal mariadb-user=guacamole \
    --from-literal mariadb-password="$DB_PASSWORD" \
    --from-literal mariadb-hostname=mariadb-host \
    --from-literal mariadb-database=guacamole \
    --from-literal mariadb-port=3306 \
    --namespace=guacamole

echo "Creating mariadb"
kubectl apply -f mariadb-deployment.yaml

# Get pod name for the new mariadb pod
POD=$(kubectl get pod --namespace=guacamole -l app=mariadb -o jsonpath="{.items[0].metadata.name}")

echo "Waiting for pod to be ready"
kubectl wait --namespace=guacamole --for=condition=Ready pod/"$POD"

echo "Create mariadb service"
kubectl apply -f mariadb-service.yaml

# THE BELOW DOES NOT WORK WHEN RUNNING AS A SCRIPT BUT WORKS WHEN RUNNING DIRECTLY

# echo "Setting up guacamole database"
# Create guacamole db
# kubectl exec --namespace=guacamole "$POD" -- mariadb -uroot -p"$DB_PASSWORD" -e "create database if not exists guacamole;"
# Create and configure the guacamole user
# kubectl exec --namespace=guacamole "$POD" -- mariadb -uroot -p"$DB_PASSWORD" -e "GRANT ALL ON guacamole.* TO 'guacamole'@'%' IDENTIFIED BY '$DB_PASSWORD';"
# Flush mariadb privileges
# kubectl exec --namespace=guacamole "$POD" -- mariadb -uroot -p"$DB_PASSWORD" -e "flush privileges;"
# Run initdb.sql script
# kubectl exec --namespace=guacamole "$POD" -- mariadb -uroot -p"$DB_PASSWORD" -Dguacamole < initdb.sql

echo "Creating guacamole deployment"
kubectl apply -f guacamole-deployment.yaml

echo "Creating guacamole service"
kubectl apply -f guacamole-service.yaml