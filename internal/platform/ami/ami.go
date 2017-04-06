package ami

import (
	"github.com/bit4bit/gami"
	"github.com/op/go-logging"
	"math/rand"
	"strings"
	"time"
)

type ami struct {
	randGenDigit int
	Host         string
	User         string
	Pass         string
}

var (
	amiClient *gami.AMIClient
	log       = logging.MustGetLogger("ami")
)

func (a *ami) randGenSuffix(i ...int) string {
	var b []rune
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	if len(i) > 0 {
		b = make([]rune, i[0])
	} else {
		b = make([]rune, a.randGenDigit)
	}
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (a *ami) Run() error {

	a.randGenDigit = 10
	var err error

	if amiClient == nil {
		if amiClient, err = gami.Dial(a.Host); err != nil {
			log.Fatal(err)
			return err
		}

		amiClient.Run()

		//log.Info("AMI Params: %v", amiClient)

		//install manager
		go func() {
			for {
				select {
				//handle network errors
				case err := <-amiClient.NetError:
					log.Error("Network Error:", err)
					//try new connection every second
					<-time.After(time.Second)
					if err := amiClient.Reconnect(); err == nil {
						//call start actions
						amiClient.Action("Events", gami.Params{"EventMask": "on"})
					}

				case err := <-amiClient.Error:
					log.Error("error:", err)
					//wait events and process
					//case ev := <-amiClient.Events:
					//	log.Error("Event Detect:", *ev)
					//	//if want type of events
					//	log.Error("EventType:", event.New(ev))
				}
			}
		}()

		if a.User != "" {
			if err = amiClient.Login(a.User, a.Pass); err != nil {
				log.Error("AMI login failed: %v", err)
				return err
			}
		}
	}

	return err
}

func (a *ami) CustomAction(action string, params map[string]string) (<-chan *gami.AMIResponse, error) {
	var actionChanResponse <-chan *gami.AMIResponse
	var amiParams = make(map[string]string)
	var err error

	amiParams["ActionID"] = strings.ToLower(action) + "-" + a.randGenSuffix()

	for k, v := range params {
		if !strings.EqualFold(k, "actionid") {
			amiParams[k] = v
		}
	}

	if actionChanResponse, err = amiClient.AsyncAction(action, amiParams); err != nil {
		return nil, err
	}

	return actionChanResponse, err
}

func (a *ami) Originate(params map[string]string, async bool) (interface{}, error) {
	var actionResponse *gami.AMIResponse
	var actionAsyncResponse <-chan *gami.AMIResponse
	var err error

	if actionAsyncResponse, err = a.CustomAction("Originate", params); err != nil {
		log.Error("AMI Action error! Error: %v, AMI Response Status: %s", err)
	}

	if !async {
		actionResponse = <-actionAsyncResponse
		return actionResponse, err
	} else {
		message := make(map[string]string)
		message["Message"] = "Originate successfully queued"
		return gami.AMIResponse{ID: params["ActionID"], Params: message, Status: "Success"}, err
	}
}

type AMI interface {
	Run() error
	CustomAction(action string, params map[string]string) (<-chan *gami.AMIResponse, error)
	Originate(params map[string]string, async bool) (interface{}, error)
}

func GetAMI(host, user, pass string) AMI {
	return &ami{
		Host: host,
		User: user,
		Pass: pass,
	}
}
