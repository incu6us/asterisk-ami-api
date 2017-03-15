package main

import (
	"github.com/op/go-logging"

	"github.com/incu6us/asterisk-ami-api/internal/utils/config"

	"net/http"
	"github.com/incu6us/asterisk-ami-api/internal/platform/api"
)

func main() {

	var log = logging.MustGetLogger("main")
	var conf = config.GetConfig()

	var err error

	if err = http.ListenAndServe(conf.General.Listen, api.NewHandler()); err != nil {
		log.Fatal(err)
	}
}
