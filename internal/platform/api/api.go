package api

import (
	"github.com/gorilla/mux"
	"github.com/incu6us/asterisk-ami-api/internal/platform/api/handler"
	"net/http"
)

const (
	PATH_PREFIX = "/api/v1/"
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
		if api.Name != "ready" {
			router.
			PathPrefix(PATH_PREFIX).
				Methods(api.Method).
				Path(api.Pattern).
				Name(api.Name).
				Handler(api.HandlerFunc)
		}else{
			router.
			PathPrefix(PATH_PREFIX).
				Methods(api.Method).
				Path(api.Pattern).
				Name(api.Name).
				Handler(api.HandlerFunc)
		}
	}

	return router
}

var apis = APIs{
	API{
		"callSipToMSISDN",
		"GET",
		"/call/{SIPID}/{MSISDN}", // ?async=false # default
		handler.GetHandler().CallFromSipToMSISDN,
	},
	API{
		"sendSMS",
		"POST",
		"/modem/send/sms/{modem}/{MSISDN}",
		handler.GetHandler().SendSms,
	},
	API{
		"ready",
		"GET",
		"/ready",
		handler.GetHandler().Ready,
	},
}
