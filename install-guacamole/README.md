# Installing Guacamole

The script installs Guacamole, guacd, and postgres. 

## Requirements
- Running Kubernetes cluster
- Access to running kubectl commands

## Install
- Move the `init-guac.sh` script to a server with access to running kubectl commands.
- Run `chmod +x init-guac.sh` to make it executable.
- Run `./init-guac.sh`.
    - You will be asked to create a password for the postgres database.

## Check Setup
- Run `kubectl get pods -n guacamole` to check the status of the components. 

## Connecting to Guacamole
- The script outputs the address for connecting to guacamole. 
- It can also be found by running `kubectl get svc -n guacamole` and using the form `http://<public-server-ip>:guacamole-exposed-nodeport/guacamole`.
- The default username and password is `guacadmin`.

## Connecting to Kali via Guacamole Interface
- Make sure that there is a Kali container running. 
- Get the cluster IP and port for the Kali container using `kubectl get services`.
- Create a new connection in Guacamole using the cluster IP and port of the Kali. 
  - Protocol: RDP
  - Network -> hostname: Kali cluster IP
  - Network -> port: Kali port
  - Authentication -> username and password: "Kali"
  - Select "Ignore server certificate"