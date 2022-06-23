package database

import (
	"fmt"
	"stori-service/src/libs/env"
	"stori-service/src/libs/logger"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	ormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db   *gorm.DB
	once sync.Once
)

//CreateTrainingConnectionString returns the connection string based on environment variables
func CreateTrainingConnectionString() string {
	//db config vars
	dbHost := env.StoriServicePostgresqlHost
	dbPort := env.StoriServicePostgresqlPort
	dbName := env.StoriServicePostgresqlName
	dbUser := env.StoriServicePostgresqlUsername
	dbPassword := env.StoriServicePostgresqlPassword
	dbSSLMode := env.StoriServicePostgresqlSSLMode
	if env.AppEnv == "testing" {
		dbName = env.StoriServicePostgresqlNameTest
	}
	//Make connection string with interpolation
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
	return connectionString
}

/*
SetupTrainingGormDB open the pool connection in db var and return it
*/
func SetupTrainingGormDB() *gorm.DB {
	once.Do(func() {
		config := &gorm.Config{
			Logger: ormlogger.Default.LogMode(ormlogger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}
		//connect to db
		var dbError error
		db, dbError = gorm.Open(postgres.Open(CreateTrainingConnectionString()), config)
		for dbError != nil {
			logger.GetInstance().Error("Failed to connect to own-database")
			time.Sleep(env.StoriServiceSecondsBetweenAttempts)
			logger.GetInstance().Info("Retrying...")
			db, dbError = gorm.Open(postgres.Open(CreateTrainingConnectionString()), config)
		}
		logger.GetInstance().Info("Connected to own-database!")
		setConnectionMaxLifetime(db, 0) //To be reused forever
	})
	return db
}

/*
GetTrainingGormConnection return db pointer which already have an open connection
*/
func GetTrainingGormConnection() *gorm.DB {
	return SetupTrainingGormDB()
}
