#based on example from https://hub.docker.com/_/golang
FROM golang:1.20

WORKDIR /usr/src/haaukins-revamp

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/haaukins-revamp .

# Generate and serve Swagger docs
RUN go get github.com/gorilla/mux
RUN go get github.com/go-openapi/runtime/middleware
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init


CMD ["haaukins-revamp"]
