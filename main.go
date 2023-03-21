package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"k8-project/models"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	//namespaceController := models.BaseHandlerNamespace()
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "docs1"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/docs1", sh1)

	namespace := r.PathPrefix("/namespace").Subrouter()
	namespace.HandleFunc("/", models.PostNamespace).Methods("POST")
	//namespace.HandleFunc("/", getNamespace).Methods("GET")
	//namespace.HandleFunc("/{id}", deleteNamespace).Methods("DELETE")

	challenge := r.PathPrefix("/challenge").Subrouter()
	challenge.HandleFunc("/", createChallenge).Methods("POST")

	http.Handle("/", r)
	s := &http.Server{
		Addr: fmt.Sprintf("%s:%s", "localhost", "5000"),
	}
	s.ListenAndServe()
}

func createChallenge(writer http.ResponseWriter, request *http.Request) {

}

func createNamespace(w http.ResponseWriter, request *http.Request) {
	response := models.Challenge{}
	response.Namespace = "hey"

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
