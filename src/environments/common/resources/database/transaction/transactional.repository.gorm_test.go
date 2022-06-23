package database

import (
	"os"
	"stori-service/src/libs/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	database.SetupTrainingGormDB()
	code := m.Run()
	os.Exit(code)
}

func TestTransactionalGORMRepository(t *testing.T) {
	t.Run("Begin", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			t.Run("With connection", func(t *testing.T) {
				// Fixture
				connection := database.GetTrainingGormConnection()
				basetx := connection.Begin()
				baseRepository := TransactionalGORMRepository{basetx, 0}

				// Action
				tx := baseRepository.Begin(basetx).(*gorm.DB)

				// Assert data
				assert.NotNil(t, tx)
				assert.Equal(t, basetx, tx)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Without connection", func(t *testing.T) {
				// Fixture
				connection := database.GetTrainingGormConnection()
				baseRepository := TransactionalGORMRepository{connection, 0}

				// Action
				tx := baseRepository.Begin(nil).(*gorm.DB)

				// Assert data
				assert.NotNil(t, tx)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Create table", func(t *testing.T) {
				// Fixture
				connection := database.GetTrainingGormConnection()
				baseRepository := TransactionalGORMRepository{connection, 0}

				// Action
				tx := baseRepository.Begin(nil).(*gorm.DB)
				tx = tx.Exec(`DROP TABLE IF EXISTS "test";
				CREATE TABLE "test"(
					"id" serial
				);`)
				tx = tx.Exec(`SELECT * FROM "test"`)
				connection = connection.Exec(`SELECT * FROM "test"`)

				// Assert data
				assert.NoError(t, tx.Error)
				assert.Error(t, connection.Error, "Should fail because transaccion is not committed")

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("Commit", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Create table", func(t *testing.T) {
				// Fixture
				connection := database.GetTrainingGormConnection()
				baseRepository := TransactionalGORMRepository{connection, 0}
				tx := baseRepository.Begin(nil).(*gorm.DB)
				tx.Exec(`DROP TABLE IF EXISTS "test";
					CREATE TABLE "test"(
						"id" serial
					);`)

				// Actions
				baseRepository.Commit()
				query := connection.Exec(`SELECT * FROM "test"`)

				// Assert data
				assert.NoError(t, query.Error)

				t.Cleanup(func() {
					connection.Exec(`DROP TABLE IF EXISTS "test"`)
				})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Commit two times", func(t *testing.T) {
				// Fixture
				connection := database.GetTrainingGormConnection()
				baseRepository := TransactionalGORMRepository{connection, 0}
				baseRepository.Begin(nil)
				err1 := baseRepository.Commit()

				// Action
				err2 := baseRepository.Commit()

				// Assert data
				assert.NoError(t, err1)
				assert.Error(t, err2, "Should fail because transaction was committed previously")
			})
		})
	})
	t.Run("Rollback", func(t *testing.T) {
		t.Run("Should pass", func(t *testing.T) {
			t.Run("Create table", func(t *testing.T) {
				// Fixture
				connection := database.GetTrainingGormConnection()
				baseRepository := TransactionalGORMRepository{connection, 0}

				// Action
				tx := baseRepository.Begin(nil).(*gorm.DB)
				tx = tx.Exec(`DROP TABLE IF EXISTS "test";
					CREATE TABLE "test"(
						"id" serial
					);`)
				err1 := tx.Exec(`SELECT * FROM "test"`).Error
				err2 := baseRepository.Rollback()
				err3 := tx.Exec(`SELECT * FROM "test"`).Error

				// Assert data
				assert.NoError(t, err1)
				assert.NoError(t, err2)
				assert.Error(t, err3, "Should fail because transaction was rolled back")

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
		t.Run("Should fail", func(t *testing.T) {
			t.Run("Rollback two times", func(t *testing.T) {
				// Fixture
				connection := database.GetTrainingGormConnection()
				baseRepository := TransactionalGORMRepository{connection, 0}
				baseRepository.Begin(nil)
				err1 := baseRepository.Rollback()

				// Action
				err2 := baseRepository.Rollback()

				// Assert data
				assert.NoError(t, err1)
				assert.Error(t, err2, "Should fail because transaction was rolled back previously")
			})
			t.Run("Unsupported driver on RollbackTo", func(t *testing.T) {
				// Fixture
				fakeConnection := &gorm.DB{
					Config: &gorm.Config{
						ConnPool: &gorm.PreparedStmtDB{},
					},
				}
				baseRepository := TransactionalGORMRepository{fakeConnection, 1} //It has a savepoint

				// Action
				err := baseRepository.Rollback() //It's a rollbackTO

				// Assert data
				assert.Error(t, err, "Should fail because transaction was commited previously")
			})
		})
	})
	t.Run("SavePoint", func(t *testing.T) {
		t.Run("Should pass", func(t *testing.T) {
			t.Run("With rollback", func(t *testing.T) {
				// Fixture
				var countBeforeSavePoint, countAfterSavePoint, countAfterRollback int64
				connection := database.GetTrainingGormConnection()
				baseRepository := TransactionalGORMRepository{connection, 0}

				// Action
				tx := baseRepository.Begin(nil).(*gorm.DB)
				tx.Exec(`CREATE TABLE "test"(
						"id" serial
				);`)
				err1 := tx.Table("test").Count(&countBeforeSavePoint).Error
				err2 := baseRepository.SavePoint()
				tx.Table("test").Create(map[string]interface{}{"id": 1})
				err3 := tx.Table("test").Count(&countAfterSavePoint).Error
				err4 := baseRepository.Rollback()
				err5 := tx.Table("test").Count(&countAfterRollback).Error

				// Assert data
				assert.NoError(t, err1)
				assert.NoError(t, err2)
				assert.NoError(t, err3)
				assert.NoError(t, err4)
				assert.NoError(t, err5)
				assert.Equal(t, int64(0), countBeforeSavePoint)
				assert.Equal(t, int64(1), countAfterSavePoint)
				assert.Equal(t, int64(0), countAfterRollback)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
		t.Run("Should fail", func(t *testing.T) {
			t.Run("Unsupported driver", func(t *testing.T) {
				// Fixture
				fakeConnection := &gorm.DB{
					Config: &gorm.Config{
						ConnPool: &gorm.PreparedStmtDB{},
					},
				}
				baseRepository := TransactionalGORMRepository{fakeConnection, 0}

				// Action
				err2 := baseRepository.SavePoint()

				// Assert data
				assert.Error(t, err2, "Should fail because it's a fake connection")
			})
		})
	})
}
