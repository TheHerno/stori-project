package src

import (
	"net/http"
	clientRouter "stori-service/src/environments/client/resources/router"
	"stori-service/src/libs/env"
	myErrors "stori-service/src/libs/errors"
	"stori-service/src/libs/middleware"
	"stori-service/src/libs/sentry"
	"stori-service/src/utils"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/*
SetupHandler returns the handler with all routes and middlewares using mux
*/
func SetupHandler() *http.Handler {
	muxRouter := mux.NewRouter()

	credentialsOk := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{})
	originsOk := handlers.AllowedOrigins(strings.Split(env.WhiteList, ","))
	methodsOk := handlers.AllowedMethods([]string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS", "HEAD"})
	exposeHeadersOk := handlers.ExposedHeaders([]string{
		"X-notifications-unreaded",
		"X-pagination-total-count",
		"X-pagination-page-count",
		"X-pagination-current-page",
		"X-pagination-page-size",
		"needs-action",
	})

	muxRouter.Use(handlers.CORS(originsOk, headersOk, methodsOk, exposeHeadersOk, credentialsOk))
	muxRouter.Use(middleware.LanguageMiddleware)
	settingRoutes(muxRouter)
	customNotFoundHanlder(muxRouter)
	sentryHandler := sentry.Handler()
	handler := sentryHandler.Handle(muxRouter)
	handler = handlers.RecoveryHandler()(handler)

	return &handler
}

/*
settingRoutes takes a pointer to Router and call all environment routers passing its prefix
*/
func settingRoutes(muxRouter *mux.Router) {
	clientRouter.SetupClientRoutes(muxRouter.PathPrefix("/v1/client").Subrouter())
	pingEndpoint(muxRouter)
}

//pingEndpoint is a public endpoint to check the status of this running instance
func pingEndpoint(muxRouter *mux.Router) {
	muxRouter.HandleFunc("/ping", func(response http.ResponseWriter, request *http.Request) {
		bodyResponse := struct {
			Service  string    `json:"service"`
			Datetime time.Time `json:"datetime"`
		}{
			Service:  "stori-service is online",
			Datetime: time.Now(),
		}
		utils.MakeSuccessResponse(response, bodyResponse, http.StatusOK, http.StatusText(http.StatusOK))
	})
}

//customNotFoundHanlder overrides the default not found handler on router
func customNotFoundHanlder(muxRouter *mux.Router) {
	muxRouter.NotFoundHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		utils.MakeErrorResponse(response, myErrors.ErrURLNotFound)
	})
}
