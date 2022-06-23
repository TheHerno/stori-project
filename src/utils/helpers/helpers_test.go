package helpers

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	callStack := make([]string, 0)
	//Fixture
	handler := func(response http.ResponseWriter, request *http.Request) {
		callStack = append(callStack, "handler")
		response.WriteHeader(http.StatusTeapot)
	}
	firstMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callStack = append(callStack, "firstMiddleware")
			h.ServeHTTP(w, r)
		})
	}
	secondMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callStack = append(callStack, "secondMiddleware")
			h.ServeHTTP(w, r)
		})
	}

	//This line do the actual work
	result := Middleware(
		http.HandlerFunc(handler),
		firstMiddleware,
		secondMiddleware,
	)
	ts := httptest.NewServer(result)
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL, nil)
	res, _ := ts.Client().Do(req)

	//Data Assertion
	assert.Equal(t, http.StatusTeapot, res.StatusCode)
	assert.Equal(t, "firstMiddleware", callStack[0])
	assert.Equal(t, "secondMiddleware", callStack[1])
	assert.Equal(t, "handler", callStack[2])
}

func TestIDFromRequestToInt(t *testing.T) {
	//Fixture
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL, nil)
	t.Run("Should success on", func(t *testing.T) {
		// Fixture
		vars := map[string]string{"id": "1234"}
		reqWithVars := mux.SetURLVars(req, vars)

		// action
		got, err := IDFromRequestToInt(reqWithVars)

		// assertion
		assert.Equal(t, 1234, got)
		assert.NoError(t, err)
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("Not numeric ID", func(t *testing.T) {
			// Fixture
			vars := map[string]string{"id": "1234a"}
			reqWithVars := mux.SetURLVars(req, vars)

			// action
			got, err := IDFromRequestToInt(reqWithVars)

			// assertion
			assert.Equal(t, 0, got)
			assert.Error(t, err)
		})
	})
}

func TestPointerToString(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		result := PointerToString("")
		//Data Assertion
		assert.Nil(t, result)
		assert.IsType(t, new(string), result)
	})
	t.Run("Non-empty string", func(t *testing.T) {
		result := PointerToString("test")
		//Data Assertion
		assert.NotNil(t, result)
		assert.Equal(t, "test", *result)
		assert.IsType(t, new(string), result)
	})
}

func TestStringInSlice(t *testing.T) {
	//Fixture
	list := []string{"Foo", "Bar", "Foo Bar", "Mock", "Testing", "With spaces"}
	type arg struct {
		ToFind string
		List   []string
	}
	t.Run("Should success on", func(t *testing.T) {

		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomIndex := r1.Intn(len(list)) //[0, len)

		testCases := []struct {
			TestName string
			Input    arg
		}{
			{
				TestName: "First word",
				Input: arg{
					ToFind: list[0],
					List:   list,
				},
			},
			{
				TestName: "Last word",
				Input: arg{
					ToFind: list[len(list)-1],
					List:   list,
				},
			},
			{
				TestName: "Random word",
				Input: arg{
					ToFind: list[randomIndex],
					List:   list,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.TestName, func(t *testing.T) {
				found := StringInSlice(tC.Input.ToFind, tC.Input.List)
				assert.True(t, found)
			})
		}
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Input    arg
		}{
			{
				TestName: "Empty list",
				Input: arg{
					ToFind: list[0],
					List:   []string{},
				},
			},
			{
				TestName: "Empty string to find",
				Input: arg{
					ToFind: "",
					List:   list,
				},
			},
			{
				TestName: "Not found",
				Input: arg{
					ToFind: "Test not found",
					List:   list,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.TestName, func(t *testing.T) {
				found := StringInSlice(tC.Input.ToFind, tC.Input.List)
				assert.False(t, found)
			})
		}
	})
}
