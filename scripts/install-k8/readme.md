# Installing Kubernetes on Bare Metal

The scripts install Kubernetes with Kubeadm and the Calico network plugin. One script is used to initialize the control-plane node and the other can be used to join worker-nodes. 

## Requirements
The server needs to have the required resources for kubeadm:
- Compatible Linux host
- 2 GB or more RAM per machine 
- 2 CPUs or more
- Full network connectivity between all machines in the cluster
- Unique hostname, MAC address, and product_uuid for every node
- Certain ports are open (see link below)
- Swap disabled (also happens during the script)

Read more [here](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)

## Setup Control-Plane Node
- Move the `control-plane-setup.sh` and `common-setup.sh` scripts to the server that will act as the control-plane node.
- Run `chmod +x control-plane-setup.sh` to make it executable.
- Run `./control-plane-setup.sh` (this script will also run `common-setup.sh`).
  - You will be asked to enter the server's IP and whether you want to use the default OS and crio version.
- Save the `kubeadm join` command at the end of the script for the worker node setup. 

## Setup Worker Node(s)
- Move the `worker-node-setup.sh` and `common-setup.sh` scripts to the server that will act as the worker node.
- Run `chmod +x worker-node-setup.sh` to make it executable.
- Run `./worker-node-setup.sh` (this script will also run `common-setup.sh`)
  - Use the same (default) OS and crio version as the control-plane.
- Run the saved `kubeadm join` command from the control-plane setup to initialize and join the worker node.
- The process can be repeated for additional worker nodes. 

## Check Setup
- Run `kubectl get po -n kube-system` to check the general setup.
- Run `kubectl get nodes` to check if the worker node has been joined.