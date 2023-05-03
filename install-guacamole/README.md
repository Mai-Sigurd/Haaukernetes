# Installing Guacamole

The script installs Guacamole, guacd, and postgres. 

## Requirements
- Running Kubernetes cluster
- Access to running kubectl commands

## Install
- Move the `init-guac.sh` script to a server with access to running kubectl commands.
- Run `chmod +x init-guac.sh` to make it executable.
- Run `./init-guac.sh`.
    - You will be asked to create a password for the postgres database and a username + password for the admin Guacamole user.

## Check Setup
- Run `kubectl get pods -n guacamole` to check the status of the components. 

## Connecting to Guacamole
- The script outputs the address for connecting to guacamole. 
- It can also be found by running `kubectl get svc -n guacamole` and using the form `http://<public-server-ip>:guacamole-exposed-nodeport/guacamole`.
- The default username and password is `guacadmin`.

## Manually Connecting to Kali via Guacamole Interface
- Make sure that there is a Kali container running. 
- Get the cluster IP and port for the Kali container using `kubectl get services`.
- Create a new connection in Guacamole using Kali's cluster IP and port. 
  - Protocol: RDP
  - Network -> hostname: Kali cluster IP
  - Network -> port: Kali port
  - Authentication -> username and password: Kali
  - Select "Ignore server certificate"

**Note:** To automatically create both a Kali, Guacamole user, and Guacamole connection use the `Post Kali` endpoint in the `src` API.

## Removing Guacamole Components
- Delete the namespace components using `kubectl delete namespace guacamole`.
- Delete the persistent volume using `kubectl delete pv postgres-pv-volume`.