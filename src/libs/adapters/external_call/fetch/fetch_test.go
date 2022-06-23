package fetch

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"stori-service/src/libs/env"
	myErrors "stori-service/src/libs/errors"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestDoRequest(t *testing.T) {
	t.Run("Multiple 4XX responses, should not open", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusNotFound) //404
		}))
		defer ts.Close() //Server will be closed when test ends
		req, _ := http.NewRequest("GET", ts.URL, nil)

		var st gobreaker.Settings
		st.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= consecutiveFailuresToOpenCB
		}
		circuitBreakerToTest := gobreaker.NewCircuitBreaker(st)

		paramsEndpointCB := paramsEndpoint.circuitBreaker

		//First 400 fail, it should not open
		response, err := doRequest(circuitBreakerToTest, req)
		assert.Nil(t, response)
		assert.Error(t, err)
		assert.Equal(t, "closed", circuitBreakerToTest.State().String())
		assert.Equal(t, "closed", paramsEndpointCB.State().String()) //Others CBs should remain closed

		//Second 400 fail, it should not open
		response, err = doRequest(circuitBreakerToTest, req)
		assert.Nil(t, response)
		assert.Error(t, err)
		assert.Equal(t, "closed", circuitBreakerToTest.State().String())
		assert.Equal(t, "closed", paramsEndpointCB.State().String()) //Others CBs should remain closed

		//Thirth 400 fail, it should not open
		response, err = doRequest(circuitBreakerToTest, req)
		assert.Nil(t, response)
		assert.Error(t, err)
		assert.Equal(t, "closed", circuitBreakerToTest.State().String())
		assert.Equal(t, "closed", paramsEndpointCB.State().String()) //Others CBs should remain closed
	})
	t.Run("no such host", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
		}))
		ts.Close() //Close the server for simulate no such host
		req, _ := http.NewRequest("GET", ts.URL, nil)

		var st gobreaker.Settings
		st.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= consecutiveFailuresToOpenCB
		}
		circuitBreakerToTest := gobreaker.NewCircuitBreaker(st)

		paramsEndpointCB := paramsEndpoint.circuitBreaker
		response, err := doRequest(circuitBreakerToTest, req)
		assert.Nil(t, response)
		assert.EqualError(t, err, myErrors.ErrConnectionProvider.Error())
		time.Sleep(time.Second / 10) //Wait one second, because the CB is instantly opened asynchronously
		assert.Equal(t, "open", circuitBreakerToTest.State().String())
		assert.Equal(t, "closed", paramsEndpointCB.State().String()) //Others CBs should remain closed
	})
	t.Run("Multiple 5XX responses", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close() //Server will be closed when test ends
		req, _ := http.NewRequest("GET", ts.URL, nil)

		circuitBreakerToTest := paramsEndpoint.circuitBreaker
		paramsEndpointCB := paramsEndpoint.circuitBreaker

		//First 500 fail, it should not open
		response, err := doRequest(circuitBreakerToTest, req)
		assert.Nil(t, response)
		assert.Error(t, err)
		assert.Equal(t, "closed", circuitBreakerToTest.State().String())
		assert.Equal(t, "closed", paramsEndpointCB.State().String()) //Others CBs should remain closed

		//Second 500 fail, it should not open
		response, err = doRequest(circuitBreakerToTest, req)
		assert.Nil(t, response)
		assert.Error(t, err)
		assert.Equal(t, "closed", circuitBreakerToTest.State().String())
		assert.Equal(t, "closed", paramsEndpointCB.State().String()) //Others CBs should remain closed

		//Thirth 500 fail, it should open
		response, err = doRequest(circuitBreakerToTest, req)
		assert.Nil(t, response)
		assert.Error(t, err)
		assert.Equal(t, "open", circuitBreakerToTest.State().String())
		assert.Equal(t, "closed", paramsEndpointCB.State().String()) //Others CBs should remain closed

		//Fourth 500 fail, it came from circuit breaker
		response, err = doRequest(circuitBreakerToTest, req)
		assert.Nil(t, response)
		assert.EqualError(t, err, myErrors.ErrConnectionProvider.Error())
		assert.Equal(t, "open", circuitBreakerToTest.State().String())
		assert.Equal(t, "closed", paramsEndpointCB.State().String()) //Others CBs should remain closed
	})
}

func TestBuildRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Endpoint with default options", func(t *testing.T) {
			endpoint := &endpoint{
				method: http.MethodGet,
				path:   "/test",
			}

			req, err := buildRequest(endpoint, Options{})

			//Data Assertion
			assert.Equal(t, env.ParamsURL+"/test", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Content-type"))
			assert.Equal(t, http.MethodGet, req.Method)
			assert.Empty(t, req.Body)
			assert.NoError(t, err)
		})
		t.Run("endpoint with body", func(t *testing.T) {
			type BodyTest struct {
				Test string
			}

			bodyToReq := &BodyTest{
				Test: "testing value",
			}
			endpoint := &endpoint{
				method: http.MethodPost,
				path:   "/test",
			}
			req, err := buildRequest(endpoint, Options{Body: bodyToReq})

			//Data Assertion
			assert.NoError(t, err)
			assert.Equal(t, env.ParamsURL+"/test", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Content-type"))
			assert.Equal(t, http.MethodPost, req.Method)

			var bodyFromReq *BodyTest
			err2 := json.NewDecoder(req.Body).Decode(&bodyFromReq)
			assert.NoError(t, err2)
			assert.EqualValues(t, bodyFromReq, bodyToReq)
		})
		t.Run("endpoint with headers", func(t *testing.T) {
			headers := map[string]string{
				"header-test": "test value",
			}

			endpoint := &endpoint{
				method: http.MethodPost,
				path:   "/test",
			}
			req, err := buildRequest(endpoint, Options{Headers: headers})

			//Data Assertion
			assert.NoError(t, err)
			assert.Equal(t, env.ParamsURL+"/test", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Content-type"))
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Empty(t, req.Body)
			assert.Equal(t, "test value", req.Header.Get("header-test"))
		})
		t.Run("endpoint with path variables", func(t *testing.T) {
			var pathVariables []interface{}
			pathVariables = append(pathVariables, "value1", "value2")

			endpoint := &endpoint{
				method: http.MethodPost,
				path:   "/test/%s/foo/%s",
			}
			req, err := buildRequest(endpoint, Options{PathVariables: pathVariables})

			//Data Assertion
			assert.NoError(t, err)
			assert.Equal(t, env.ParamsURL+"/test/value1/foo/value2", req.URL.String(), "should have value1 and value2")
			assert.Equal(t, "application/json", req.Header.Get("Content-type"))
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Empty(t, req.Body)
		})
	})
	t.Run("Fail", func(t *testing.T) {
		t.Run("Wrong request body", func(t *testing.T) {
			endpoint := &endpoint{
				method: http.MethodPost,
				path:   "/test",
			}
			value := make(chan int) //chan can't be serialize, it'll cause and error with json.Marshal
			req, err := buildRequest(endpoint, Options{Body: value})

			//Data Assertion
			assert.Nil(t, req)
			assert.Error(t, err)
		})
		t.Run("Unknown method", func(t *testing.T) {
			endpoint := &endpoint{
				method: "unknown method",
				path:   "/test",
			}
			req, err := buildRequest(endpoint, Options{})

			//Data Assertion
			assert.Nil(t, req)
			assert.Error(t, err)
		})
	})
}

func TestOpenManuallyCB(t *testing.T) {
	t.Run("Should success on executing the function and waiting ", func(t *testing.T) {
		var st gobreaker.Settings
		//Open the circuit breaker when it has at least 1 consecutive failures
		st.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 1
		}
		cb := gobreaker.NewCircuitBreaker(st)
		assert.Equal(t, "closed", cb.State().String())
		openManuallyCB(cb)
		time.Sleep(time.Second / 10) //Wait .1 second, because the CB is instantly opened asynchronously
		assert.Equal(t, "open", cb.State().String())
	})
}

func TestGetParameteresponse(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Find the Max Stock Movement Param", func(t *testing.T) {
			//Mocking the endpoint response
			defer gock.Off()
			gock.New(env.ParamsURL).
				Get("/v1/console/parameters/value/attr/max_stock_movement").
				Reply(200).
				JSON([]byte(`[]`))

			response, err := GetParameterResponse("max_stock_movement")
			//Data Assertion
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, "/params/v1/console/parameters/value/attr/max_stock_movement", response.Request.URL.Path)
			assert.Equal(t, paramsEndpoint.method, response.Request.Method)
		})
	})
}
