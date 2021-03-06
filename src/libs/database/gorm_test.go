package database

import (
	"stori-service/src/libs/env"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func resetOnceStori() {
	once = sync.Once{}
}

func TestSetupStoriGormDB(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			resetOnceStori()
			db := SetupStoriGormDB()
			sqlDB, _ := db.DB()
			errPing := sqlDB.Ping()
			//Data Assertion
			assert.NotNil(t, db)
			assert.NoError(t, errPing)
			t.Cleanup(func() {
			})
		})
		t.Run("Wait for postgres", func(t *testing.T) {
			// Smaller time & wrong DB name
			oldDelta := env.StoriServiceSecondsBetweenAttempts
			env.StoriServiceSecondsBetweenAttempts = time.Second / 2
			oldValue := env.StoriServicePostgresqlNameTest
			env.StoriServicePostgresqlNameTest = "Stori_SERVICE_POSTGRESQL_NAME_not_found"
			var db *gorm.DB
			var errPing error
			wait := make(chan bool)
			go func() {
				resetOnceStori()
				db = SetupStoriGormDB()
				sqlDB, _ := db.DB()
				errPing = sqlDB.Ping()
				wait <- true
			}()
			time.Sleep(env.StoriServiceSecondsBetweenAttempts)
			env.StoriServicePostgresqlNameTest = oldValue
			<-wait

			//Data Assertion
			assert.NotNil(t, db)
			assert.NoError(t, errPing)
			t.Cleanup(func() {
				env.StoriServicePostgresqlNameTest = oldValue
				env.StoriServiceSecondsBetweenAttempts = oldDelta
			})
		})
	})
}

func TestGetStoriGormConnection(t *testing.T) {
	t.Run("Should success when the connection is already open", func(t *testing.T) {
		resetOnceStori()
		db := SetupStoriGormDB()
		dbSingleton := GetStoriGormConnection()
		sqlDB, _ := dbSingleton.DB()
		errPing := sqlDB.Ping()
		//Data Assertion
		assert.Equal(t, db, dbSingleton)
		assert.NoError(t, errPing)
	})
}
