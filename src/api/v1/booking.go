package v1

import (
	"net/http"
	"time"
	"encoding/json"
	"nomad/api/src/resources"
	"nomad/api/src/models"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"gopkg.in/mgo.v2/bson"
	"nomad/api/src/sendmail/templates"
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

	ws.Route(ws.PATCH("/{booking_id}").To(u.Update).
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
	booking.Status = models.BOOKING_STATUS_NEW
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

	var bookingObj models.Booking

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"_id": booking.Id}},
		bson.M{"$lookup": bson.M{ "from": "hotels", "localField": "hotel_id", "foreignField": "id", "as": "hotel"}},
		bson.M{"$lookup": bson.M{ "from": "locations", "localField": "location_id", "foreignField": "id", "as": "location"}},
	}

	err = collection.Pipe(pipeline).One(&bookingObj)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Pipeline failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Pipeline failed")
		return
	}

	u.Resources.Log.Debug().Msgf("bookingObj: %+v", bookingObj)

	err = u.Resources.Mail.Send( booking.Email,
		"Booking is Processed!",
		templates.Processed(bookingObj))
	if err != nil {
		u.Resources.Log.Panic().Msgf("Error send mail to %s: %s", booking.Email, err.Error())
	}

	response.WriteHeaderAndEntity(http.StatusCreated, bookingObj)
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


	pipeline := []bson.M{
		bson.M{"$lookup": bson.M{ "from": "hotels", "localField": "hotel_id", "foreignField": "id", "as": "hotel"}},
		bson.M{"$lookup": bson.M{ "from": "locations", "localField": "location_id", "foreignField": "id", "as": "location"}},
	}
	err = collection.Pipe(pipeline).All(&results)

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

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(bookingId)}},
		bson.M{"$lookup": bson.M{ "from": "hotels", "localField": "hotel_id", "foreignField": "id", "as": "hotel"}},
		bson.M{"$lookup": bson.M{ "from": "locations", "localField": "location_id", "foreignField": "id", "as": "location"}},
	}

	err = collection.Pipe(pipeline).One(&results)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Pipeline failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Pipeline failed")
		return
	}

	WriteSuccessResponse(response, results, "")
}

func (u *BookingController) Update(request *restful.Request, response *restful.Response) {

	decoder := json.NewDecoder(request.Request.Body)
	var params interface{}
	err := decoder.Decode(&params)
	if err != nil {
		panic(err)
	}
	defer request.Request.Body.Close()
	u.Resources.Log.Debug().Msgf("params: %+v", params)

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

	err = collection.Update(bson.M{"_id": bson.ObjectIdHex(bookingId)}, bson.M{"$set": params})
	if err != nil {
		u.Resources.Log.Debug().Msgf("Update booking error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}
	u.Resources.Log.Debug().Msgf("results: %s", results)

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(bookingId)}},
		bson.M{"$lookup": bson.M{ "from": "hotels", "localField": "hotel_id", "foreignField": "id", "as": "hotel"}},
		bson.M{"$lookup": bson.M{ "from": "locations", "localField": "location_id", "foreignField": "id", "as": "location"}},
	}
	pipe := collection.Pipe(pipeline)

	err = pipe.One(&results)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Pipeline failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Pipeline failed")
		return
	}

	params2 := params.(map[string]interface{})

	if int(params2["status"].(float64)) == models.BOOKING_STATUS_CONSENT || int(params2["status"].(float64)) == models.BOOKING_STATUS_REFUSAL {

		var title string

		switch int(params2["status"].(float64)) {
		case models.BOOKING_STATUS_CONSENT:
			title = "Booking is Consent!"
		case models.BOOKING_STATUS_REFUSAL:
			title = "Booking is Refusal!"
		}

		err = u.Resources.Mail.Send( results.Email,
			title,
			templates.Processed(results))
		if err != nil {
			u.Resources.Log.Panic().Msgf("Error send mail to %s: %s", results.Email, err.Error())
		}
	}

	WriteSuccessResponse(response, results, "")
}