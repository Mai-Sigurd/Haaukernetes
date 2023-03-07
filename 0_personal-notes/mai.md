
# Mai
## build challenges
find web mappen og cat readmes i admin-logon og heartbleed
og kør docker kommando

## minikube
minikube start

## load challenges ind i minikube
minikube image load logon

## Go run main.go
go run main.go
det kører uden videre, men det kører jo rent faktisk ikke, får ikke nogle print statements

## delete deployment

    kubectl delete deployment haaukins-deployment

## se minikube dashboard
kør i terminal
minikube dashboard
alt er grønt !!

## open vnc
minikube service vnc-expose