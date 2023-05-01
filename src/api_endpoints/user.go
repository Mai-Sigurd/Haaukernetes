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

type UserInfo struct {
	// UserInfo pods
	// in: array
	Pods []string `json:"pods"`
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

// GetUser godoc
// @Summary Retrieves all challenges, as well as Kalis or wireguards running for a user
// @Produce json
// @Param user body User true "User"
// @Success 200 {object} UserInfo
// @Router /user/ [get]
func (c Controller) GetUser(ctx *gin.Context) {
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
		ctx.JSON(200, UserInfo{Pods: result})
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
