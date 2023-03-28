#!/bin/bash

[[ $EUID -ne 0 ]] && echo "This script must be run as root/sudo" && exit 1
echo "Starting worker-node setup"
chmod +x common-setup.sh

./common-setup.sh

echo "Worker-node setup done!"
echo "Now join the node to the cluster by running the kubeadm join command from the control-plane setup"