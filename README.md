# Haaukins-Kubernetes-bachelor-project

[![go-static-check](https://github.com/Mai-Sigurd/Haaukins-Kubernetes-bachelor-project/actions/workflows/go-static-check.yml/badge.svg)](https://github.com/Mai-Sigurd/Haaukins-Kubernetes-bachelor-project/actions/workflows/go-static-check.yml)

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

Swagger UI can be seen at: localhost:5000/docs/index.html#/

# Frontend
Navigate to the commandline_frontend folder and run
``go run main.go``
