package params

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"stori-service/src/libs/adapters/external_call/fetch"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetParamInt(t *testing.T) {
	// fixture
	paramName := "test_param"
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Getting param", func(t *testing.T) {
			// inject spy
			GetParameter = func(paramName string) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":12}`)),
				}, nil
			}
			// execute
			param, err := GetParamInt(paramName)
			// assert
			assert.NoError(t, err)
			assert.Equal(t, 12, param)

			t.Cleanup(func() {
				GetParameter = fetch.GetParameterResponse
			})
		})
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("Error in body", func(t *testing.T) {
			// inject spy
			GetParameter = func(paramName string) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"errors":"["error":"F in the chat"]"}`)),
				}, nil
			}
			// execute
			param, err := GetParamInt(paramName)
			// assert
			assert.Zero(t, param)
			assert.Error(t, err)

			t.Cleanup(func() {
				GetParameter = fetch.GetParameterResponse
			})
		})
		t.Run("HTTP error status", func(t *testing.T) {
			// inject spy
			GetParameter = func(paramName string) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":12}`)),
				}, nil
			}
			// execute
			param, err := GetParamInt(paramName)
			// assert
			assert.Zero(t, param)
			assert.Error(t, err)

			t.Cleanup(func() {
				GetParameter = fetch.GetParameterResponse
			})
		})
	})
}
