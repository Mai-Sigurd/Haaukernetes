package apis

import (
	"github.com/gin-gonic/gin"
	"k8-project/deployments"
	"k8-project/services"
	"k8s.io/client-go/kubernetes"
)

type Challenge struct {
	Port          int32  `json:"port"`
	ChallengeName string `json:"challengeName"`
	Namespace     string `json:"namespace"`
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
		podLabels := make(map[string]string)
		podLabels["app"] = body.ChallengeName
		podLabels["type"] = "challenge"
		deployments.CreateDeployment(*c.ClientSet, body.Namespace, body.ChallengeName, body.Port, podLabels)
		services.CreateService(*c.ClientSet, body.Namespace, body.ChallengeName, body.Port)

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

func deleteChallenge(clientSet kubernetes.Clientset, teamName string, challengeName string, ctx *gin.Context, body DelChallenge) {
	ctx.JSON(200, body)
	if !deployments.CheckIfDeploymentExists(clientSet, teamName, challengeName) {
		message := "Challenge " + challengeName + " is not turned on"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		deploymentDeleteStatus := deployments.DeleteDeployment(clientSet, teamName, challengeName)
		serviceDeleteStatus := services.DeleteService(clientSet, teamName, challengeName)
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
