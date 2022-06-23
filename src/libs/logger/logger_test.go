package logger

import (
	"errors"
	"os"
	"testing"
	"time"

	"stori-service/src/libs/env"
	customSpies "stori-service/src/utils/test/spy"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-extras/elogrus.v7"
)

func TestSetupLogger(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("With valid elastic client (mock)", func(t *testing.T) {
			mockElogrus := customSpies.NewElogrus()
			newAsyncElasticHookWithFunc = mockElogrus.MockNewAsyncElasticHookWithFunc(&elogrus.ElasticHook{}, nil)

			//Action
			SetupLogger()

			//Data Assertion
			assert.Equal(t, 1, mockElogrus.Calls)
			assert.NotEqual(t, logger, logrus.NewEntry(logrus.StandardLogger()), "Should not be the default logger")

			t.Cleanup(func() {
				logger = logrus.NewEntry(logrus.StandardLogger())
				newAsyncElasticHookWithFunc = elogrus.NewAsyncElasticHookWithFunc
			})
		})
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("Without url", func(t *testing.T) {
			oldValue := env.EventLoggerURL
			os.Unsetenv("EVENT_LOGGER_URL")

			//Action
			SetupLogger()

			//Data Assertion
			assert.Equal(t, logger, logrus.NewEntry(logrus.StandardLogger()), "Should be the default logger")
			t.Cleanup(func() {
				os.Setenv("EVENT_LOGGER_URL", oldValue)
				logger = logrus.NewEntry(logrus.StandardLogger())
			})
		})
		t.Run("Prepare elastic client fails", func(t *testing.T) {
			mockElasticsearch := customSpies.NewElasticsearch()
			newClient = mockElasticsearch.MockNewClient(nil, errors.New("some error"))

			//Action
			SetupLogger()

			//Data Assertion
			assert.Equal(t, 1, mockElasticsearch.Calls)
			assert.Equal(t, logger, logrus.NewEntry(logrus.StandardLogger()), "Should be the default logger")
			t.Cleanup(func() {
				logger = logrus.NewEntry(logrus.StandardLogger())
				newClient = elasticsearch.NewClient
			})
		})
	})
}

func TestIndexName(t *testing.T) {

	//Action
	index := indexName()

	//Data Assertion
	assert.Contains(t, index, "logs-")
	assert.Contains(t, index, time.Now().Format("2006.01.02"))
}

func TestGetInstance(t *testing.T) {

	//Action
	instance := GetInstance()

	//Data Assertion
	assert.IsType(t, instance, &logrus.Entry{})
}
