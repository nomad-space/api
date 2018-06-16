package v1

import (
	//"fmt"
	"net/http"
	//"nomad/api/src/app/filters"
	"nomad/api/src/resources"
	"nomad/api/src/models"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"gopkg.in/mgo.v2/bson"
	"time"
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
		ApiVersion("v1")
		//Consumes(restful.MIME_JSON, restful.DefaultResponseMimeType).
		//Produces(restful.MIME_JSON, restful.DefaultResponseMimeType)
		//Param(ws.HeaderParameter("Authorization", "Bearer JWT token"))

	ws.Route(ws.POST("/").To(u.Create).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		//Filter(filters.ValidateJWT).
		Reads(models.Booking{}).
		Consumes(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Produces(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Doc("Create"))

	ws.Route(ws.GET("/list").To(u.List).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		//Filter(filters.ValidateJWT).
		//Consumes(restful.MIME_JSON).
		//Produces(restful.MIME_JSON).
		Doc("List"))

	ws.Route(ws.GET("/{booking_id}").To(u.Info).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("booking_id", "identifier of the booking").DataType("string")).
		//Filter(filters.ValidateJWT).
		//Consumes(restful.MIME_JSON).
		//Produces(restful.MIME_JSON).
		Doc("Info"))

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

	booking.Id = bson.NewObjectId()
	booking.UpdatedAt = time.Now()
	booking.CreatedAt = booking.UpdatedAt
	u.Resources.Log.Debug().Msgf("Before insert booking: %+v", booking)

	defer session.Close()
	err = collection.Insert(&booking)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Create booking error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}

	response.WriteHeaderAndEntity(http.StatusCreated, booking)
}

func (u *BookingController) List(request *restful.Request, response *restful.Response) {

	collection, session, err := u.Resources.Mongo.BookingCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return
	}
	defer session.Close()

	var results []models.Booking

	err = collection.Find(bson.M{}).All(&results)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Find booking error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}
	u.Resources.Log.Debug().Msgf("results: %s", results)

	WriteSuccessResponse(response, results, "")
}

func (u *BookingController) Info(request *restful.Request, response *restful.Response) {

	bookingId := request.PathParameter("booking_id")
	u.Resources.Log.Debug().Msgf("bookingId: %+v", bookingId)

	collection, session, err := u.Resources.Mongo.BookingCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return

	}
	defer session.Close()

	var results models.Booking

	err = collection.Find(bson.M{"_id": bson.ObjectIdHex(bookingId)}).One(&results)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Find booking error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}
	u.Resources.Log.Debug().Msgf("results: %s", results)

	WriteSuccessResponse(response, results, "")
}