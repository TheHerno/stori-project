package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	customMocks "stori-service/src/utils/test/mock"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	validClientToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IlNhY2hhQGdtYWlsLmNvbSIsIm5hbWUiOiJTYWNoYSBpZ3VhZ3VhdSIsInJvbGUiOiJjbGllbnQiLCJ1c2VyX2lkIjo0fQ.WK2MjJy9gc9XNCfaGBwW1nqlxBgb9YzgpxygPKF7FHU"
	validAdminToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1lc3NpQGdtYWlsLmNvbSIsIm5hbWUiOiJMaW9uZWwgTWVzc2kiLCJyb2xlIjoiYWRtaW4iLCJ1c2VyX2lkIjoxfQ.BT1wjRkMj6Zverq-T8WWBuXbEyRBKNi86UEfYx0GrCM"
	// user name with 203 characters
	invalidUserToken      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWEiLCJyb2xlIjoiYWRtaW4iLCJpZCI6MSwiZW1haWwiOiJtZXNzaUBnbWFpbC5jb20ifQ.sp_3m0M8A4O9SKqk0T6XEV6NrFvDyR3H6xCr8dZBSQw"
	invalidSignatureToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiTGlvbmVsIE1lc3NpIiwicm9sZSI6ImFkbWluIiwiaWQiOjEsImVtYWlsIjoibWVzc2lAZ21haWwuY29tIn0.7tkQwd6YC_-nTTVwtVcuAfdmQorb3ub1dfzgyWAuveRjuHUZZ_sNSZ5tXEYhjP6hFbtg4jTwxh-xYqCI9LWC5A"
)

func TestGetEnv(t *testing.T) {
	t.Run("Should panic", func(t *testing.T) {
		secret := os.Getenv("AUTH_APP_SESSION_SECRET")
		os.Setenv("AUTH_APP_SESSION_SECRET", "")
		assert.Panics(t, getEnv)
		os.Setenv("AUTH_APP_SESSION_SECRET", secret)
	})
	t.Run("Should not panic", func(t *testing.T) {
		assert.NotPanics(t, getEnv)
	})
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		authMiddleware := NewAuthMiddleware()
		testCases := []struct {
			name    string
			token   string
			handler func(http.Handler) http.Handler
		}{
			{
				name:    "admin",
				token:   validAdminToken,
				handler: authMiddleware.HandlerAdmin(),
			},
			{
				name:    "client",
				token:   validClientToken,
				handler: authMiddleware.HandlerClient(),
			},
		}
		for _, tC := range testCases {
			t.Run(tC.name, func(t *testing.T) {
				mockHTTPHandler := new(customMocks.MockHTTPHandler)
				mockHTTPHandler.On("ServeHTTP", mock.AnythingOfType("*http.response"), mock.AnythingOfType("*http.Request")).Return()
				ts := httptest.NewServer(tC.handler(mockHTTPHandler))
				req, _ := http.NewRequest("GET", ts.URL, nil)
				req.Header.Add("Authorization", "Bearer "+tC.token)
				resp, _ := ts.Client().Do(req)
				//Data Assertion
				assert.NotNil(t, resp)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				// Mock Assertion: Behavioral
				mockHTTPHandler.AssertExpectations(t)
				mockHTTPHandler.AssertNumberOfCalls(t, "ServeHTTP", 1)
			})
		}
	})

	t.Run("Should fail on", func(t *testing.T) {
		authMiddleware := NewAuthMiddleware()
		testCases := []struct {
			name    string
			token   string
			handler func(http.Handler) http.Handler
		}{
			{
				name:    "Client on admin handler",
				token:   validClientToken,
				handler: authMiddleware.HandlerAdmin(),
			},
			{
				name:    "Admin on client handler",
				token:   validAdminToken,
				handler: authMiddleware.HandlerClient(),
			},
			{
				name:    "Invalid user token",
				token:   invalidUserToken,
				handler: authMiddleware.HandlerAdmin(),
			},
			{
				name:    "Empty token",
				token:   "",
				handler: authMiddleware.HandlerClient(),
			},
			{
				name:    "Token modification: +length",
				token:   validClientToken + "UK48KG4I",
				handler: authMiddleware.HandlerClient(),
			},
			{
				name:    "Token modification: -length",
				token:   validClientToken[:100],
				handler: authMiddleware.HandlerClient(),
			},
			{
				name:    "Token modification: content",
				token:   validAdminToken[:len(validAdminToken)-1] + "z",
				handler: authMiddleware.HandlerClient(),
			},
			{
				name:    "Token with another signature",
				token:   invalidSignatureToken,
				handler: authMiddleware.HandlerAdmin(),
			},
		}
		for _, tC := range testCases {
			t.Run(tC.name, func(t *testing.T) {
				mockHTTPHandler := new(customMocks.MockHTTPHandler)
				ts := httptest.NewServer(tC.handler(mockHTTPHandler))

				req, _ := http.NewRequest("GET", ts.URL, nil)
				req.Header.Add("Authorization", "Bearer "+tC.token)
				resp, _ := ts.Client().Do(req)

				//Data Assertion
				if assert.NotNil(t, resp) {
					assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

					// Mock Assertion: Behavioral
					mockHTTPHandler.AssertNumberOfCalls(t, "ServeHTTP", 0)
				}
			})
		}
		t.Run("No authorization token", func(t *testing.T) {
			mockHTTPHandler := new(customMocks.MockHTTPHandler)
			ts := httptest.NewServer(authMiddleware.HandlerAdmin()(mockHTTPHandler))

			req, _ := http.NewRequest("GET", ts.URL, nil)
			req.Header.Add("ImNotAuthorization", "Bearer "+validAdminToken)
			resp, _ := ts.Client().Do(req)

			//Data Assertion
			if assert.NotNil(t, resp) {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

				// Mock Assertion: Behavioral
				mockHTTPHandler.AssertNumberOfCalls(t, "ServeHTTP", 0)
			}
		})
	})
}

func TestUserValidation(t *testing.T) {
	// fixture
	validEmail := "messi@gmail.com"
	validName := "Lionel Messi"
	validRole := "admin"
	validID := 1
	longString := strings.Repeat("a", 201)
	shortString := "a"
	validUser := &User{
		UserID: validID,
		Email:  validEmail,
		Name:   validName,
		Role:   validRole,
	}
	t.Run("Should success on", func(t *testing.T) {
		err := validUser.Validate()
		assert.NoError(t, err)
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			name string
			user *User
		}{
			{
				name: "Empty email",
				user: &User{
					UserID: validID,
					Name:   validName,
					Email:  "",
					Role:   validRole,
				},
			},
			{
				name: "Invalid ID",
				user: &User{
					UserID: 0,
					Name:   validName,
					Email:  validEmail,
					Role:   validRole,
				},
			},
			{
				name: "Empty name",
				user: &User{
					UserID: validID,
					Name:   "",
					Email:  validEmail,
					Role:   validRole,
				},
			},
			{
				name: "Empty role",
				user: &User{
					UserID: validID,
					Name:   validName,
					Email:  validEmail,
					Role:   "",
				},
			},
			{
				name: "Long role",
				user: &User{
					UserID: validID,
					Name:   validName,
					Email:  validEmail,
					Role:   longString,
				},
			},
			{
				name: "Short role",
				user: &User{
					UserID: validID,
					Name:   validName,
					Email:  validEmail,
					Role:   shortString,
				},
			},
			{
				name: "Long Name",
				user: &User{
					UserID: validID,
					Name:   longString,
					Email:  validEmail,
					Role:   validRole,
				},
			},
			{
				name: "Short Name",
				user: &User{
					UserID: validID,
					Name:   shortString,
					Email:  validEmail,
					Role:   validRole,
				},
			},
			{
				name: "Short Email",
				user: &User{
					UserID: validID,
					Name:   validName,
					Email:  shortString,
					Role:   validRole,
				},
			},
			{
				name: "Long email",
				user: &User{
					UserID: validID,
					Name:   validName,
					Email:  longString,
					Role:   validRole,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.name, func(t *testing.T) {
				err := tC.user.Validate()
				assert.Error(t, err)
			})
		}
	})
}
