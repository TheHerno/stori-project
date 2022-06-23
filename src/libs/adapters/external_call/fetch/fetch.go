package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"stori-service/src/libs/env"
	myErrors "stori-service/src/libs/errors"
	"stori-service/src/libs/logger"
	"time"

	"github.com/sony/gobreaker"
)

type endpoint struct {
	method         string
	path           string
	circuitBreaker *gobreaker.CircuitBreaker
}

//Options is a struct for passing dynamic values to the request
type Options struct {
	Body          interface{}
	PathVariables []interface{}
	Headers       map[string]string
}

const (
	consecutiveFailuresToOpenCB uint32 = 3
)

var (
	defaultHTTPClient http.Client = http.Client{}
	//Client to fetch bind API, with default client without certs
	Client         *http.Client = &defaultHTTPClient
	paramsEndpoint              = endpoint{
		method: http.MethodGet,
		path:   "/v1/console/parameters/value/attr/%s",
	}
)

/*
init is like a constructor for this package, only is called (implicit) once
*/
func init() {
	Client = &http.Client{}
	setupCircuitBreaker()
}

/*
GetParameterResponse exported function then calls buildRequest with the correct pathVariables and finally calls doRequest
*/
func GetParameterResponse(parameter string) (*http.Response, error) {
	options := &Options{}
	options.PathVariables = []interface{}{parameter}
	req, err := buildRequest(&paramsEndpoint, *options)
	if err != nil {
		return nil, err
	}
	return doRequest(paramsEndpoint.circuitBreaker, req)
}

/*
Creates all circuit breaker with settings (this is only called once)
*/
func setupCircuitBreaker() {
	var st gobreaker.Settings
	st.Name = "Fetch params"

	//Open the circuit breaker when it has at least 3 consecutive failures
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		return counts.ConsecutiveFailures >= consecutiveFailuresToOpenCB
	}

	//is the period of the open state, after which the state of CircuitBreaker becomes half-open
	st.Timeout = 30 * time.Second

	st.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		logger.GetInstance().Warningf("%s has changed from %s to %s", name, from, to)
	}

	cb := gobreaker.NewCircuitBreaker(st)
	paramsEndpoint.circuitBreaker = cb
}

/*
buildRequest receives endpoint and options, first it checks if there are path variables
that need to be interpolated in path, then it checks if options has body so it marshales and assigns it
creates the request, sets default and dynamic headers.
*/
func buildRequest(endpoint *endpoint, options Options) (*http.Request, error) {

	endpointPath := endpoint.path
	if len(options.PathVariables) > 0 {
		endpointPath = fmt.Sprintf(endpointPath, options.PathVariables...) //Interpolating path variables into endpoint path
	}
	fullPath := fmt.Sprint(env.ParamsURL, endpointPath)
	var body io.Reader = nil
	if options.Body != nil {
		bodyJSON, err := json.Marshal(&options.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(bodyJSON)
	}

	req, err := http.NewRequest(endpoint.method, fullPath, body)
	if err != nil {
		return nil, err
	}
	//Default header
	req.Header.Set("Content-type", "application/json")

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

/*
doRequest receives a circuit breaker (CB) and a request, then makes the request with circuit breaker
and checks for http status that opens or not the CB
*/
func doRequest(cb *gobreaker.CircuitBreaker, req *http.Request) (*http.Response, error) {
	logger.GetInstance().Infof("Request to: %s, URL: %s", req.Method, req.URL.String())
	var generalError error = nil
	response, err := cb.Execute(func() (interface{}, error) {
		resp, err := Client.Do(req)
		if err != nil {
			logger.GetInstance().Errorf("Error response from bind: %s", err)
			//Errors like timeout, no such host will open the circuit instantly
			go openManuallyCB(cb)
			generalError = myErrors.ErrConnectionProvider
			return nil, generalError
		}
		logger.GetInstance().Infof("Response status: %s", resp.Status)

		if resp.StatusCode >= 500 { //These errors will count as failure for circuit breaker
			generalError = errors.New(resp.Status)
			return nil, generalError
		}

		if resp.StatusCode >= 400 { //These errors will NOT count as failure for circuit breaker
			generalError = fmt.Errorf("API status error: %s", resp.Status)
			return nil, nil
		}

		return resp, nil
	})

	if generalError != nil {
		return nil, generalError
	}

	if errors.Is(err, gobreaker.ErrOpenState) {
		//It could be a fallback here, because circuit breaker is open.
		return nil, myErrors.ErrConnectionProvider
	}

	return response.(*http.Response), nil
}

//openManuallyCB Hack the cb by opening it manually simulating fail requests
func openManuallyCB(cb *gobreaker.CircuitBreaker) {
	for i := 0; i < int(consecutiveFailuresToOpenCB); i++ {
		go cb.Execute(func() (interface{}, error) {
			return nil, errors.New("dummy error")
		})
	}
}
