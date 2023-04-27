#!/bin/bash

set -e

echo "##### Creating Grafana user"
read -p "Enter username: " username

while true; do
    read -s -p "Enter password: " password
    echo
    read -s -p "Repeat password: " password2
    echo
    if [ "$password" = "$password2" ]; then
        break
    else
        echo "Passwords do not match. Please try again."
    fi
done

# Install Helm if necessary
if command -v helm &> /dev/null
then
    echo "##### Helm is already installed"
else
    echo "##### Installing Helm"
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
    chmod 700 get_helm.sh
    ./get_helm.sh
fi

echo "##### Adding prometheus-community repo to Helm"
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add stable https://charts.helm.sh/stable
helm repo update

echo "##### Creating monitoring namespace"
kubectl create namespace monitoring

echo "##### Installing kube-prometheus stack and exposing on NodePort"
helm install kube-prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --set prometheus.service.nodePort=30000 \
  --set prometheus.service.type=NodePort \
  --set grafana.service.nodePort=31000 \
  --set grafana.service.type=NodePort \
  --set alertmanager.service.nodePort=32000 \
  --set alertmanager.service.type=NodePort \
  --set prometheus-node-exporter.service.nodePort=32001 \
  --set prometheus-node-exporter.service.type=NodePort \
  --set grafana.adminUser="$username" \
  --set grafana.adminPassword="$password" \

echo ""
echo "##### Finished Install!"
echo "##### Check status with kubectl get all -n monitoring"