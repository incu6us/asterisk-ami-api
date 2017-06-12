package api

import (
	"github.com/gorilla/mux"
	"github.com/incu6us/asterisk-ami-api/internal/platform/api/handler"
	"net/http"
	"time"
)

const (
	PATH_PREFIX  = "/api/v1/"
	HTTP_TIMEOUT = 120 * time.Second
)

type API struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type APIs []API

func NewHandler() http.Handler {

	router := mux.NewRouter().StrictSlash(true)
	middlewareHandler := http.TimeoutHandler(router, HTTP_TIMEOUT, "Server timedout!")

	for _, api := range apis {
		if api.Name != "ready" {
			router.
				PathPrefix(PATH_PREFIX).
				Methods(api.Method).
				Path(api.Pattern).
				Name(api.Name).
				Handler(api.HandlerFunc)
		} else {
			router.
				PathPrefix(PATH_PREFIX).
				Methods(api.Method).
				Path(api.Pattern).
				Name(api.Name).
				Handler(api.HandlerFunc)
		}
	}

	return middlewareHandler
}

var apis = APIs{
	API{
		"callSipToMSISDN",
		"GET",
		"/call/{SIPID}/{MSISDN}", // ?async=false # default
		handler.GetHandler().CallFromSipToMSISDN,
	},
	API{
		"playbackAdvertisement",
		"GET",
		"/playback/{MSISDN}/{FILE}", // ?async=false # default
		handler.GetHandler().PlaybackAdvertisement,
	},
	API{
		"sendSMS",
		"POST",
		"/modem/send/sms/{modem}/{MSISDN}",
		handler.GetHandler().SendSms,
	},
	API{
		"getStatByMSISDN",
		"GET",
		"/cdr/search/{MSISDN}", //?startdate=2017-04-26&enddate=2017-04-26 || ?actionid=a3a17318-fc48-4aea-835d-619e0cda585e
		handler.GetHandler().GetStatByMSISDN,
	},
	API{
		"ready",
		"GET",
		"/ready",
		handler.GetHandler().Ready,
	},
}
