# Running
## Without docker
Make sure that you have a valid k8s config file and export it's path to KUBECONFIG env variable, like so
``export KUBECONFIG=$HOME/.kube/config``
and that you have a valid digitalocean secret (in dockerconfig.json format - for image repository) in a DO_SECRET_PATH env variable, like so
``export DO_SECRET_PATH=$HOME/do_secret``
Run the program from the root directory
``go run main.go``


## With docker
Build the docker image from the provided Dockerfile
``docker build -t haaukins-revamp .``
Run the image, providing the k8s config and digitalocean secret through bind mounted volumes and exposing the app on the hosts port 33333
``docker run -v ~/.kube/config:/kube/config --env KUBECONFIG=/kube/config -v ~/do_secret:/secret/do_secret --env DO_SECRET_PATH=/secret/do_secret -p 33333:33333 -d haaukins-revamp``
OBS: this seems to cause issues with e.g. minikube as the k8s config contains several other paths that can't 
be resolved with the current ``docker run`` setup - consider just running it raw with minikube.

# Requirements
Go-Swagger
  https://goswagger.io/install.html

``go get github.com/gorilla/mux``

``go get github.com/go-openapi/runtime/middleware``

# Generating New Swagger
## Requirements
``go install github.com/swaggo/swag/cmd/swag@latest``

## Run
``swag init ``

Swagger UI can be seen at: URL/docs/index.html#/
Where URL is either localhost:33333 or public url:33333

# Frontend
Navigate to the commandline_frontend folder and run
``go run main.go``

Test
