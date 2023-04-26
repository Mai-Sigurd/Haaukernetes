# Installing Prometheus and Grafana

The script installs the Kube-prometheus stack (including Grafana), which can be used for monitoring a Kubernetes cluster. 

## Requirements
- Running Kubernetes cluster
- Access to running kubectl commands 

Read more about the Helm chart [here](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)

## Install
- Move the `install-prometheus-grafana.sh` script to a server with access to running kubectl commands. 
- Run `chmod +x install-prometheus-grafana.sh` to make it executable.
- Run `./install-prometheus-grafana.sh`.
  - You will be asked to create a username and password for the Grafana dashboard.

## Check Setup
- Run `kubectl get all --namespace monitoring` to check the status of all components.

## Access Dashboards
The dashboards are available at:
- Prometheus: http://\<nodeIP>:30000
- Grafana: http://\<nodeIP>:31000
- AlertManager: http://\<nodeIP>:32000 
