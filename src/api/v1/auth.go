package v1

import (
	"net/http"
	"fmt"

	"nomad/api/src/resources"
	"nomad/api/src/models"
	"nomad/api/src/app/filters"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AuthController struct {
	Resources	*resources.Resources
}

var instanceAuthController *AuthController

func InitAuthController() *AuthController {
	if instanceAuthController == nil {
		res, _ := resources.GetInstance()

		instanceAuthController = &AuthController{Resources: res}
	}
	return instanceAuthController
}


type ParamsSignin struct {
	Email			string			`json:"email" description:"email of the user" valid:"required"`
	Password 		string			`json:"password" description:"password of the user" valid:"required"`
}

func (u *AuthController) WebService() *restful.WebService {
	ws := new(restful.WebService)
	tags := []string{"Auth"}
	ws.Path("/v1/auth").
		ApiVersion("v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
		//Param(ws.HeaderParameter("X-Access-Token", "Уникальный ключ пользователя"))

	ws.Route(ws.POST("/signin").To(u.Signin).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		//Param(ws.QueryParameter("email", "user email").DataType("string")).
		//Param(ws.QueryParameter("password", "user password").DataType("string")).
		//Reads(ParamsSignin{}).
		//Param(ws.BodyParameter("ParamsSignin", "params for signin").DataType("ParamsSignin")).
		Doc("Signin"))

	ws.Route(ws.GET("/refresh").To(u.Refresh).
		Filter(filters.ValidateJWT).
		Consumes(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Produces(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Param(ws.HeaderParameter("Authorization", "Bearer JWT token")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Refresh"))

	ws.Route(ws.GET("/signout").To(u.Signout).
		Filter(filters.ValidateJWT).
		Consumes(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Produces(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Param(ws.HeaderParameter("Authorization", "Bearer JWT token")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Signout"))

	ws.Route(ws.POST("/signup").To(u.Signup).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Reads(models.Users{}).
		Param(ws.BodyParameter("User", "user model").DataType("models.Users")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Signup"))

	//ws.Route(ws.POST("/signup/confirm").To(u.SignupConfirm).
	//	Metadata(restfulspec.KeyOpenAPITags, tags).
	//	Doc("Signup confirm"))

	return ws
}

func (u *AuthController) Signin(request *restful.Request, response *restful.Response) {

	params := ParamsSignin{}
	err := request.ReadEntity(&params)
	if err != nil {
		WriteErrorResponse(response, http.StatusNotAcceptable, err.Error())
		return
	}

	fmt.Printf("params: %+v\n", params)

	password := params.Password
	email := params.Email

	user := models.Users{}

	u.Resources.Log.Debug().Msgf("UserCollection().Find...")
	collection, session, err := u.Resources.Mongo.UserCollectionAndSession()
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return
	}
	defer session.Close()

	err = collection.Find(bson.M{"email": email}).One(&user)
	u.Resources.Log.Debug().Msgf("find error: %+v\n", err)
	if err != nil {
		if err == mgo.ErrNotFound {
			WriteErrorResponse(response, http.StatusBadRequest, "User not found")
		} else {
			u.Resources.Log.Debug().Msgf("find user error: %s", err.Error())
			WriteErrorResponse(response, http.StatusInternalServerError, err.Error())
		}
		return
	}

	u.Resources.Log.Debug().Msgf("CompareHashAndPassword starting...")
	err = user.CompareHashAndPassword(password)
	u.Resources.Log.Debug().Msgf("CompareHashAndPassword end")
	if err != nil {
		u.Resources.Log.Debug().Msgf("CompareHashAndPassword error: %s", err.Error())
		WriteErrorResponse(response, http.StatusUnauthorized, "Compare hash password error")
		return
	}

	tokenString, err := user.GetJWT()
	if err != nil {
		u.Resources.Log.Debug().Msgf("JWT error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "JWT error")
		return
	}

	u.Resources.Log.Debug().Msgf("result: %+v\n", user)

	response.ResponseWriter.Header().Set("Authorization", "Bearer "+tokenString)
	response.ResponseWriter.Header().Set("access-control-allow-headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	response.ResponseWriter.Header().Set("access-control-expose-headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	response.WriteAsJson(map[string]string{"token": tokenString})
}

func (u *AuthController) Refresh(request *restful.Request, response *restful.Response) {

	tokenString := request.Request.Context().Value("JwtToken");

	if tokenString == nil {
		WriteErrorResponse(response, http.StatusUnauthorized, "Not found JWT token")
		return
	}

	newTokenString, err := models.UpdateJWTToken(tokenString.(string))
	if err != nil {
		u.Resources.Log.Debug().Msgf("Parse JWT Token error: %s", err.Error())
		WriteErrorResponse(response, http.StatusUnauthorized, err.Error())
		return
	}
	response.WriteAsJson(map[string]string{"token": newTokenString})
}

func (u *AuthController) Signout(request *restful.Request, response *restful.Response) {
	WriteSuccessResponse(response, nil, "Logged out Successfully.")
}

func (u *AuthController) Signup(request *restful.Request, response *restful.Response) {

	user := models.Users{}
	err := request.ReadEntity(&user)
	if err != nil {
		WriteErrorResponse(response, http.StatusNotAcceptable, err.Error())
		return
	}

	_, err = user.GenerateHashPassword()
	if err != nil {
		u.Resources.Log.Fatal().Msg(err.Error())
		WriteErrorResponse(response, http.StatusNotAcceptable, err.Error())
		return
	}

	collection, session, err := u.Resources.Mongo.UserCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return

	}
	defer session.Close()
	err = collection.Insert(&user)
	if err != nil {
		if mgo.IsDup(err) {
			u.Resources.Log.Debug().Msgf("Email %s is duplicate", user.Email)
			WriteErrorResponse(response, http.StatusBadRequest, fmt.Sprintf("Email %s is duplicate", user.Email))
			return
		} else {
			u.Resources.Log.Debug().Msgf("Create user error: %s", err.Error())
			WriteErrorResponse(response, http.StatusInternalServerError, "")
			return
		}
	}

	response.WriteHeaderAndEntity(http.StatusCreated, nil)
}