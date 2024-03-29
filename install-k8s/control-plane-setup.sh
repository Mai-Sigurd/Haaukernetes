#!/bin/bash

set -e

[[ $EUID -ne 0 ]] && echo "This script must be run as root" && exit 1
echo "##### Starting control-plane setup"

PUBLIC_IP=$(ip -f inet -o addr show eth0|cut -d\  -f 7 | cut -d/ -f 1 | head -n 1)
APISERVER=$(hostname -s)
NODENAME=$(hostname -s)
CIDR="10.244.0.0/16"

apt-get update -y
apt-get install -y apt-transport-https ca-certificates curl

echo "##### Disabling swap"
swapoff -a
(
  crontab -l 2>/dev/null
  echo "@reboot /sbin/swapoff -a"
) | crontab - || true


echo "##### Enabling bridged traffic on node"
cat <<EOF | tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

modprobe overlay
modprobe br_netfilter

# Required sysctl parameters (persist across reboots)
cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

# Apply parameters without reboot
sysctl --system

echo "##### Installing Containerd runtime"
curl -Lo /tmp/containerd-1.6.9-linux-amd64.tar.gz https://github.com/containerd/containerd/releases/download/v1.6.9/containerd-1.6.9-linux-amd64.tar.gz
tar Cxzvf /usr/local /tmp/containerd-1.6.9-linux-amd64.tar.gz

curl -Lo /tmp/runc.amd64 https://github.com/opencontainers/runc/releases/download/v1.1.4/runc.amd64
install -m 755 /tmp/runc.amd64 /usr/local/sbin/runc

curl -Lo /tmp/cni-plugins-linux-amd64-v1.1.1.tgz https://github.com/containernetworking/plugins/releases/download/v1.1.1/cni-plugins-linux-amd64-v1.1.1.tgz
mkdir -p /opt/cni/bin
tar Cxzvf /opt/cni/bin /tmp/cni-plugins-linux-amd64-v1.1.1.tgz

# Remove temporary files
rm /tmp/containerd-1.6.9-linux-amd64.tar.gz /tmp/runc.amd64 /tmp/cni-plugins-linux-amd64-v1.1.1.tgz

mkdir /etc/containerd
containerd config default | tee /etc/containerd/config.toml
sed -i 's/SystemdCgroup \= false/SystemdCgroup \= true/g' /etc/containerd/config.toml

curl -Lo /etc/systemd/system/containerd.service https://raw.githubusercontent.com/containerd/containerd/main/containerd.service
systemctl daemon-reload
systemctl enable --now containerd

echo "##### Installing Kubeadm, Kubelet, and Kubectl"
curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list

apt-get update
apt-get install -y kubelet kubeadm kubectl

echo "##### Initializing Kubeadm on control-plane node"

kubeadm init --apiserver-advertise-address=$PUBLIC_IP \
                  --apiserver-cert-extra-sans=$APISERVER \
                  --pod-network-cidr=$CIDR \
                  --node-name $NODENAME \
                  --ignore-preflight-errors Swap

echo "##### Creating kubeconfig file"
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

echo "##### Installing tigera operator and calico"
curl -Lo /tmp/tigera-operator.yaml https://raw.githubusercontent.com/projectcalico/calico/v3.24.5/manifests/tigera-operator.yaml
kubectl create -f /tmp/tigera-operator.yaml

# Update calico custom resource yaml to match pod cidr network
curl -Lo /tmp/custom-resources.yaml https://raw.githubusercontent.com/projectcalico/calico/v3.24.5/manifests/custom-resources.yaml
sed -i "s|192.168.0.0/16|$CIDR|" /tmp/custom-resources.yaml
kubectl create -f /tmp/custom-resources.yaml

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