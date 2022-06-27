package env

import (
	"os"
	"strconv"
	"time"
)

var (

	// AppEnv Application Environment
	AppEnv string

	// EnvironmentName Environment Name
	EnvironmentName string

	// ServiceName Service Name
	ServiceName string

	// ServiceVersion Service version
	ServiceVersion string

	// StoriServiceSecondsBetweenAttempts StoriService Interval in Seconds between attempts
	StoriServiceSecondsBetweenAttempts time.Duration

	// StoriServicePostgresqlHost StoriService PostgreSQL host
	StoriServicePostgresqlHost string

	// StoriServicePostgresqlPort StoriService PostgreSQL port
	StoriServicePostgresqlPort string

	// StoriServicePostgresqlName StoriService PostgreSQL name
	StoriServicePostgresqlName string

	// StoriServicePostgresqlNameTest StoriService PostgreSQL name Test
	StoriServicePostgresqlNameTest string

	// StoriServicePostgresqlUsername StoriService PostgreSQL app username
	StoriServicePostgresqlUsername string

	// StoriServicePostgresqlPassword StoriService PostgreSQL app password
	StoriServicePostgresqlPassword string

	// StoriServicePostgresqlSSLMode StoriService PostgreSQL ssl mode
	StoriServicePostgresqlSSLMode string

	// StoriServiceGrpcPort StoriService gRPC port
	StoriServiceGrpcPort string

	// StoriServiceRestPort StoriService Rest port
	StoriServiceRestPort string

	// WhiteList White List
	WhiteList string

	// External services

	// EventLoggerURL Logger service URL
	EventLoggerURL string

	// EventLoggerUser Logger service user
	EventLoggerUser string

	// EventLoggerPassword Logger service password
	EventLoggerPassword string

	// FileRoute Params service URL
	FileRoute string
)

func init() {
	// App Environment
	AppEnv = os.Getenv("APP_ENV")

	// Environment Name
	EnvironmentName = os.Getenv("ENVIRONMENT_NAME")
	// Service Name
	ServiceName = os.Getenv("SERVICE_NAME")
	// Service Version
	ServiceVersion = os.Getenv("VERSION")

	// StoriService - gRPC
	StoriServiceGrpcPort = os.Getenv("STORI_SERVICE_GRPC_PORT")

	// StoriService - Rest
	StoriServiceRestPort = os.Getenv("STORI_SERVICE_REST_PORT")

	// StoriService Interval in Seconds Between Attempts
	var seconds int
	processIntEnvVar(&seconds, "STORI_SERVICE_SECONDS_BETWEEN_ATTEMPTS", 60)
	StoriServiceSecondsBetweenAttempts = time.Duration(seconds) * time.Second

	// StoriService - PostgreSQL
	StoriServicePostgresqlHost = os.Getenv("STORI_SERVICE_POSTGRESQL_HOST")
	StoriServicePostgresqlPort = os.Getenv("STORI_SERVICE_POSTGRESQL_PORT")
	StoriServicePostgresqlName = os.Getenv("STORI_SERVICE_POSTGRESQL_NAME")
	StoriServicePostgresqlNameTest = os.Getenv("STORI_SERVICE_POSTGRESQL_NAME_TEST")
	StoriServicePostgresqlUsername = os.Getenv("STORI_SERVICE_POSTGRESQL_USERNAME")
	StoriServicePostgresqlPassword = os.Getenv("STORI_SERVICE_POSTGRESQL_PASSWORD")
	StoriServicePostgresqlSSLMode = os.Getenv("STORI_SERVICE_POSTGRESQL_SSLMODE")

	// Logger service
	EventLoggerURL = os.Getenv("EVENT_LOGGER_URL")
	EventLoggerUser = os.Getenv("EVENT_LOGGER_USER")
	EventLoggerPassword = os.Getenv("EVENT_LOGGER_PASSWORD")

	// White list
	WhiteList = os.Getenv("WHITE_LIST")

	// Params service
	FileRoute = os.Getenv("FILE_ROUTE")
}

// processIntEnvVar gets environment variable from os and parses it to int
func processIntEnvVar(intVar *int, envKey string, defaultValue int) {
	var err error
	*intVar, err = strconv.Atoi(os.Getenv(envKey))
	if err != nil {
		*intVar = defaultValue
	}
}
