Setup

Servers 
- Control plane: Ubuntu 22.04 (LTS) x64, 2 vCPUs, 4 GB memory, 80 GB disk
- Worker node: Ubuntu 22.04 (LTS) x64, 4 vCPUs, 8 GB memory, 160 GB disk

- Control plane with monitoring and the Go program.
  - Tainted so that nothing gets scheduled on it. 
  - Monitoring on it though because it uses some resources. 
- Worker node 
  - Untainted. 

### Note
- We do not use our API in the test. We call the funcs directly.

