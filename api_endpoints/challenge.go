package api_endpoints

import (
	"github.com/gin-gonic/gin"
	"k8-project/challenge"
	"k8s.io/client-go/kubernetes"
)

type Challenge struct {
	Ports         []int32 `json:"ports"`
	ChallengeName string  `json:"challengeName"`
	Namespace     string  `json:"namespace"`
}

type DelChallenge struct {
	ChallengeName string `json:"challengeName"`
	Namespace     string `json:"namespace"`
}

type DelRespChallenge struct {
	ChallengeName string `json:"challengeName"`
	Namespace     string `json:"namespace"`
	Message       string `json:"message"`
}

// PostChallenge godoc
// @Summary Creates challenge based in a given namespace
// @Produce json
// @Param namespace body Challenge true "Challenge"
// @Success 200 {object} Challenge
// @Router /challenge/ [post]
func (c Controller) PostChallenge(ctx *gin.Context) {
	var body Challenge
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		challenge.CreateChallenge(*c.ClientSet, body.Namespace, body.ChallengeName, body.ChallengeName, body.Ports)
		ctx.JSON(200, body)
	}
}

// DeleteChallenge godoc
// @Summary Deletes challenge in a namespace
// @Produce json
// @Param challenge body DelChallenge true "Challenge"
// @Success 200 {object} DelRespChallenge
// @Router /challenge/ [delete]
func (c Controller) DeleteChallenge(ctx *gin.Context) {
	var body DelChallenge
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	}
	deleteChallenge(*c.ClientSet, body.Namespace, body.ChallengeName, ctx, body)
}

func deleteChallenge(clientSet kubernetes.Clientset, namespace string, challengeName string, ctx *gin.Context, body DelChallenge) {
	ctx.JSON(200, body)
	if !challenge.Challenge_exists(clientSet, namespace, challengeName) {
		message := "Challenge " + challengeName + " is not turned on"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		deploymentDeleteStatus, serviceDeleteStatus := challenge.DeleteChallenge(clientSet, namespace, challengeName)
		if deploymentDeleteStatus && serviceDeleteStatus {
			message := "Challenge " + challengeName + " turned off"
			resp := DelRespChallenge{
				ChallengeName: body.ChallengeName,
				Namespace:     body.ChallengeName,
				Message:       message,
			}
			ctx.JSON(200, resp)
		} else {
			message := "Challenge " + challengeName + "could not be turned off"
			ctx.JSON(400, ErrorResponse{Message: message})
		}
	}
}
