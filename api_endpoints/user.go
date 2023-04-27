package api_endpoints

import (
	"k8s-project/namespaces"

	"github.com/gin-gonic/gin"
)

type Users struct {
	// Users names
	// in: array
	Names []string `json:"names"`
}

type Challenges struct {
	// Challenges names
	// in: array
	Names []string `json:"names"`
}

// GetUser godoc
// @Summary Retrieves user based on given name
// @Produce json
// @Param user body User true "User"
// @Success 200 {object} User
// @Router /user/ [get]
func (c Controller) GetUser(ctx *gin.Context) {
	var body User
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		name := body.Name
		if !namespaces.NamespaceExists(*c.ClientSet, name) {
			message := "Sorry user " + name + " does not exist"
			ctx.JSON(400, ErrorResponse{Message: message})
		} else {
			ctx.JSON(200, User{name})
		}
	}
}

// GetUsers godoc
// @Summary Retrieves all users
// @Produce json
// @Success 200 {object} Users
// @Router /users [get]
func (c Controller) GetUsers(ctx *gin.Context) {
	result, err := namespaces.GetNamespaces(*c.ClientSet)
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	} else {

		ctx.JSON(200, Users{Names: result})
	}
}

// GetUserChallenges godoc
// @Summary Retrieves all challenges, as well as Kalis or wireguards running for a user
// @Produce json
// @Param user body User true "User"
// @Success 200 {object} Challenges
// @Router /user/challenges/ [get]
func (c Controller) GetUserChallenges(ctx *gin.Context) {
	var body User
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		name := body.Name
		result, err1 := namespaces.GetNamespacePods(*c.ClientSet, name)
		if err1 != nil {
			ctx.JSON(400, ErrorResponse{Message: err.Error()})
		}
		ctx.JSON(200, Challenges{Names: result})
	}
}

// PostUser godoc
// @Summary Creates user based on given name
// @Produce json
// @Param user body User true "User"
// @Success 200 {object} User
// @Router /user/ [post]
func (c Controller) PostUser(ctx *gin.Context) {
	var body User
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else if namespaces.NamespaceExists(*c.ClientSet, body.Name) {
		message := "Sorry user " + body.Name + " already exists"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		errKubernetes := namespaces.PostNamespace(*c.ClientSet, body.Name)
		if errKubernetes != nil {
			ctx.JSON(400, ErrorResponse{Message: err.Error()})
		}
		ctx.JSON(200, body)
	}
}

// DeleteUser godoc
// @Summary Deletes user based on given name
// @Produce json
// @Param user body User true "User"
// @Success 200
// @Router /user/ [delete]
func (c Controller) DeleteUser(ctx *gin.Context) {
	var body User
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		name := body.Name
		err1 := namespaces.DeleteNamespace(*c.ClientSet, name)
		if err1 != nil {
			ctx.JSON(400, ErrorResponse{Message: err.Error()})
		} else {
			ctx.JSON(200, "User "+name+" Deleted")
		}
	}
}
