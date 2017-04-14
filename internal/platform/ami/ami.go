package ami

import (
	"github.com/bit4bit/gami"
	"github.com/bit4bit/gami/event"
	"log"
	"math/rand"
	"strings"
	"time"
)

type ami struct {
	randGenDigit int
	host         string
	user         string
	pass         string
}

var (
	amiClient *gami.AMIClient
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
		if amiClient, err = gami.Dial(a.host); err != nil {
			log.Panic(err)
			return err
		}

		amiClient.Run()
		defer amiClient.Close()
		//log.Info("AMI Params: %v", amiClient)

		//install manager
		go func() {
			for {
				select {
				//handle network errors
				case err := <-amiClient.NetError:
					log.Println("Network Error:", err)
					//try new connection every second
					<-time.After(time.Second * 1)
					if err := amiClient.Reconnect(); err == nil {
						//call start actions
						amiClient.Action("Events", gami.Params{"EventMask": "on"})
					}

				case err := <-amiClient.Error:
					log.Panic("error:", err)
				//wait events and process
				case ev := <-amiClient.Events:
					log.Println("Event Detect:", *ev)
					//if want type of events
					log.Println("EventType:", event.New(ev))
				}
			}
		}()

		if a.user != "" {
			if err = amiClient.Login(a.user, a.pass); err != nil {
				log.Panicf("AMI login failed: %v", err)
				return err
			}
		}
	}

	return err
}

func (a *ami) CustomAction(action string, params map[string]string) (*gami.AMIResponse, error) {
	var amiParams = make(map[string]string)

	amiParams["ActionID"] = strings.ToLower(action) + "-" + a.randGenSuffix()

	for k, v := range params {
		if !strings.EqualFold(k, "actionid") {
			amiParams[k] = v
		}
	}

	actionChanResponse, err := amiClient.Action(action, amiParams)
	if err != nil {
		return nil, err
	}

	return actionChanResponse, err
}

//func (a *ami) Originate(params map[string]string, async bool) (interface{}, error) {
func (a *ami) Originate(params map[string]string) (*gami.AMIResponse, error) {

	actionResponse, err := a.CustomAction("Originate", params)
	if err != nil {
		log.Printf("AMI Action error! Error: %v, AMI Response Status: %s", err, actionResponse)
		return nil, err
	}

	return actionResponse, nil
}

type AMI interface {
	Run() error
	CustomAction(action string, params map[string]string) (*gami.AMIResponse, error)
	Originate(params map[string]string) (*gami.AMIResponse, error)
}

func GetAMI(host, user, pass string) AMI {
	return &ami{
		host: host,
		user: user,
		pass: pass,
	}
}
