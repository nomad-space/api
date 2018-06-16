package v1

import (
	"fmt"
	"net/http"
	"nomad/api/src/resources"
	"nomad/api/src/models"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"gopkg.in/mgo.v2/bson"
)

type HotelController struct {
	Resources	*resources.Resources
}

var instanceHotelController *HotelController

func InitHotelController() *HotelController {
	if instanceHotelController == nil {
		res, _ := resources.GetInstance()
		instanceHotelController = &HotelController{Resources: res}
	}
	return instanceHotelController
}

func (u *HotelController) WebService() *restful.WebService {
	ws := new(restful.WebService)
	tags := []string{"Hotel"}
	ws.Path("/v1/hotel").
		ApiVersion("v1")
		//Consumes(restful.MIME_JSON, restful.DefaultResponseMimeType).
		//Produces(restful.MIME_JSON, restful.DefaultResponseMimeType)
		//Param(ws.HeaderParameter("Authorization", "Bearer JWT token"))

	ws.Route(ws.POST("/").To(u.Create).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		//Filter(filters.ValidateJWT).
		Reads(models.Hotel{}).
		Consumes(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Produces(restful.MIME_JSON, restful.DefaultResponseMimeType).
		Doc("Create"))

	ws.Route(ws.GET("/{hotel_id}").To(u.Info).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("hotel_id", "identifier of the hotel").DataType("string")).
		//Filter(filters.ValidateJWT).
		//Consumes(restful.MIME_JSON).
		//Produces(restful.MIME_JSON).
		Doc("Info"))

	return ws
}

func (u *HotelController) Create(request *restful.Request, response *restful.Response) {

	hotel := models.Hotel{}
	err := request.ReadEntity(&hotel)
	if err != nil {
		WriteErrorResponse(response, http.StatusNotAcceptable, err.Error())
		return
	}

	collection, session, err := u.Resources.Mongo.HotelCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return

	}
	defer session.Close()
	info, err := collection.Upsert(bson.M{"id": hotel.Id}, bson.M{"$set": &hotel})
	if err != nil {
		u.Resources.Log.Debug().Msgf("Create hotel error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}
	fmt.Printf("info: %+v\n", info)

	response.WriteHeaderAndEntity(http.StatusCreated, nil)
}

func (u *HotelController) Info(request *restful.Request, response *restful.Response) {

	hotelId := request.PathParameter("hotel_id")
	u.Resources.Log.Debug().Msgf("hotelId: %+v", hotelId)

	collection, session, err := u.Resources.Mongo.HotelCollectionAndSession();
	if err != nil {
		u.Resources.Log.Debug().Msgf("Connection failed: %s", err.Error())
		WriteErrorResponse(response, http.StatusBadRequest, "Connection failed")
		return

	}
	defer session.Close()

	var results models.Hotel

	err = collection.Find(bson.M{"_id": bson.ObjectIdHex(hotelId)}).One(&results)
	if err != nil {
		u.Resources.Log.Debug().Msgf("Find hotel error: %s", err.Error())
		WriteErrorResponse(response, http.StatusInternalServerError, "")
		return
	}
	u.Resources.Log.Debug().Msgf("results: %s", results)

	WriteSuccessResponse(response, results, "")
}