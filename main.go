package main

import (
	"github.com/incu6us/asterisk-ami-api/internal/platform/api"
	"github.com/incu6us/asterisk-ami-api/internal/utils/config"
	"log"
	"net/http"
)

func main() {

	var conf = config.GetConfig()
	var err error

	srv := &http.Server{
		Addr: conf.General.Listen,
		//ReadTimeout:  api.HTTP_TIMEOUT,
		//WriteTimeout: api.HTTP_TIMEOUT,
		Handler: api.NewHandler(),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
