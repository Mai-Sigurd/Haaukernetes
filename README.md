[![go-static-check](https://github.com/Mai-Sigurd/Haaukernetes/actions/workflows/go-static-check.yml/badge.svg)](https://github.com/Mai-Sigurd/Haaukernetes/actions/workflows/go-static-check.yml)
![wordart](https://user-images.githubusercontent.com/63573363/236185676-91977a5f-8378-4edf-b3b9-46364418eb30.png)

# Haaukernetes: Haaukins Kubernetes Bachelor Project

Haaukernetes simulates parts of the [Haaukins](https://docs.haaukins.com/) CTF platform created by Aalborg University but using Kubernetes.
After setting up a Kubernetes cluster, adding Apache Guacamole, and running the Go program, the user can interact with an API to create users that can access CTF challenges via a VPN or in their browser.
The project consists of different parts (details about setup are available in the respective READMEs):

### install-k8s
Scripts for installing a Kubernetes cluster with the Calico network plugin. One script for setting up the control plane and one for setting up worker nodes.

*A functioning cluster is required for the program to run. However, it can be setup via other methods as long as the program has access to the kubeadm config file.*

### install-guacamole
Scripts for installing [Apache Guacamole](https://guacamole.apache.org/), a remote desktop gateway, on the cluster. 

*Guacamole is required if the program should be able to handle browser Kali browser connections.*

### install-monitoring
Scripts for installing the [Kube-prometheus stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack) (including Grafana) on the cluster, which can be used for monitoring the cluster.

*Monitoring is not required for the program to run.*

### report
Files related to the project report: data analysis R scripts, test data, and graphs based on the data. 

### src
Program written in Go that can be used to create CTF users and connections. It exposes an API and uses the [Kubernetes go client](https://github.com/kubernetes/client-go/) to communicate with the cluster.

The directory also includes tests for estimating the resource use of different container setups in terms of amount of users, connections, and challenges. 
