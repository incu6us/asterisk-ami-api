package handler

import (
	"encoding/json"
	"github.com/bit4bit/gami"
	"github.com/gorilla/mux"
	"github.com/incu6us/asterisk-ami-api/internal/platform/ami"
	"github.com/incu6us/asterisk-ami-api/internal/utils/config"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
	"strconv"
)

type apiHandler struct {
	ContentType string
	amiClient   ami.AMI
}

type response struct {
	Result interface{} `json:"Result"`
}

const (
	CONTENT_TYPE = "application/json"
)

var (
	//amiResponse *gami.AMIResponse
	hendler *apiHandler
	l       = logging.MustGetLogger("handler")
	conf    = config.GetConfig()
)

func (a *apiHandler) amiInit() {
	var err error
	var host = conf.Ami.Host + ":" + strconv.Itoa(conf.Ami.Port)

	a.amiClient = ami.GetAMI(host, conf.Ami.Username, conf.Ami.Password)
	if err = a.amiClient.Run(); err != nil {
		l.Error("Error:", err)
	} else {
		l.Info("AMI connection established")
	}

}

func (a *apiHandler) setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", a.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (a apiHandler) print(w http.ResponseWriter, r *http.Request, message interface{}) {
	a.setJsonHeader(w)

	if encodeError := json.NewEncoder(w).Encode(response{message}); encodeError != nil {
		l.Warning("Parse message error", encodeError)
	}
}

func (a *apiHandler) Test(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a.print(w, r, vars)
}

func (a *apiHandler) CallFromSipToMSISDN(w http.ResponseWriter, r *http.Request) {

	var err error

	vars := mux.Vars(r)

	sipId := vars["SIPID"]
	msisdn := vars["MSISDN"]
	async, _ := strconv.ParseBool(r.URL.Query().Get("async"))

	var amiResponse interface{}

	l.Debug("vars", vars, async)

	var params = make(map[string]string)
	params["Channel"] = "SIP/" + sipId
	params["CallerID"] = "manual_" + msisdn
	params["MaxRetries"] = "0"
	params["RetryTime"] = "1"
	params["WaitTime"] = "20"
	params["Context"] = conf.Asterisk.Context
	params["Exten"] = msisdn
	params["Priority"] = "1"

	if async {
		params["Async"] = "true"
	}

	l.Debug("Originate: %v", params)

	if amiResponse, err = a.amiClient.Originate(params); err != nil {
		l.Error("AMI Action error! Error: %v, AMI Response Status: %s", err)
		a.print(w, r, err)
		return
	}

	a.print(w, r, amiResponse)

}

func (a *apiHandler) PlaybackAdvertisement(w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)

	audioFile := vars["FILE"]
	msisdn := vars["MSISDN"]
	async, _ := strconv.ParseBool(r.URL.Query().Get("async"))

	var amiResponse interface{}

	l.Debug("vars", vars, async)

	var params = make(map[string]string)
	params["Channel"] = "local/" + msisdn + "@" + conf.Asterisk.Context
	params["CallerID"] = "playback_" + msisdn
	params["MaxRetries"] = "5"
	params["RetryTime"] = "10"
	params["WaitTime"] = "20"
	params["Context"] = conf.Asterisk.PlaybackContext
	params["Priority"] = "1"
	params["Variable"] = "AudioFile=" + audioFile

	if async {
		params["Async"] = "true"
	}

	l.Debug("Originate: %v", params)

	if amiResponse, err = a.amiClient.Originate(params); err != nil {
		l.Error("AMI Action error! Error: %v, AMI Response Status: %s", err)
		a.print(w, r, err)
		return
	}

	a.print(w, r, amiResponse)

}

func (a *apiHandler) SendSms(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()

	var err error
	var body []byte
	var amiResponse <-chan *gami.AMIResponse

	vars := mux.Vars(r)

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		a.print(w, r, err)
		return
	}

	var params = make(map[string]string)
	params["Device"] = vars["modem"]
	params["Number"] = vars["MSISDN"]
	params["Message"] = string(body)

	l.Debug("Send SMS: %v", params)

	if amiResponse, err = a.amiClient.CustomAction("DongleSendSMS", params); err != nil {
		l.Error("AMI Action error! Error: %v, AMI Response Status: %s", err)
		a.print(w, r, err)
		return
	}

	resp := <-amiResponse
	a.print(w, r, resp)
}

// simple check which improve, that server is running
func (a *apiHandler) Ready(w http.ResponseWriter, r *http.Request) {
	a.print(w, r, "Service is up and running")
}

type ApiHandler interface {
	Test(w http.ResponseWriter, r *http.Request)
	CallFromSipToMSISDN(http.ResponseWriter, *http.Request)
	PlaybackAdvertisement(http.ResponseWriter, *http.Request)
	SendSms(w http.ResponseWriter, r *http.Request)
	Ready(w http.ResponseWriter, r *http.Request)
}

func GetHandler() ApiHandler {

	if hendler == nil {
		hendler = &apiHandler{ContentType: CONTENT_TYPE}
		hendler.amiInit()
	}

	return hendler
}
