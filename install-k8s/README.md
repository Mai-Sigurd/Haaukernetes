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
- Move the `control-plane-setup.sh` script to the server that will act as the control-plane node.
- Run `chmod +x control-plane-setup.sh` to make it executable.
- Run `./control-plane-setup.sh`.
- Save the `kubeadm join` command at the end of the script for the worker node setup. 

## Setup Worker Node(s)
- Move the `worker-node-setup.sh`script to the server that will act as the worker node.
- Run `chmod +x worker-node-setup.sh` to make it executable.
- Run `./worker-node-setup.sh`.
- Run the saved `kubeadm join` command from the control-plane setup to initialize and join the worker node.
- The process can be repeated for additional worker nodes. 
- The `kubeadm join` command contains a token that is only valid for 24 hours, so a new command can be generated with `kubeadm token create --print-join-command` 

## Check Setup
- Run `kubectl get pods -n kube-system` to check the general setup.
- Run `kubectl get nodes` to check if the worker node has been joined.

## Running kubectl on a Worker-Node
To be able to run `kubectl` commands on other nodes than the control-plane, they must have the kubeconfig file: 

- On the control-plane node run `cat $HOME/.kube/config` and copy the output.
- On a worker-node run `mkdir -p $HOME/.kube` and cd into the new directory.
- Create a `config` file and paste the output from the first step. 
- You should now be able to e.g. run `kubectl get nodes` on the worker-node.
