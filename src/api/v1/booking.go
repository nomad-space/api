package v1

import (
	"fmt"
	"net/http"
	"nomad/api/src/app/filters"
	"nomad/api/src/resources"
	"nomad/api/src/models"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

type BookingController struct {
	Resources	*resources.Resources
}

var instanceBookingController *BookingController

func InitBookingController() *BookingController {
	if instanceBookingController == nil {
		res, _ := resources.GetInstance()
		instanceBookingController = &BookingController{Resources: res}
	}
	return instanceBookingController
}

func (u *BookingController) WebService() *restful.WebService {
	ws := new(restful.WebService)
	tags := []string{"Booking"}
	ws.Path("/v1/booking").
		ApiVersion("v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Param(ws.HeaderParameter("Authorization", "Bearer JWT token"))

	ws.Route(ws.POST("/").To(u.Create).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Filter(filters.ValidateJWT).
		//Consumes(restful.MIME_JSON).
		//Produces(restful.MIME_JSON).
		Doc("Create"))

	ws.Route(ws.GET("/list").To(u.List).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Filter(filters.ValidateJWT).
		//Consumes(restful.MIME_JSON).
		//Produces(restful.MIME_JSON).
		Doc("List"))

	return ws
}

func (u *BookingController) Create(request *restful.Request, response *restful.Response) {

	booking := models.Booking{}
	err := request.ReadEntity(&booking)
	if err != nil {
		WriteErrorResponse(response, http.StatusNotAcceptable, err.Error())
		return
	}

	collection, session, err := u.Resources.Mongo.BookingCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return

	}
	defer session.Close()
	err = collection.Insert(&booking)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Create booking error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}

	response.WriteHeaderAndEntity(http.StatusCreated, nil)
}

func (u *BookingController) List(request *restful.Request, response *restful.Response) {

	booking := models.Booking{}
	err := request.ReadEntity(&booking)
	if err != nil {
		WriteErrorResponse(response, http.StatusNotAcceptable, err.Error())
		return
	}

	collection, session, err := u.Resources.Mongo.BookingCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return

	}
	defer session.Close()

	var results []models.Booking

	err = collection.Find(nil).All(&results)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Create booking error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}

	WriteSuccessResponse(response, results, "")
}