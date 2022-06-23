package product

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/libs/errors"
	"stori-service/src/utils/helpers"
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
	Product *entity.Product
	Err     error
	Date    time.Time
}

var trueValue = true
var falseValue = false
var products = []entity.Product{
	{
		ProductID:   1,
		Name:        "Product 1",
		Slug:        "product-1",
		Description: helpers.PointerToString("Descripción 1"),
		Enabled:     &trueValue,
	},
	{
		ProductID:   2,
		Name:        "Product 2",
		Slug:        "product-2",
		Description: helpers.PointerToString("Descripción 2"),
		Enabled:     &trueValue,
	},
	{
		ProductID:   3,
		Name:        "Product 3",
		Slug:        "product-3",
		Description: helpers.PointerToString("Descripción 3"),
		Enabled:     &trueValue,
	},
	{
		ProductID:   4,
		Name:        "Product 4",
		Slug:        "product-4",
		Description: helpers.PointerToString("Descripción 4"),
		Enabled:     &falseValue,
	},
}

/*
	Fixtures: four products
*/
func addFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.Product{}) // cleaning products
	tx.Create(products)
}

func TestGormRepository(t *testing.T) {
	t.Run("FindAndLockByID", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			// fixture
			productIDToFind := products[1].ProductID
			channelResult1 := make(chan result)
			channelResult2 := make(chan result)

			connection := database.GetTrainingGormConnection()
			addFixtures(connection)
			tx := connection.Begin()
			rProduct := NewProductGormRepo(tx)

			go func() {
				got, err := rProduct.FindAndLockByID(productIDToFind)
				channelResult1 <- result{got, err, time.Now()}
				time.Sleep(500 * time.Millisecond)
				tx.Commit()
			}()

			go func() {
				err2 := connection.Table("product").Where("product_id = ?", productIDToFind).
					Update("description", "Productitooo").Error
				channelResult2 <- result{nil, err2, time.Now()}
			}()

			// assertions
			result1 := <-channelResult1
			assert.NoError(t, result1.Err)
			assert.True(t, cmp.Equal(result1.Product, &products[1], cmpopts.IgnoreTypes(time.Time{})))
			// tries to update but still nil because the lock is not released
			result2 := <-channelResult2
			assert.Nil(t, result2.Product)
			assert.NoError(t, result2.Err)
			assert.True(t, result2.Date.After(result1.Date))

			t.Cleanup(func() {
				tx.Rollback()
				connection.Unscoped().Where("1=1").Delete(&entity.Product{})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Movement doesn't exists", func(t *testing.T) {
				// fixture
				productIDToFind := 78
				connection := database.GetTrainingGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rProduct := NewProductGormRepo(tx)

				// action
				got, err := rProduct.FindAndLockByID(productIDToFind)

				// assertions
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.EqualError(t, err, errors.ErrNotFound.Error())

				t.Cleanup(func() {
					tx.Rollback()
					connection.Unscoped().Where("1=1").Delete(&entity.Product{})
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rProduct := NewProductGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Product{}) //Cleaning movements
				tx.Migrator().DropTable(&entity.Product{})

				got, err := rProduct.FindAndLockByID(1)

				//Data Assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
		t.Run("Clone", func(t *testing.T) {
			db := database.GetTrainingGormConnection()
			rProduct := NewProductGormRepo(db)

			clone := rProduct.Clone()

			assert.NotNil(t, clone)
			assert.Equal(t, rProduct, clone)
		})
	})
}
