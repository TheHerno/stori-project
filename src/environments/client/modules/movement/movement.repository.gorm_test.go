package movement

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/utils/constant"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// setup
	database.SetupStoriGormDB()
	code := m.Run()
	os.Exit(code)
}

type result struct {
	Movement *entity.Movement
	Err      error
	Date     time.Time
}

var customers = []entity.Customer{
	{
		CustomerID: 1,
		Name:       "User 1",
		Email:      "test1@hotmail.com",
	},
	{
		CustomerID: 2,
		Name:       "User 2",
		Email:      "test2@hotmail.com",
	},
	{
		CustomerID: 3,
		Name:       "User 3",
		Email:      "test3@hotmail.com",
	},
	{
		CustomerID: 4,
		Name:       "User 4",
		Email:      "test4@hotmail.com",
	},
}

var movements = []entity.Movement{
	{
		MovementID: 1,
		CustomerID: 1,
		Quantity:   10,
		Available:  10,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 2,
		CustomerID: 1,
		Quantity:   5,
		Available:  5,
		Type:       constant.OutcomeType,
	},
	{
		MovementID: 3,
		CustomerID: 1,
		Quantity:   10,
		Available:  15,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 4,
		CustomerID: 4,
		Quantity:   17,
		Available:  17,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 5,
		CustomerID: 1,
		Quantity:   20,
		Available:  20,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 6,
		CustomerID: 1,
		Quantity:   100,
		Available:  100,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 7,
		CustomerID: 1,
		Quantity:   75,
		Available:  25,
		Type:       constant.OutcomeType,
	},
	{
		MovementID: 8,
		CustomerID: 1,
		Quantity:   12,
		Available:  12,
		Type:       constant.IncomeType,
	},
}

/*
	Foreign Fixures
*/

func addForeignFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.Customer{}) // cleaning users
	tx.Create(customers)
}

/*
addFixtures adds fixtures to the database
*/
func addFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning users
	tx.Create(movements)
}

func TestGormRepository(t *testing.T) {
	t.Run("BulkCreate", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Creating a batch of movements", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				addForeignFixtures(tx)
				rMovement := NewMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning the table

				err := rMovement.BulkCreate(movements)

				// data assertion
				assert.NoError(t, err)

				// database assertion
				var movementCount int64
				tx.Model(&entity.Movement{}).Count(&movementCount)
				assert.Equal(t, int64(len(movements)), movementCount)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Table doesn't exists", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				rMovement := NewMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning the table
				tx.Migrator().DropTable(&entity.Movement{})

				err := rMovement.BulkCreate(movements)

				// data assertion
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("GetLastMovementByCustomerID", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Getting last movement", func(t *testing.T) {
				// fixture
				customerIDToFind := customers[0].CustomerID
				channelResult1 := make(chan result)
				channelResult2 := make(chan result)

				connection := database.GetStoriGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rMovement := NewMovementGormRepo(tx)

				go func() {
					got, err := rMovement.GetLastMovementByCustomerID(customerIDToFind)
					channelResult1 <- result{got, err, time.Now()}
					time.Sleep(500 * time.Millisecond)
					tx.Commit()
				}()

				go func() {
					err2 := connection.Table("movement").Where("customer_id = ?", customerIDToFind).
						Update("customer_id", "5").Error
					channelResult2 <- result{nil, err2, time.Now()}
				}()

				// assertions
				result1 := <-channelResult1
				assert.NoError(t, result1.Err)
				assert.True(t, cmp.Equal(result1.Movement, &movements[0], cmpopts.IgnoreTypes(time.Time{})))
				// tries to update but still nil because the lock is not released
				result2 := <-channelResult2
				assert.Nil(t, result2.Movement)
				assert.NoError(t, result2.Err)
				assert.True(t, result2.Date.After(result1.Date))

				t.Cleanup(func() {
					tx.Rollback()
					connection.Unscoped().Where("1=1").Delete(&entity.Customer{})
				})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Customer doesn't exists", func(t *testing.T) {
				// fixture
				customerIDToFind := 78
				connection := database.GetStoriGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rMovement := NewMovementGormRepo(tx)

				// action
				got, err := rMovement.GetLastMovementByCustomerID(customerIDToFind)

				// assertions
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
					connection.Unscoped().Where("1=1").Delete(&entity.Customer{})
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				rMovement := NewMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) //Cleaning customers
				tx.Migrator().DropTable(&entity.Movement{})

				got, err := rMovement.GetLastMovementByCustomerID(1)

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
		db := database.GetStoriGormConnection()
		rMovement := NewMovementGormRepo(db)

		clone := rMovement.Clone()

		assert.NotNil(t, clone)
		assert.Equal(t, rMovement, clone)
	})
}
