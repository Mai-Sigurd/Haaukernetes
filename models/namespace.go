package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// swagger:model Namespace
type Namespace struct {
	// Namespace name
	// in: string
	Name string `json:"name"`
}

type ReqNamespace struct {
	// name: Name
	// in: formData
	// type: string
	// required: true
	Name string `json:"name" validate:"required,min=2,max=100,alpha_space"`
}

// swagger:route POST /namespace/  addNamespace
// Create a new Namespace
// responses:
//
//	401: CommonError
//	200: GetNamespace
func PostNamespace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Namespace{}
	response.Name = "hello"

	decoder := json.NewDecoder(r.Body)
	var reqNamespace ReqNamespace
	err := decoder.Decode(&reqNamespace)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler("invalid_requuest"))
		return
	}

	json.NewEncoder(w).Encode(response)
}
