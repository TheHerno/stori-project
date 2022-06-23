package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

/*
AuthServiceMiddleware is a IAuthServiceMiddleware mock
*/
type AuthServiceMiddleware struct {
	mock.Mock
}

//Handler mock method
func (mock *AuthServiceMiddleware) Handler() func(next http.Handler) http.Handler {
	mock.Called()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			next.ServeHTTP(writer, request)
		})
	}
}
