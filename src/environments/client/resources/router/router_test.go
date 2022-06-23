package router

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSetupClientRoutes(t *testing.T) {
	t.Run("Should not panics", func(t *testing.T) {
		muxRouter := mux.NewRouter()
		assert.NotPanics(t, func() { SetupClientRoutes(muxRouter) })
	})
}
