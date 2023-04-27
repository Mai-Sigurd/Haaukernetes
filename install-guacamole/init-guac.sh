#!/bin/bash

# Exit if any command fails
set -e

SERVICE_NAME="postgres"
NAMESPACE="guacamole"
POSTGRES_IP=""
POSTGRES_PORT=5432
DB_PASSWORD=""

while true; do
    read -s -p "Enter new password for DB: " DB_PASSWORD
    echo
    read -s -p "Repeat password: " password2
    echo
    if [ "$DB_PASSWORD" = "$password2" ]; then
        break
    else
        echo "Passwords do not match. Please try again."
    fi
done

echo "##### Creating guacamole namespace"
kubectl create namespace $NAMESPACE

echo "##### Creating guacamole secret"
kubectl create secret generic guacamole \
    --from-literal postgres-user=guacamole \
    --from-literal postgres-password="$DB_PASSWORD" \
    --from-literal postgres-database=guacamole \
    --from-literal postgres-port=$POSTGRES_PORT \
    --namespace=guacamole

echo "##### Setting up postgres"
kubectl apply -f postgres.yaml

# Get pod name for the new postgres pod
POD=$(kubectl get pod --namespace=guacamole -l app=postgres -o jsonpath="{.items[0].metadata.name}")

echo "##### Waiting for postgres pod to be ready"
kubectl wait --namespace=guacamole --for=condition=Ready pod/"$POD"

echo "##### Run DB init script"
kubectl exec -it "$POD" -n guacamole -- psql -h localhost -d guacamole -U guacamole -p $POSTGRES_PORT < initdb.sql

echo "##### Waiting for $SERVICE_NAME service to become available"
while true; do
  # Check if postgres service exists
  if kubectl get service $SERVICE_NAME -n $NAMESPACE >/dev/null 2>&1; then
    # Retrieve service IP address
    POSTGRES_IP=$(kubectl get service $SERVICE_NAME -n $NAMESPACE -o jsonpath='{.spec.clusterIP}')
    if [[ ! -z "$POSTGRES_IP" ]]; then
      echo "##### Service $SERVICE_NAME is now available"
      break
    fi
  fi
  sleep 3
done

echo "##### Updating guacamole secret with postgres IP"
kubectl patch secret guacamole -n guacamole --type='json' -p='[{"op": "add", "path": "/data/postgres-hostname", "value": "'$POSTGRES_IP'"}]'

echo "##### Setting up guacamole"
kubectl apply -f guacamole.yaml

echo ""
echo "##### Installation complete!"

PUBLIC_IP=$(ip -f inet -o addr show eth0|cut -d\  -f 7 | cut -d/ -f 1 | head -n 1)
PORT=$(kubectl get service guacamole -n guacamole -o=jsonpath='{.spec.ports[0].nodePort}')

echo "You can access guacamole on ${PUBLIC_IP}:${PORT}/guacamole"
echo "The default username is: guacadmin"
echo "The default password is: guacadmin"
