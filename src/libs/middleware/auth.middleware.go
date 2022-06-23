package middleware

import (
	"context"
	"net/http"
	"stori-service/src/libs/validator"
	"strings"

	"github.com/golang-jwt/jwt"
)

/*
User model for Middleware user
*/
type User struct {
	jwt.StandardClaims
	UserID int    `json:"user_id" validate:"required,gt=0"`
	Name   string `json:"name" validate:"required,min=2,max=200"`
	Email  string `json:"email" validate:"required,email,min=2,max=200"`
	Role   string `json:"role" validate:"required,min=2,max=50"`
}

/*
Validate validates the User struct
*/
func (user *User) Validate() error {
	if err := validator.ValidateStruct(user); err != nil {
		return err
	}
	return nil
}

type contextKey int

const (
	ContextKeyUser contextKey = iota
)

var (
	authKey string
)

// getEnv gets the environment variable
func getEnv() {
	authKey = "123"
	if authKey == "" {
		panic("no AUTH_APP_SESSION_SECRET")
	}
}

// panic if the key is not set
func init() {
	getEnv()
}

//IAuthMiddleware interface for auth middleware, useful for mocking too
type IAuthMiddleware interface {
	HandlerAdmin() func(next http.Handler) http.Handler
	HandlerClient() func(next http.Handler) http.Handler
}

type authMiddleware struct{}

//NewAuth0Middleware is a constructor for middleware struct
func NewAuthMiddleware() IAuthMiddleware {
	return new(authMiddleware)
}

/*
Handler can be applied to a route to require authentication.
If the jwt is valid then the request will continue to the endpoint
otherwise it will return a 401 status
*/
func (auth *authMiddleware) handler(role string) func(next http.Handler) http.Handler {
	// grab te jwt from the request
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// grab the jwt from the request
			// if it is empty return a 401
			token := extractJWT(r)
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := parseJWT(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			err = user.Validate()
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if user.Role != role {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextKeyUser, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// HandlerAdmin is a middleware that checks if the user is an admin
func (auth *authMiddleware) HandlerAdmin() func(next http.Handler) http.Handler {
	return auth.handler("admin")
}

// HandlerCLient is a middleware that checks if the user is a client
func (auth *authMiddleware) HandlerClient() func(next http.Handler) http.Handler {
	return auth.handler("client")
}

// extractJWT gets the jwt from the request header
func extractJWT(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	// split the header into 2 parts
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// parseJWT parses the token and returns the claims
func parseJWT(token string) (*User, error) {
	// parse the jwt
	parsedToken, err := jwt.ParseWithClaims(token, &User{}, func(token *jwt.Token) (interface{}, error) {
		// validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(authKey), nil
	})
	if err != nil {
		return nil, err
	}
	// validate the claims
	if claims, ok := parsedToken.Claims.(*User); ok && parsedToken.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
