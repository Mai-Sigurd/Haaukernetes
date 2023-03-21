package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Challenge struct {
	Port      int64  `json:"port"`
	ImageName string `json:"imageName"`
	Namespace string `json:"namespace"`
}

type DelChallenge struct {
	ImageName string `json:"imageName"`
	Namespace string `json:"namespace"`
}

// PostChallenge godoc
// @Summary Creates challenge based in a given namespace
// @Produce json
// @Param namespace body Challenge true "Challenge"
// @Success 200 {object} Challenge
// @Router /challenge/ [post]
func PostChallenge(ctx *gin.Context) {
	// TODO
	var challengeBody Challenge
	if err := ctx.BindJSON(&challengeBody); err != nil {
		// TODO
	}
	fmt.Println(challengeBody)
	ctx.JSON(200, challengeBody)
}

// DeleteChallenge godoc
// @Summary Deletes challenge in a namespace
// @Produce json
// @Param challenge body DelChallenge true "Challenge"
// @Success 200
// @Router /challenge/ [delete]
func DeleteChallenge(ctx *gin.Context) {
	// TODO
	var challengeBody DelChallenge
	if err := ctx.BindJSON(&challengeBody); err != nil {
		//TODO
	}
	fmt.Println(challengeBody)
}
