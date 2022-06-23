package sentry

import (
	"net/http"
	"testing"

	"stori-service/src/libs/env"
	customMock "stori-service/src/utils/test/mock"

	"github.com/gorilla/handlers"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureDSN = "http://0@localhost/0"
)

func TestSetupSentry(t *testing.T) {
	prodDSN = fixtureDSN
	t.Run("Should success on", func(t *testing.T) {
		t.Run("With empty dsn", func(t *testing.T) {
			prodDSN = ""

			assert.NotPanics(t, SetupSentry)

			t.Cleanup(func() {
				prodDSN = fixtureDSN
			})
		})
		t.Run("With dsn on production", func(t *testing.T) {
			env.AppEnv = "production"

			assert.NotPanics(t, SetupSentry)

			t.Cleanup(func() {
				env.AppEnv = "testing"
			})
		})
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("Invalid dsn", func(t *testing.T) {
			env.AppEnv = "production"
			prodDSN = "invalid_dsn"

			assert.Panics(t, SetupSentry)

			t.Cleanup(func() {
				env.AppEnv = "testing"
				prodDSN = fixtureDSN
			})
		})
	})
}

func TestHandler(t *testing.T) {
	prodDSN = fixtureDSN
	t.Run("Should succeed on", func(t *testing.T) {
		t.Run("dev env", func(t *testing.T) {
			//Fixture
			SetupSentry()
			handlerFn := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				panic("mock panic")
			})

			sentryHandler := Handler().Handle(handlerFn)
			handler := handlers.RecoveryHandler()(sentryHandler)

			// Prepare request
			response := customMock.MHTTPHandle("GET", "/", handler.ServeHTTP, "", nil, nil)

			//Data Assertion
			assert.NotNil(t, response)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		})
		t.Run("prod env", func(t *testing.T) {
			env.AppEnv = "production"

			//Fixture
			SetupSentry()
			handlerFn := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				panic("mock panic")
			})

			sentryHandler := Handler().Handle(handlerFn)
			handler := handlers.RecoveryHandler()(sentryHandler)

			// Prepare request
			response := customMock.MHTTPHandle("GET", "/", handler.ServeHTTP, "", nil, nil)

			//Data Assertion
			assert.NotNil(t, response)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

			t.Cleanup(func() {
				env.AppEnv = "testing"
			})
		})
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("missing dsn on prod", func(t *testing.T) {
			env.AppEnv = "production"
			prodDSN = ""

			assert.Panics(t, SetupSentry)

			t.Cleanup(func() {
				env.AppEnv = "testing"
				prodDSN = fixtureDSN
			})
		})
	})
}
