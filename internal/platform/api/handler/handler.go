package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/incu6us/asterisk-ami-api/internal/platform/ami"
	"github.com/incu6us/asterisk-ami-api/internal/utils/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type apiHandler struct {
	ContentType string
	//amiClient   ami.AMI
}

type response struct {
	Result interface{} `json:"Result"`
}

const (
	CONTENT_TYPE = "application/json"
)

var (
	handler   *apiHandler
	conf      = config.GetConfig()
	amiClient ami.AMI
)

func (a *apiHandler) amiInit() {
	var err error
	var host = conf.Ami.Host + ":" + strconv.Itoa(conf.Ami.Port)

	amiClient = ami.GetAMI(host, conf.Ami.Username, conf.Ami.Password)
	if err = amiClient.Run(); err != nil {
		log.Println("Error:", err)
	} else {
		log.Println("AMI connection established")
	}

}

func (a *apiHandler) setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", a.ContentType)
	w.WriteHeader(http.StatusOK)
}

func (a apiHandler) print(w http.ResponseWriter, r *http.Request, message interface{}) {
	a.setJsonHeader(w)

	if encodeError := json.NewEncoder(w).Encode(response{message}); encodeError != nil {
		log.Println("Parse message error", encodeError)
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

	//var amiResponse interface{}

	log.Println("vars", vars, async)

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

	log.Println("Originate: %v", params)

	amiResponse, err := amiClient.Originate(params)
	if err != nil {
		log.Panicf("AMI Action error! Error: %v, AMI Response Status: %s", err, amiResponse)
		a.print(w, r, err.Error())
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

	log.Println("vars", vars, async)

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

	log.Printf("Originate: %v", params)

	log.Println(amiClient)
	amiResponse, err := amiClient.Originate(params)
	if err != nil {
		log.Panicf("AMI Action error! Error: %v, AMI Response Status: %s", err, amiResponse)
		a.print(w, r, err.Error())
		return
	}

	log.Println("---")
	a.print(w, r, amiResponse)
}

func (a *apiHandler) SendSms(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()

	var err error
	var body []byte
	//var amiResponse *gami.AMIResponse

	vars := mux.Vars(r)

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		a.print(w, r, err)
		return
	}

	var params = make(map[string]string)
	params["Device"] = vars["modem"]
	params["Number"] = vars["MSISDN"]
	params["Message"] = string(body)

	log.Printf("Send SMS: %v", params)

	amiResponse, err := amiClient.CustomAction("DongleSendSMS", params)
	if err != nil {
		log.Panicf("AMI Action error! Error: %v, AMI Response Status: %s", err, amiResponse)
		a.print(w, r, err.Error())
		return
	}

	a.print(w, r, amiResponse)
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

	if handler == nil {
		handler = &apiHandler{ContentType: CONTENT_TYPE}
		handler.amiInit()
	}

	return handler
}
