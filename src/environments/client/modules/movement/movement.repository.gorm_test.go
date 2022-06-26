package movement

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/utils/constant"
	"testing"

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
	Fixtures: four movements
*/

func addFixtures(tx *gorm.DB) {
	addForeignFixtures(tx)
	tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning movements
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
	t.Run("FindLastMovement", func(t *testing.T) {
		//Fixture
		customerID := 1
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Finding last movement", func(t *testing.T) {
				tx := database.GetStoriGormConnection().Begin()
				addFixtures(tx)
				rMovement := NewMovementGormRepo(tx)

				got, err := rMovement.FindLastMovementByCustomerID(customerID)

				//Data Assertion
				assert.NoError(t, err)
				assert.True(t, cmp.Equal(got, &movements[3], cmpopts.IgnoreFields(entity.Movement{}, "UpdatedAt", "CreatedAt")))

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Table doesn't exists", func(t *testing.T) {
				tx := database.GetStoriGormConnection().Begin()
				rMovement := NewMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning the table
				tx.Migrator().DropTable(&entity.Movement{})

				got, err := rMovement.FindLastMovementByCustomerID(customerID)

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
