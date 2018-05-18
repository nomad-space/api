package filters

import (
	"time"
	"net/http"
	"strings"
	"errors"
	"context"

	"nomad/api/src/models"
	"nomad/api/src/resources"
	restful "github.com/emicklei/go-restful"
)

func LoggerFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	t := time.Now()
	d := time.Now().Sub(t)

	res, err := resources.GetInstance()
	if err != nil {
		res.Log.Fatal().Msg(err.Error())
	}

	res.Log.Debug().
		Str("method", req.Request.Method).
		Str("uri", req.Request.URL.RequestURI()).
		Int("status_code", resp.StatusCode()).
		Int("content_length", resp.ContentLength()).
		Float64("time2", float64(d/time.Nanosecond)).
		Float64("time", float64(d/time.Millisecond)).
		Msgf("%s %s", req.Request.Method, req.Request.URL.RequestURI())

	chain.ProcessFilter(req, resp)
}

func ValidateJWT(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	accessToken := req.HeaderParameter("Authorization")

	if accessToken == "" {
		resp.WriteError(http.StatusUnauthorized, errors.New("auth header is empty"))
		return
	}

	parts := strings.SplitN(accessToken, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		resp.WriteError(http.StatusUnauthorized, errors.New("auth header is invalid"))
		return
	}

	tokenString := parts[1]

	claims, err := models.ParseJWTToken(tokenString)
	if err != nil {
		resp.WriteError(http.StatusUnauthorized, err)
		return
	}

	ctx := context.WithValue(req.Request.Context(), "JwtClaims", claims)
	ctx = context.WithValue(ctx, "JwtToken", tokenString)
	req.Request = req.Request.WithContext(ctx)

	chain.ProcessFilter(req, resp)
}