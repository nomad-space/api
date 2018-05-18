package v1

import (
	"fmt"
	"nomad/api/src/app/filters"
	"nomad/api/src/resources"
	"nomad/api/src/models"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

type UserController struct {
	Resources	*resources.Resources
}

var instanceUserController *UserController

func InitUserController() *UserController {
	if instanceUserController == nil {
		res, _ := resources.GetInstance()
		instanceUserController = &UserController{Resources: res}
	}
	return instanceUserController
}


type ParamsProfile struct {
	Email			string			`json:"email" description:"email of the user" valid:"required"`
	Password 		string			`json:"password" description:"password of the user" valid:"required"`
}

func (u *UserController) WebService() *restful.WebService {
	ws := new(restful.WebService)
	tags := []string{"User"}
	ws.Path("/v1/user").
		ApiVersion("v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Param(ws.HeaderParameter("Authorization", "Bearer JWT token"))

	ws.Route(ws.GET("/profile").To(u.Profile).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Filter(filters.ValidateJWT).
		//Consumes(restful.MIME_JSON).
		//Produces(restful.MIME_JSON).
		Doc("Profile"))

	return ws
}

func (u *UserController) Profile(request *restful.Request, response *restful.Response) {

	fmt.Printf("Profile action\n")

	JwtClaims := request.Request.Context().Value("JwtClaims")

	jwt := JwtClaims.(*models.JwtClaims)

	user := jwt.GetModel()

	WriteSuccessResponse(response, user, "")
}