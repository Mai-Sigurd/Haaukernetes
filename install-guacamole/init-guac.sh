#!/bin/bash

# Exit if any command fails
set -e

SERVICE_NAME="postgres"
NAMESPACE="guacamole"
POSTGRES_IP=""
POSTGRES_PORT=5432
DB_PASSWORD=""
GUAC_PASSWORD=""
GUAC_USERNAME=""

while true; do
    read -s -p "Enter new password for postgres DB: " DB_PASSWORD
    echo
    read -s -p "Repeat password: " password2
    echo
    if [ "$DB_PASSWORD" = "$password2" ]; then
        break
    else
        echo "Passwords do not match. Please try again."
    fi
done

while true; do
    read -s -p "Enter new password for guacamole admin: " GUAC_PASSWORD
    echo
    read -s -p "Repeat password: " password2
    echo
    if [ "$GUAC_PASSWORD" = "$password2" ]; then
        break
    else
        echo "Passwords do not match. Please try again."
    fi
done

echo "##### Creating guacamole namespace"
kubectl create namespace $NAMESPACE

echo "##### Creating Guacamole secret"
kubectl create secret generic guacamole \
    --from-literal guac-user=$GUAC_USERNAME \
    --from-literal guac-password=$GUAC_PASSWORD \
    --namespace=guacamole

echo "##### Creating postgres secret"
kubectl create secret generic postgres \
    --from-literal postgres-user=guacamole \
    --from-literal postgres-password=$DB_PASSWORD \
    --from-literal postgres-database=guacamole \
    --from-literal postgres-port=$POSTGRES_PORT \
    --namespace=guacamole

echo "##### Setting up postgres"
kubectl apply -f postgres.yaml

# Get pod name for the new postgres pod
POSTGRES_POD=$(kubectl get pod --namespace=guacamole -l app=postgres -o jsonpath="{.items[0].metadata.name}")

echo "##### Waiting for postgres pod to be ready"
while true; do
    if kubectl wait --namespace=$NAMESPACE --for=condition=Ready pod/$POSTGRES_POD --timeout=300s; then
        echo "Pod $POD is ready"
        break
    else
        echo "Pod $POD is not ready yet. Retrying in 5 seconds..."
        sleep 3
    fi
done

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

# To avoid connection errors when running the database init script
sleep 10

echo "##### Run DB init script"
POSTGRES_CONNECTION_STRING="postgresql://guacamole:${DB_PASSWORD}@${POSTGRES_IP}:${POSTGRES_PORT}/guacamole"
kubectl exec -it "$POSTGRES_POD" -n guacamole -- psql "$POSTGRES_CONNECTION_STRING" < initdb.sql

echo "##### Updating guacamole secret with postgres IP"
POSTGRES_IP_ENCODED=$(echo -n "$POSTGRES_IP" | base64)
kubectl patch secret postgres -n guacamole --type='json' -p='[{"op": "add", "path": "/data/postgres-hostname", "value": "'$POSTGRES_IP_ENCODED'"}]'

echo "##### Setting up guacamole"
kubectl apply -f guacamole.yaml

sleep 20
echo "##### Changing default guacadmin password"
PUBLIC_IP=$(ip -f inet -o addr show eth0|cut -d\  -f 7 | cut -d/ -f 1 | head -n 1)
PORT=$(kubectl get service guacamole -n guacamole -o=jsonpath='{.spec.ports[0].nodePort}')

RESP=$(curl -X POST ${PUBLIC_IP}:${PORT}/guacamole/api/tokens -d "username=guacadmin&password=guacadmin&attributes=" -H "Content-Type: application/x-www-form-urlencoded")
TOKEN=$(echo $RESP | grep -o '"authToken":"[^"]*' | cut -d '"' -f4)

curl -X PUT ${PUBLIC_IP}:${PORT}/guacamole/api/session/data/postgresql/users/guacadmin/password?token=${TOKEN} -d '{"oldPassword":"guacadmin", "newPassword": "'${GUAC_PASSWORD}'"}' -H "Content-Type: application/json"

echo ""
echo "##### Installation complete!"


echo "You can access guacamole on ${PUBLIC_IP}:${PORT}/guacamole"
echo "The default username is: guacadmin which should be used with the password you just set"