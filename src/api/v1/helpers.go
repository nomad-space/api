package v1

import (
	"github.com/emicklei/go-restful"
)

type ResponseSuccess struct {
	Data		interface{} 		`json:"data"`
	Msg			string				`json:"msg"`
	Status		string 				`json:"status"`
}
type ResponseError struct {
	Error		string 				`json:"error"`
	Msg			string				`json:"msg"`
	Status		string 				`json:"status"`
}

func WriteSuccessResponse(response *restful.Response, data interface{}, msg string){
	response.WriteAsJson(ResponseSuccess{data, msg, "success"})
}
func WriteErrorResponse(response *restful.Response, httpStatus int, msg string){
	response.WriteHeaderAndEntity(httpStatus, ResponseError{msg, msg, "error"})
}