package ami

import (
	"errors"
	"github.com/ivahaev/amigo"
	"github.com/ivahaev/amigo/uuid"
	"log"
	"strings"
)

type amiAmigo struct {
	randGenDigit int
	host         string
	user         string
	pass         string
}

var (
	amigoClient *amigo.Amigo
)

func (a *amiAmigo) Run() error {
	var err error

	host := strings.Split(a.host, ":")[0]
	port := strings.Split(a.host, ":")[1]

	settings := &amigo.Settings{Username: a.user, Password: a.pass, Host: host, Port: port}
	amigoClient = amigo.New(settings)

	amigoClient.Connect()

	// Listen for connection events
	amigoClient.On("connect", func(message string) {
		log.Println("Connected", message)
	})
	amigoClient.On("error", func(message string) {
		log.Println("Connection error:", message)
		err = errors.New(message)
	})

	return err
}

func (a *amiAmigo) CustomAction(action string, params map[string]string) (map[string]string, error) {
	params["Action"] = action
	resp, err := amigoClient.Action(params)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}

func (a *amiAmigo) Originate(params map[string]string) (map[string]string, error) {
	params["ActionID"] = uuid.NewV4()

	if _, ok := params["Variable"]; !ok {
		params["Variable"] = "ActionID=" + params["ActionID"]
	}
	params["Action"] = "Originate"

	log.Printf("Originate: %v", params)

	resp, err := amigoClient.Action(params)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}

type AMIA interface {
	Run() error
	CustomAction(action string, params map[string]string) (map[string]string, error)
	Originate(params map[string]string) (map[string]string, error)
}

func GetAMIAmigo(host, user, pass string) AMIA {
	return &amiAmigo{
		host: host,
		user: user,
		pass: pass,
	}
}
