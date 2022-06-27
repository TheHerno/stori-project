package customer

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
	database.SetupStoriGormDB()
	code := m.Run()
	os.Exit(code)
}

type result struct {
	Customer *entity.Customer
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

/*
	Fixtures: four customers
*/
func addFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.Customer{}) // cleaning customers
	tx.Create(customers)
}

func TestCustomerRepository(t *testing.T) {
	t.Run("FindAndLockByCustomerID", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Finding a customer", func(t *testing.T) {
				customerIDToFind := customers[1].CustomerID
				connection := database.GetStoriGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rCustomer := NewCustomerGormRepo(tx)

				got, err := rCustomer.FindByCustomerID(customerIDToFind)

				assert.NoError(t, err)
				assert.True(t, cmp.Equal(got, &customers[1], cmpopts.IgnoreTypes(time.Time{})))
				t.Cleanup(func() {
					tx.Rollback()
					connection.Unscoped().Where("1=1").Delete(&entity.Customer{})
				})
			})
		})
	})
	t.Run("FindAndLockByCustomerID", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Finding and locking a customer", func(t *testing.T) {
				// fixture
				customerIDToFind := customers[1].CustomerID
				channelResult1 := make(chan result)
				channelResult2 := make(chan result)

				connection := database.GetStoriGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rCustomer := NewCustomerGormRepo(tx)

				go func() {
					got, err := rCustomer.FindAndLockByCustomerID(customerIDToFind)
					channelResult1 <- result{got, err, time.Now()}
					time.Sleep(500 * time.Millisecond)
					tx.Commit()
				}()

				go func() {
					err2 := connection.Table("customer").Where("customer_id = ?", customerIDToFind).
						Update("email", "testtest@hotmail.com").Error
					channelResult2 <- result{nil, err2, time.Now()}
				}()

				// assertions
				result1 := <-channelResult1
				assert.NoError(t, result1.Err)
				assert.True(t, cmp.Equal(result1.Customer, &customers[1], cmpopts.IgnoreTypes(time.Time{})))
				// tries to update but still nil because the lock is not released
				result2 := <-channelResult2
				assert.Nil(t, result2.Customer)
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
				rCustomer := NewCustomerGormRepo(tx)

				// action
				got, err := rCustomer.FindAndLockByCustomerID(customerIDToFind)

				// assertions
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.EqualError(t, err, errors.ErrNotFound.Error())

				t.Cleanup(func() {
					tx.Rollback()
					connection.Unscoped().Where("1=1").Delete(&entity.Customer{})
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				rCustomer := NewCustomerGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Customer{}) //Cleaning customers
				tx.Migrator().DropTable(&entity.Customer{})

				got, err := rCustomer.FindAndLockByCustomerID(1)

				//Data Assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
		t.Run("Clone", func(t *testing.T) {
			db := database.GetStoriGormConnection()
			rCustomer := NewCustomerGormRepo(db)

			clone := rCustomer.Clone()

			assert.NotNil(t, clone)
			assert.Equal(t, rCustomer, clone)
		})
	})
}
