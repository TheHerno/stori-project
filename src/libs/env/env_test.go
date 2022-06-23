package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	value := 21
	stringValue := fmt.Sprint(value)
	defaultValue := 60
	environmentKey := "ENV_TESTING"

	t.Run("processIntEnvVar", func(t *testing.T) {
		t.Run("Valid int", func(t *testing.T) {
			// Fixture
			var variable int
			os.Setenv(environmentKey, stringValue)

			// Action
			processIntEnvVar(&variable, environmentKey, defaultValue)

			// Assert Data
			assert.Equal(t, variable, value)
		})
		t.Run("Invalid int", func(t *testing.T) {
			// Fixture
			var variable int
			os.Setenv(environmentKey, "No int")

			// Action
			processIntEnvVar(&variable, environmentKey, defaultValue)

			// Assert Data
			assert.Equal(t, variable, defaultValue)
		})
		t.Run("Invalid environment key", func(t *testing.T) {
			// Fixture
			var variable int

			// Action
			processIntEnvVar(&variable, "Invalid key", defaultValue)

			// Assert Data
			assert.Equal(t, variable, defaultValue)
		})
	})
}
