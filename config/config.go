package config

import (
	"io"
	"stori-service/src/libs/database"
	"stori-service/src/libs/logger"
	"stori-service/src/libs/sentry"
)

/*
slice of dependencies, io.Closes is an interface with method Close() error
all package that makes connections implements it
*/
var dependenciesToClose []io.Closer

/*
SetupCommonDependencies calls setup for each necessary dependencies
and registers them on one slice to be closed later
*/
func SetupCommonDependencies() {
	logger.SetupLogger()
	sentry.SetupSentry()
	database.SetupStoriGormDB()
	dependenciesToClose = []io.Closer{}
}

/*
TearDownCommonDependencies iterates each dependency and calls Close method
*/
func TearDownCommonDependencies() {
	for _, dependecy := range dependenciesToClose {
		dependecy.Close()
	}
}
