package sentry

import (
	"net/http"
	"stori-service/src/libs/env"
	"stori-service/src/libs/logger"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

var prodDSN string = "" //Sentry client initialized with an empty DSN. Using noopTransport. No events will be delivered.

/*
SetupSentry init the sentry connection with some params
*/
func SetupSentry() {
	dsn := ""
	if env.AppEnv == "production" {
		dsn = prodDSN
		if dsn == "" {
			panic("missing dsn")
		}
	}
	options := sentry.ClientOptions{
		Dsn: dsn,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request
					logger.GetInstance().Warn("Request panic ", req)
				}
			}
			logger.GetInstance().Warn("Sentry event", event)
			return event
		},
		Debug:            true,
		AttachStacktrace: true,
	}
	if err := sentry.Init(options); err != nil {
		panic(err)
	}
}

/*
Handler returns a sentry handler that works as a middleware and sends only panic errors
*/
func Handler() *sentryhttp.Handler {
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
	return sentryHandler
}
