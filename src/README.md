# Running
 
 ## Cluster specific, required constants
- Needs to be set in ``utils/const.go``
- ImageRepoUrl (url for the image repository to be used)
- WireguardEndpoint (the public endpoint of the cluster)
- WireguardSubnet (the K8s pod or service CIDR)

## Without docker

**Requirements** 

- Go installed
- Docker installed
- Valid K8s config file 
- Valid DigitalOcean image repository secret (for getting Docker images from the cloud registry). This should be in `dockerconfig.json` format

**Steps**
- Export the K8s config file path to a `KUBECONFIG` env variable with ``export KUBECONFIG=$HOME/.kube/config``
- Export the DO secret file path a `DO_SECRET_PATH` env variable with ``export DO_SECRET_PATH=$HOME/do_secret``
- Export the server IP to a `SERVER_IP` env variable with ``export SERVER_IP=<YOUR-SERVER-IP>``
- Run the program from the root directory with ``go run main.go``

## With docker
- Build the docker image from the provided Dockerfile with ``docker build -t haaukernetes .``
- Run the image, providing the k8s config and digitalocean secret through bind mounted volumes and exposing the app on the hosts port 33333:
    - The local paths provided for the -v flags must point to the actual files on your system. These can vary depending on how your cluster is set up
    - The IP of your server can be inserted directly in the SERVER_IP env variable instead of the bash command that fetches it dynamically 

```bash
docker run \
--env SERVER_IP=$(hostname -I | awk '{print $1}') \
-v ~/.kube/config:/kube/config --env KUBECONFIG=/kube/config \
-v ~/do_secret:/secret/do_secret --env DO_SECRET_PATH=/secret/do_secret \
-p 33333:33333 -d haaukernetes
```

OBS: this seems to cause issues with e.g. minikube as the k8s config contains several other paths that can't 
be resolved with the current ``docker run`` setup - consider just running it raw with minikube.

# Requirements
Go-Swagger: https://goswagger.io/install.html

Install with:

``go get github.com/gorilla/mux``

``go get github.com/go-openapi/runtime/middleware``

# Generating New Swagger
## Requirements
``go install github.com/swaggo/swag/cmd/swag@latest``

## Initialize Swagger
``swag init ``

The Swagger UI is available at: 

http://\<YOUR-SERVER-IP\>:33333/docs/index.html#/

# Frontend
*The frontend allows you to interact with the Haaukernetes API via the command-line.*

Navigate to the `commandline_frontend` folder and run
``go run main.go``

# Static Checker
Run `staticcheck .` inside `src`
