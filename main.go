package main

import (
	"github.com/bit4bit/gami"
	"github.com/op/go-logging"

	"asterisk-custom-api/internal/utils/config"

	"asterisk-custom-api/internal/platform/ami"
	"strconv"
)

func main() {

	var log = logging.MustGetLogger("main")
	var conf = config.GetConfig()

	var amiClient ami.AMI
	var rsPing <-chan *gami.AMIResponse
	var err error

	var host = conf.Ami.Host + ":" + strconv.Itoa(conf.Ami.Port)

	amiClient = ami.GetAMI(host, conf.Ami.Username, conf.Ami.Password)
	if err = amiClient.Run(); err != nil {
		log.Error("Error:", err)
	}

	if rsPing, err = amiClient.Action("Ping", nil); err != nil {
		log.Fatal(err)
	} else {
		log.Debug("ping: %v", <-rsPing)
	}


	if rsPing, err = amiClient.Action("Ping", nil); err != nil {
		log.Fatal(err)
	} else {
		log.Debug("ping: %v", <-rsPing)
	}
}
