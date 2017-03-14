package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/incu6us/asterisk-ami-api/internal/platform/api/handler"
)

type API struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type APIs []API

func NewHandler() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, api := range apis {
		router.
		Methods(api.Method).
			Path(api.Pattern).
			Name(api.Name).
			Handler(api.HandlerFunc)
	}

	return router
}

var apis = APIs{
	API{
		"callSipToMSISDN",
		"get",
		"/call/{SIPID}/{MSISDN}",
		handler.GetHandler().CallFromSipToMSISDN,
	},
}
