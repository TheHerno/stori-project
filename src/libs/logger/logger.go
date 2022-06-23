package logger

import (
	"fmt"
	"os"
	"stori-service/src/libs/env"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-extras/elogrus.v7"
)

var (
	//with this we can inject a spy on unit test
	newAsyncElasticHookWithFunc = elogrus.NewAsyncElasticHookWithFunc

	//with this we can inject a spy on unit test
	newClient = elasticsearch.NewClient

	logger *logrus.Entry = logrus.NewEntry(logrus.StandardLogger()) //Default console logger
)

/*
GetInstance returns the singleton instance
*/
func GetInstance() *logrus.Entry {
	return logger
}

/*
SetupLogger setup function to instantiate ELK logger
*/
func SetupLogger() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	log.SetLevel(logrus.TraceLevel)

	client, err := prepareElasticSearchClient()
	if err != nil {
		log.Errorln("Elastic client error:", err)
		return
	}
	hostname, _ := os.Hostname()
	hook, err := newAsyncElasticHookWithFunc(client, hostname, logrus.TraceLevel, indexName)
	if err != nil {
		log.Errorln("Elastic logger error:", err)
		return
	}
	log.AddHook(hook)

	//Saving logger with default fields
	logger = log.WithFields(logrus.Fields{
		"environment":     env.EnvironmentName,
		"service":         env.ServiceName,
		"service_version": env.ServiceVersion,
	})
	logger.Debug("Logger succesfully connected")
}

/*
prepareElasticSearchClient makes the connection with elasticsearch
*/
func prepareElasticSearchClient() (*elasticsearch.Client, error) {
	return newClient(elasticsearch.Config{
		Addresses: []string{env.EventLoggerURL},
		Username:  env.EventLoggerUser,
		Password:  env.EventLoggerPassword,
	})
}

/*
indexName returns the index with dynamic name (e.g: logs-2020-10-15, logs-2020-10-16)
*/
func indexName() string {
	return fmt.Sprintf("%s-%s", "logs", time.Now().Format("2006.01.02"))
}
