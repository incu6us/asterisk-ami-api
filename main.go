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

	//var amiClient ami.AMI
	//var rsPing <-chan *gami.AMIResponse
	var err error
	//
	//var host = conf.Ami.Host + ":" + strconv.Itoa(conf.Ami.Port)
	//
	//amiClient = ami.GetAMI(host, conf.Ami.Username, conf.Ami.Password)
	//if err = amiClient.Run(); err != nil {
	//	log.Error("Error:", err)
	//}

	//if rsPing, err = amiClient.Action("Ping", nil); err != nil {
	//	log.Fatal(err)
	//} else {
	//	log.Debug("ping: %v", <-rsPing)
	//}
	//
	//
	//if rsPing, err = amiClient.Action("Ping", nil); err != nil {
	//	log.Fatal(err)
	//} else {
	//	log.Debug("ping: %v", <-rsPing)
	//}




	if err = http.ListenAndServe(conf.General.Listen, api.NewHandler()); err != nil {
		log.Fatal(err)
	}
}
