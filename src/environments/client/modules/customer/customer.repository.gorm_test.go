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
		Customerid: 1,
		Name:       "Customer 1",
	},
	{
		Customerid: 2,
		Name:       "Customer 2",
	},
	{
		Customerid: 3,
		Name:       "Customer 3",
	},
	{
		Customerid: 4,
		Name:       "Customer 4",
	},
}

/*
	Fixtures: four customers
*/
func addFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.Customer{}) // cleaning customers
	tx.Create(customers)
}

func TestGormRepository(t *testing.T) {
	t.Run("FindByCustomerid", func(t *testing.T) {
		//Fixture
		customerid := customers[0].Customerid
		t.Run("Should success on", func(t *testing.T) {
			tx := database.GetStoriGormConnection().Begin()
			addFixtures(tx)
			rCustomer := NewCustomerGormRepo(tx)

			got, err := rCustomer.FindByCustomerid(customerid)

			//Data Assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, &customers[0], cmpopts.IgnoreFields(entity.Customer{}, "UpdatedAt", "CreatedAt")))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})

		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Customer not found", func(t *testing.T) {
				//Fixture
				customerid := 9999999
				tx := database.GetStoriGormConnection().Begin()
				addFixtures(tx)

				rCustomer := NewCustomerGormRepo(tx)

				got, err := rCustomer.FindByCustomerid(customerid)

				//Data Assertion
				assert.EqualError(t, err, errors.ErrNotFound.Error())
				assert.Nil(t, got)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				rCustomer := NewCustomerGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Customer{}) //Cleaning customers
				tx.Migrator().DropTable(&entity.Customer{})

				got, err := rCustomer.FindByCustomerid(customerid)

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
			customeridToFind := customers[1].Customerid
			channelResult1 := make(chan result)
			channelResult2 := make(chan result)

			connection := database.GetStoriGormConnection()
			addFixtures(connection)
			tx := connection.Begin()
			rCustomer := NewCustomerGormRepo(tx)

			go func() {
				got, err := rCustomer.FindAndLockByCustomerid(customeridToFind)
				channelResult1 <- result{got, err, time.Now()}
				time.Sleep(500 * time.Millisecond)
				tx.Commit()
			}()

			go func() {
				err2 := connection.Table("customer").Where("customer_id = ?", customeridToFind).
					Update("name", "Customeritooo").Error
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
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Customer doesn't exists", func(t *testing.T) {
				// fixture
				customeridToFind := 78
				connection := database.GetStoriGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rCustomer := NewCustomerGormRepo(tx)

				// action
				got, err := rCustomer.FindAndLockByCustomerid(customeridToFind)

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

				got, err := rCustomer.FindAndLockByCustomerid(1)

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
		rCustomer := NewCustomerGormRepo(db)

		clone := rCustomer.Clone()

		assert.NotNil(t, clone)
		assert.Equal(t, rCustomer, clone)
	})
}
