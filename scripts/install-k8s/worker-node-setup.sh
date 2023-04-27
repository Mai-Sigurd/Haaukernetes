#!/bin/bash

set -e

[[ $EUID -ne 0 ]] && echo "This script must be run as root/sudo" && exit 1
echo "##### Starting worker-node setup"

apt-get update
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
sudo systemctl daemon-reload
systemctl enable --now containerd

echo "##### Installing Kubeadm, Kubelet, and Kubectl"
curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list

apt-get update
apt-get install -y kubelet kubeadm kubectl

echo "##### Worker-node setup done!"
echo "##### Join the node to the cluster by running the kubeadm join command from the control-plane setup"