package handler

import (
	"net/http"
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/incu6us/asterisk-ami-api/internal/utils/config"
	"github.com/gorilla/mux"
	"github.com/incu6us/asterisk-ami-api/internal/platform/ami"
	"strconv"
	"github.com/bit4bit/gami"
)

type apiHandler struct {
	ContentType string
	amiClient ami.AMI
}

type response struct {
	Result interface{} `json:"Result"`
}

const (
	CONTENT_TYPE = "application/json"
)

var (
	amiResponse *gami.AMIResponse
	err error
	log = logging.MustGetLogger("main")
	conf = config.GetConfig()
)

func (a *apiHandler) amiInit() {
	var host = conf.Ami.Host + ":" + strconv.Itoa(conf.Ami.Port)

	a.amiClient = ami.GetAMI(host, conf.Ami.Username, conf.Ami.Password)
	if err = a.amiClient.Run(); err != nil {
		log.Error("Error:", err)
	}else {
		log.Info("AMI connection established")
	}

}

func (a *apiHandler) setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", a.ContentType)
	w.WriteHeader(200)
}

func (a apiHandler) print(w http.ResponseWriter, r *http.Request, message interface{}) {
	defer r.Body.Close()
	a.setJsonHeader(w)

	if encodeError := json.NewEncoder(w).Encode(response{message}); encodeError != nil {
		log.Warning("Parse message error", encodeError)
	}
}

func (a *apiHandler) Test(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	a.print(w, r, vars)
}

func (a *apiHandler) CallFromSipToMSISDN(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	sipId := vars["SIPID"]
	msisdn := vars["MSISDN"]

	var params = make(map[string]string)
	params["Channel"] = "SIP/"+sipId
	params["CallerID"] = "manual_"+msisdn
	params["MaxRetries"] = "0"
	params["RetryTime"] = "1"
	params["WaitTime"] = "20"
	params["Context"] = conf.Asterisk.Context
	params["Exten"] = msisdn
	params["Priority"] = "1"

	log.Debug("Originate: %v", params)

	if amiResponse, err = a.amiClient.Originate(params); err != nil || amiResponse.Status == "Error" {
		log.Error("AMI Action error! Error: %v, AMI Response Status: %s", err, amiResponse.Status)
	}

	a.print(w, r, amiResponse)
}

type ApiHandler interface {
	Test(w http.ResponseWriter, r *http.Request)
	CallFromSipToMSISDN(http.ResponseWriter, *http.Request)
}

func GetHandler() ApiHandler {
	var a *apiHandler

	if a == nil {
		a = &apiHandler{ContentType: CONTENT_TYPE}
		a.amiInit()
	}

	return a
}