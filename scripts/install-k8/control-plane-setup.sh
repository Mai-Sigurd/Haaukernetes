#!/bin/bash

set -e

SERVER_IP=""
NODENAME=$(hostname -s)
POD_CIDR="192.168.0.0/16"

[[ $EUID -ne 0 ]] && echo "This script must be run as root" && exit 1
echo "##### Starting control-plane setup"
chmod +x common-setup.sh

while true; do
  echo "Input control-plane server ip:"
  read -r SERVER_IP
  if [[ $SERVER_IP =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    break
  fi
done

./common-setup.sh

echo "##### Initializing Kubeadm on control-plane node"
kubeadm init --apiserver-advertise-address=$IPADDR  --apiserver-cert-extra-sans=$IPADDR  --pod-network-cidr=$POD_CIDR --node-name $NODENAME --ignore-preflight-errors Swap

echo "##### Creating kubeconfig file"
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

echo "##### Installing calico"
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.0/manifests/calico.yaml

echo ""
echo "##### Control-plane node setup finished! To check server and calico setup run:"
echo "##### kubectl get po -n kube-system"

echo ""
echo "##### Now run worker-node-setup.sh on any worker node"
echo "##### Then join the worker node by running the following on the worker node as root:"
SETUP=$(kubeadm token create --print-join-command)
echo "$SETUP"

echo ""
echo "##### To check if the worker node has been added run:"
echo "##### kubectl get nodes"

