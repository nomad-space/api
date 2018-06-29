package v1

import (
	"net/http"
	//"nomad/api/src/app/filters"
	"nomad/api/src/resources"
	"nomad/api/src/models"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strconv"
)

type LocationController struct {
	Resources	*resources.Resources
}

var instanceLocationController *LocationController

func InitLocationController() *LocationController {
	if instanceLocationController == nil {
		res, _ := resources.GetInstance()
		instanceLocationController = &LocationController{Resources: res}
	}
	return instanceLocationController
}

func (u *LocationController) WebService() *restful.WebService {
	ws := new(restful.WebService)
	tags := []string{"Location"}
	ws.Path("/v1/location").
		ApiVersion("v1")
		//Consumes(restful.MIME_JSON, restful.DefaultResponseMimeType).
		//Produces(restful.MIME_JSON, restful.DefaultResponseMimeType)
		//Param(ws.HeaderParameter("Authorization", "Bearer JWT token"))

	ws.Route(ws.POST("/").To(u.Create).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		//Filter(filters.ValidateJWT).
		Reads(models.Location{}).
		Consumes(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Produces(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Doc("Create"))

	ws.Route(ws.GET("/list").To(u.List).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		//Filter(filters.ValidateJWT).
		//Consumes(restful.MIME_JSON).
		//Produces(restful.MIME_JSON).
		Doc("List"))

	ws.Route(ws.GET("/{location_id}").To(u.Info).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("location_id", "identifier of the location").DataType("integer")).
		//Filter(filters.ValidateJWT).
		//Consumes(restful.MIME_JSON).
		//Produces(restful.MIME_JSON).
		Doc("Info"))

	return ws
}

func (u *LocationController) Create(request *restful.Request, response *restful.Response) {

	location := models.Location{}
	err := request.ReadEntity(&location)
	if err != nil {
		WriteErrorResponse(response, http.StatusNotAcceptable, err.Error())
		return
	}

	collection, session, err := u.Resources.Mongo.LocationCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return

	}

	location.UpdatedAt = time.Now()
	location.CreatedAt = location.UpdatedAt
	u.Resources.Log.Debug().Msgf("Before insert location: %+v", location)

	defer session.Close()
	_, err = collection.Upsert(bson.M{"id": location.Id}, bson.M{"$set": &location})
	if err != nil {
		u.Resources.Log.Debug().Msgf("Create location error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}

	response.WriteHeaderAndEntity(http.StatusCreated, location)
}

func (u *LocationController) List(request *restful.Request, response *restful.Response) {

	collection, session, err := u.Resources.Mongo.LocationCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return
	}
	defer session.Close()

	var results []models.Location

	err = collection.Find(bson.M{}).All(&results)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Find location error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}
	u.Resources.Log.Debug().Msgf("results: %s", results)

	WriteSuccessResponse(response, results, "")
}

func (u *LocationController) Info(request *restful.Request, response *restful.Response) {

	locationId, err := strconv.Atoi(request.PathParameter("location_id"))
	if err != nil {
		u.Resources.Log.Debug().Msgf("Convert param locationId failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Params failed")
		return
	}

	u.Resources.Log.Debug().Msgf("locationId: %#v", locationId)

	collection, session, err := u.Resources.Mongo.LocationCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return
	}
	defer session.Close()

	var results models.Location

	err = collection.Find(bson.M{"id": int(locationId)}).One(&results)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Find location error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}
	u.Resources.Log.Debug().Msgf("results: %s", results)

	WriteSuccessResponse(response, results, "")
}