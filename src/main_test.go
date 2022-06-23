package src

import (
	"net/http"
	"net/http/httptest"
	myErrors "stori-service/src/libs/errors"
	"stori-service/src/utils"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSetupHandler(t *testing.T) {
	t.Run("Should not panics", func(t *testing.T) {
		assert.NotPanics(t, func() { SetupHandler() })
	})
}

func TestPingEndpoint(t *testing.T) {
	t.Run("Should response http 200 status", func(t *testing.T) {
		muxRouter := mux.NewRouter()
		pingEndpoint(muxRouter)
		ts := httptest.NewServer(muxRouter)
		defer ts.Close()
		req, _ := http.NewRequest("GET", ts.URL+"/ping", nil)
		res, _ := ts.Client().Do(req)

		//Data Assertion
		type data struct {
			Service  string    `json:"service"`
			Datetime time.Time `json:"datetime"`
		}

		bodyResult, _ := utils.GetBodyResponse(res, &data{})
		bodyData := bodyResult.Data.(*data)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "stori-service is online", bodyData.Service)
		assert.Empty(t, bodyResult.Errors)
		assert.Equal(t, http.StatusText(http.StatusOK), bodyResult.Message)
		assert.WithinDuration(t, time.Now(), bodyData.Datetime, 1*time.Second)
	})
}

func TestCustomNotFoundHanlder(t *testing.T) {
	t.Run("Should response a custom body on url not found", func(t *testing.T) {
		router := new(mux.Router)
		customNotFoundHanlder(router)
		ts := httptest.NewServer(router)
		defer ts.Close()

		//Action
		resp, err := http.DefaultClient.Get(ts.URL + "/not_found_url")

		//Data Assertion
		bodyResponse, _ := utils.GetBodyResponse(resp, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, myErrors.ErrURLNotFound.Error(), bodyResponse.Message)
		assert.Equal(t, myErrors.ErrURLNotFound.Error(), bodyResponse.Errors[0]["error"])
	})
}
