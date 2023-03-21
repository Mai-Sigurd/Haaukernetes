package apis

import (
	"fmt"
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

// PostChallenge godoc
// @Summary Creates challenge based in a given namespace
// @Produce json
// @Param namespace body Challenge true "Challenge"
// @Success 200 {object} Challenge
// @Router /challenge/ [post]
func (c Controller) PostChallenge(ctx *gin.Context) {
	var body Challenge
	if err := ctx.BindJSON(&body); err != nil {
		bad := fmt.Errorf("bad request")
		_ = ctx.Error(bad)
	}
	podLabels := make(map[string]string)
	podLabels["app"] = body.ChallengeName
	podLabels["type"] = "challenge"
	deployments.CreateDeployment(*c.ClientSet, body.Namespace, body.ChallengeName, body.Port, podLabels)
	services.CreateService(*c.ClientSet, body.Namespace, body.ChallengeName, body.Port)

	ctx.JSON(200, body)
}

// DeleteChallenge godoc
// @Summary Deletes challenge in a namespace
// @Produce json
// @Param challenge body DelChallenge true "Challenge"
// @Success 200
// @Router /challenge/ [delete]
func (c Controller) DeleteChallenge(ctx *gin.Context) {
	var body DelChallenge
	if err := ctx.BindJSON(&body); err != nil {
		bad := fmt.Errorf("bad request")
		_ = ctx.Error(bad)
	}
	deleteChallenge(*c.ClientSet, body.Namespace, body.ChallengeName)
	ctx.JSON(200, body)
}

func deleteChallenge(clientSet kubernetes.Clientset, teamName string, challengeName string) {
	if !deployments.CheckIfDeploymentExists(clientSet, teamName, challengeName) {
		fmt.Printf("Challenge %s is not turned on \n", challengeName)
	} else {
		deploymentDeleteStatus := deployments.DeleteDeployment(clientSet, teamName, challengeName)
		serviceDeleteStatus := services.DeleteService(clientSet, teamName, challengeName)
		if deploymentDeleteStatus && serviceDeleteStatus {
			fmt.Printf("Challenge %s turned off\n", challengeName)
		} else {
			fmt.Printf("Challenge %s could not be turned off\n", challengeName)
		}
	}
}
