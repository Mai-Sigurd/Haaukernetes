# Test Data and Graphs

Contains results from CPU and memory tests of different Haaukernetes scenarios and R files for creating graphs for the data. 

## Setup
The tests where performed on the following server setup:
- Control-plane: Ubuntu 22.04 (LTS) x64, 2 vCPUs, 4 GB memory, 80 GB disk
  - Running Kubernetes with Calico and monitoring stack. 
  - Tainted to ensure no test resources were scheduled on it. 
- Worker-node: Ubuntu 22.04 (LTS) x64, 4 vCPUs, 8 GB memory, 160 GB disk
  - Running Kubernetes with Calico and resources created by the tests. 
  - Untainted.

The data has been taken from the Grafana interface. 
