package app

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"path/filepath"

	"nomad/api/src/api/v1"
	"nomad/api/src/resources"
	"nomad/api/src/app/filters"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
)

type Book struct {
	Title  string
	Author string
}

type App struct {
	Resources   *resources.Resources
}

func (a *App) Run() {

	restful.PrettyPrintResponses = false
	restful.Filter(filters.LoggerFilter)
	//wsContainer.Filter(wsContainer.OPTIONSFilter)
	restful.Add(v1.InitAuthController().WebService())
	restful.Add(v1.InitUserController().WebService())

	// swagger

	config := restfulspec.Config{
		WebServices: restful.DefaultContainer.RegisteredWebServices(), // you control what services are visible
		APIPath:     "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	//

	//http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir(config.ApiDir))))
	//http.HandleFunc("/check_alive", func(rw http.ResponseWriter, req *http.Request) { rw.Header().Set("X-Status", "OK") })

	//

	restful.DefaultContainer.Router(restful.CurlyRouter{})

	ws := new(restful.WebService)
	ws.Route(ws.GET("/apidocs/{subpath:*}").To(staticFromPathParam))
	ws.Route(ws.GET("/apidocs/").To(func(req *restful.Request, resp *restful.Response) {

		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}

		exPath := filepath.Dir(ex)
		exPath += "/../dist/index.html"

		http.ServeFile(
			resp.ResponseWriter,
			req.Request,
			exPath)
	}))
	restful.Add(ws)

	restful.DefaultContainer.Filter(restful.DefaultContainer.OPTIONSFilter)

	cors := restful.CrossOriginResourceSharing{AllowedHeaders: []string{"Authorization"}, CookiesAllowed: true, Container: restful.DefaultContainer}
	restful.Filter(cors.Filter)
	//restful.DefaultContainer.Filter(restful.DefaultContainer.OPTIONSFilter)

	//

	addr := fmt.Sprintf(":%d", a.Resources.Config.Port)
	a.Resources.Log.Info().Msgf("Start listening on %s", addr)
	server := &http.Server{
		Addr: addr,
		Handler: restful.DefaultContainer,
	}
	log.Fatal(server.ListenAndServe())
}

func staticFromPathParam(req *restful.Request, resp *restful.Response) {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)
	exPath += "/../dist/"
	exPath += req.PathParameter("subpath")

	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		exPath)
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Nomad API",
			Description: "Main Nomad API",
			Contact: &spec.ContactInfo{
				Name:  "a.borisov",
				Email: "deeseefromcd@gmail.com",
				URL:   "https://nomad.space",
			},
			License: &spec.License{
				Name: "MIT",
				URL:  "http://mit.org",
			},
			Version: "1.0.0",
		},
	}
	//swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
	//	Name:        "Auth",
	//	Description: "Auth Api"}}}
}