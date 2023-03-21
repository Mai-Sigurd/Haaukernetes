// Package models Haaukins Api:
//
//	 version: 0.0.1
//	 title: Haaukins Api
//	Schemes: http, https
//	Host: localhost:5000
//	BasePath: /
//	Produces:
//	  - application/json
//
// swagger:meta
package models

// swagger:model Challenge
type Challenge struct {
	// Port of the challenge
	// in: int64
	Port int64 `json:"port"`
	// ImageName of the challenge docker image
	// in: string
	ImageName string `json:"imageName"`
	// Namespace where the challenge is running
	// in: string
	Namespace string `json:"namespace"`
}

// swagger:route POST /challenge/ addChallenge
// Create a new company
//
// security:
// - apiKey: []
// responses:
//
//	401: CommonError
//	200: GetCompany
func postChallenge() {

}

func deleteChallenge() {

}
