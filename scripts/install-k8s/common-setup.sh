#!/bin/bash

set -e

OS="xUbuntu_20.04"
CRIO_VERSION="1.23"

while true; do
  read -p "Use default OS $OS for crio (y/n)" yn
  case $yn in
  [yY])
    break ;;
  [nN])
    read -p "Enter preferred OS" OS
    break
    ;;
  *) echo invalid response ;;
  esac
done

while true; do
  read -p "Use default crio version $CRIO_VERSION (y/n)" yn
  case $yn in
  [yY])
    break ;;
  [nN])
    read -p "Enter preferred version" CRIO_VERSION
    break
    ;;
  *) echo invalid response ;;
  esac
done

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

echo "##### Disabling swap"
swapoff -a
(
  crontab -l 2>/dev/null
  echo "@reboot /sbin/swapoff -a"
) | crontab - || true

echo "##### Installing crio runtime"
cat <<EOF | tee /etc/modules-load.d/crio.conf
overlay
br_netfilter
EOF

# Required sysctl parameters (persist across reboots)
cat <<EOF | tee /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

modprobe overlay
modprobe br_netfilter

cat <<EOF | tee /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

sysctl --system

echo "##### Enabling crio repositories for version $CRIO_VERSION"
cat <<EOF | tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/$OS/ /
EOF
cat <<EOF | tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable:cri-o:$CRIO_VERSION.list
deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable:/cri-o:/$CRIO_VERSION/$OS/ /
EOF

echo "##### Adding gpg keys"
curl -L https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable:cri-o:$CRIO_VERSION/$OS/Release.key | apt-key --keyring /etc/apt/trusted.gpg.d/libcontainers.gpg add -
curl -L https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/$OS/Release.key | apt-key --keyring /etc/apt/trusted.gpg.d/libcontainers.gpg add -

echo "##### Update and install crio and crio-tools"
apt-get update
apt-get install cri-o cri-o-runc cri-tools -y
systemctl daemon-reload
systemctl enable crio --now

echo "##### Installing Kubeadm, Kubelet, and Kubectl"
apt-get update
apt-get install -y apt-transport-https ca-certificates curl
curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list

apt-get update -y
apt-get install -y kubelet kubeadm kubectl

echo "##### Adding node IP to KUBELET_EXTRA_ARGS"
IFACE=eth1  # change to eth1 for DO's private network
DROPLET_IP_ADDRESS=$(ip addr show dev $IFACE | awk 'match($0,/inet (([0-9]|\.)+).* scope global/,a) { print a[1]; exit }')
cat > /etc/default/kubelet << EOF
KUBELET_EXTRA_ARGS=--node-ip=$DROPLET_IP_ADDRESS
EOF

systemctl daemon-reload
systemctl restart kubelet