
# Mai

## minikube
``minikube start``

## minikube og nye docker images
``eval $(minikube -p minikube docker-env)``

i den mappe hvor din dockerfile er
``docker build . -t <image name>``

check dit image er i minikube
``minikube image ls --format table``

## load challenges ind i minikube hvis ovenstående ikke virker
``minikube image load logon``

## Go run main.go
``go run main.go``

## delete deployment

    kubectl delete deployment haaukins-deployment

## se minikube dashboard

``minikube dashboard``

## open vnc
``minikube service vnc-expose``