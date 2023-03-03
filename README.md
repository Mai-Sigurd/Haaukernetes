# Haaukins-Kubernetes-bachelor-project


# Run 
minikube start
go run main.go


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

## noget med at challenges kører op en ip idresse, jeg ved ikke hvor. 

## Static checking 
Run `go install honnef.co/go/tools/cmd/staticcheck@latest` to install the checker
Run `staticcheck ./...` to check all packages. No output means that there are no problems.  