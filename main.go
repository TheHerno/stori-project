package main

import (
	"fmt"
	"net/http"
	"time"

	"stori-service/config"
	"stori-service/src"
	"stori-service/src/libs/env"
	"stori-service/src/libs/logger"
)

func main() {
	config.SetupCommonDependencies()
	defer config.TearDownCommonDependencies()
	handler := src.SetupHandler()

	host := fmt.Sprint(":", env.StoriServiceRestPort)
	srv := &http.Server{
		Handler:      *handler,
		Addr:         host,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}
	go srv.ListenAndServe()
	logger.GetInstance().Info("Server Listening on ", host)

	select {} //Infinite waiting
}
