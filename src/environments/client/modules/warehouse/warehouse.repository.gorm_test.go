package warehouse

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/libs/errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// setup
	database.SetupTrainingGormDB()
	code := m.Run()
	os.Exit(code)
}

type result struct {
	Warehouse *entity.Warehouse
	Err       error
	Date      time.Time
}

var warehouses = []entity.Warehouse{
	{
		WarehouseID: 1,
		Name:        "Warehouse 1",
		Address:     "Address 1",
		UserID:      2,
	},
	{
		WarehouseID: 2,
		Name:        "Warehouse 2",
		Address:     "Address 2",
		UserID:      1,
	},
	{
		WarehouseID: 3,
		Name:        "Warehouse 3",
		Address:     "Address 3",
		UserID:      3,
	},
	{
		WarehouseID: 4,
		Name:        "Warehouse 4",
		Address:     "Address 4",
		UserID:      4,
	},
}

/*
	Fixtures: four warehouses
*/
func addFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning warehouses
	tx.Create(warehouses)
}

func TestGormRepository(t *testing.T) {
	t.Run("FindByUserID", func(t *testing.T) {
		//Fixture
		userID := warehouses[0].UserID
		t.Run("Should success on", func(t *testing.T) {
			tx := database.GetTrainingGormConnection().Begin()
			addFixtures(tx)
			rWarehouse := NewWarehouseGormRepo(tx)

			got, err := rWarehouse.FindByUserID(userID)

			//Data Assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, &warehouses[0], cmpopts.IgnoreFields(entity.Warehouse{}, "UpdatedAt", "CreatedAt")))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})

		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Warehouse not found", func(t *testing.T) {
				//Fixture
				userID := 9999999
				tx := database.GetTrainingGormConnection().Begin()
				addFixtures(tx)

				rWarehouse := NewWarehouseGormRepo(tx)

				got, err := rWarehouse.FindByUserID(userID)

				//Data Assertion
				assert.EqualError(t, err, errors.ErrNotFound.Error())
				assert.Nil(t, got)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rWarehouse := NewWarehouseGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) //Cleaning warehouses
				tx.Migrator().DropTable(&entity.Warehouse{})

				got, err := rWarehouse.FindByUserID(userID)

				//Data Assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("FindAndLockByID", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			// fixture
			userIDToFind := warehouses[1].UserID
			channelResult1 := make(chan result)
			channelResult2 := make(chan result)

			connection := database.GetTrainingGormConnection()
			addFixtures(connection)
			tx := connection.Begin()
			rWarehouse := NewWarehouseGormRepo(tx)

			go func() {
				got, err := rWarehouse.FindAndLockByUserID(userIDToFind)
				channelResult1 <- result{got, err, time.Now()}
				time.Sleep(500 * time.Millisecond)
				tx.Commit()
			}()

			go func() {
				err2 := connection.Table("warehouse").Where("user_id = ?", userIDToFind).
					Update("name", "Warehouseitooo").Error
				channelResult2 <- result{nil, err2, time.Now()}
			}()

			// assertions
			result1 := <-channelResult1
			assert.NoError(t, result1.Err)
			assert.True(t, cmp.Equal(result1.Warehouse, &warehouses[1], cmpopts.IgnoreTypes(time.Time{})))
			// tries to update but still nil because the lock is not released
			result2 := <-channelResult2
			assert.Nil(t, result2.Warehouse)
			assert.NoError(t, result2.Err)
			assert.True(t, result2.Date.After(result1.Date))

			t.Cleanup(func() {
				tx.Rollback()
				connection.Unscoped().Where("1=1").Delete(&entity.Warehouse{})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Warehouse doesn't exists", func(t *testing.T) {
				// fixture
				userIDToFind := 78
				connection := database.GetTrainingGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rWarehouse := NewWarehouseGormRepo(tx)

				// action
				got, err := rWarehouse.FindAndLockByUserID(userIDToFind)

				// assertions
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.EqualError(t, err, errors.ErrNotFound.Error())

				t.Cleanup(func() {
					tx.Rollback()
					connection.Unscoped().Where("1=1").Delete(&entity.Warehouse{})
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rWarehouse := NewWarehouseGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) //Cleaning warehouses
				tx.Migrator().DropTable(&entity.Warehouse{})

				got, err := rWarehouse.FindAndLockByUserID(1)

				//Data Assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("Clone", func(t *testing.T) {
		db := database.GetTrainingGormConnection()
		rWarehouse := NewWarehouseGormRepo(db)

		clone := rWarehouse.Clone()

		assert.NotNil(t, clone)
		assert.Equal(t, rWarehouse, clone)
	})
}
