# Haaukins-Kubernetes-bachelor-project

# Requirements
- Go-Swagger
  https://goswagger.io/install.html
-  ``go get github.com/gorilla/mux``
-  ``go get github.com/go-openapi/runtime/middleware``


# Run 
`` minikube start ``

``go run main.go``

# generating new swagger
## requirements
``go install github.com/swaggo/swag/cmd/swag@latest``

## run
``swag init ``

## Swagger UI can be seen at
localhost:5000/docs/index.html#/

# Frontend
Navigate to the commandline_frontend folder and run
``go run main.go``
