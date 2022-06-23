package helpers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
Middleware (this function) makes adding more than one layer of middleware easy
by specifying them as a list. It will run the first specified middleware first.
*/
func Middleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	length := len(middlewares)
	for i := range middlewares {
		handler = middlewares[length-1-i](handler)
	}
	return handler
}

/*
IDFromRequestToInt returns the ID from the request as an int.
*/
func IDFromRequestToInt(request *http.Request) (int, error) {
	warehouseId, err := strconv.Atoi(mux.Vars(request)["id"])
	return warehouseId, err
}

//PointerToString is a helper to create (inline) pointers to string value, returns nil if string is empty
func PointerToString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

/*
StringInSlice receives a string and a slice
then tries to find the string in the list
*/
func StringInSlice(toFind string, list []string) bool {
	for _, validItems := range list {
		if toFind == validItems {
			return true
		}
	}
	return false
}
