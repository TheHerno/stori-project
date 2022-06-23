package product

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/utils/helpers"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// setup
	database.SetupTrainingGormDB()
	code := m.Run()
	os.Exit(code)
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
	t.Run("Index", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			connection := database.GetTrainingGormConnection()
			tx := connection.Begin()
			addFixtures(tx)
			rProduct := NewProductGormRepo(tx)

			testCases := []struct {
				TestName          string
				Pagination        *dto.Pagination
				Expected          []entity.Product
				ExpectedPageCount int
			}{
				{
					TestName:          "All in one page",
					Pagination:        dto.NewPagination(1, 20, 0),
					Expected:          products,
					ExpectedPageCount: 1,
				},
				{
					TestName:          "With offset",
					Pagination:        dto.NewPagination(2, 2, 0),
					Expected:          products[2:],
					ExpectedPageCount: 2,
				},
				{
					TestName:          "With small page_size",
					Pagination:        dto.NewPagination(1, 1, 0),
					Expected:          products[:1],
					ExpectedPageCount: 4,
				},
				{
					TestName:          "With small page_size and second page",
					Pagination:        dto.NewPagination(2, 1, 0),
					Expected:          products[1:2],
					ExpectedPageCount: 4,
				},
			}

			for _, testCase := range testCases {
				t.Run(testCase.TestName, func(t *testing.T) {
					got, err := rProduct.Index(testCase.Pagination)

					// data assertion
					assert.NoError(t, err)
					assert.Len(t, got, len(testCase.Expected))
					assert.True(t, cmp.Equal(testCase.Expected, got, cmpopts.IgnoreTypes(time.Time{})))
					assert.Equal(t, int64(4), testCase.Pagination.TotalCount)
					assert.Equal(t, testCase.ExpectedPageCount, testCase.Pagination.PageCount())
				})
			}

			t.Cleanup(func() {
				tx.Rollback()
			})
		})

		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Table doesn't exists", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rProduct := NewProductGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Product{}) // cleaning the table
				tx.Migrator().DropTable(&entity.Product{})           // droping table

				pagination := dto.NewPagination(1, 2, 0)
				got, err := rProduct.Index(pagination)

				// data assert
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})

	t.Run("Update", func(t *testing.T) {
		// fixture
		newName := "Product 1 edited"
		newDescription := "Description editada 123"
		newEnabled := &falseValue
		productToUpdate := &entity.Product{}
		copier.Copy(productToUpdate, &products[0])
		productToUpdate.Name = newName
		productToUpdate.Description = &newDescription
		productToUpdate.Enabled = newEnabled

		t.Run("Should success on", func(t *testing.T) {
			connection := database.GetTrainingGormConnection()
			tx := connection.Begin()
			addFixtures(tx)
			rProduct := NewProductGormRepo(tx)

			got, err := rProduct.Update(productToUpdate)
			// data assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, productToUpdate, cmpopts.IgnoreTypes(time.Time{})))
			assert.Equal(t, &newDescription, got.Description)
			assert.Equal(t, newName, got.Name)
			assert.Equal(t, newEnabled, got.Enabled)
			assert.True(t, got.UpdatedAt.After(products[0].UpdatedAt))

			// database assertion
			productUpdated := &entity.Product{}
			tx.Take(productUpdated, got)
			assert.True(t, cmp.Equal(got, productUpdated, cmpopts.IgnoreTypes(time.Time{})))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})

		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Product validation fails", func(t *testing.T) {
				// fixture

				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				addFixtures(tx)
				rProduct := NewProductGormRepo(tx)

				// Fixture
				invalidProductToUpdate := &entity.Product{}
				copier.Copy(invalidProductToUpdate, &products[0])
				invalidProductToUpdate.Name = ""
				invalidProductToUpdate.Description = helpers.PointerToString(strings.Repeat("a", 501))
				invalidProductToUpdate.Slug = strings.Repeat("a", 306)
				got, err := rProduct.Update(invalidProductToUpdate)

				// data assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Product not found", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rProduct := NewProductGormRepo(tx)

				got, err := rProduct.Update(&entity.Product{
					ProductID: 999,
				})

				// data assertion
				assert.EqualError(t, err, errors.ErrNotFound.Error())
				assert.Nil(t, got)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Table doesn't exists", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rProduct := NewProductGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Product{}) // cleaning the table
				tx.Migrator().DropTable(&entity.Product{})

				got, err := rProduct.Update(productToUpdate)

				// data assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})

		t.Run("Create", func(t *testing.T) {
			// fixture
			productToCreate := &entity.Product{}
			copier.Copy(productToCreate, &products[0])

			t.Run("Should success on", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rProduct := NewProductGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Product{}) // cleaning the table

				got, err := rProduct.Create(productToCreate)

				// data assertion
				assert.NoError(t, err)
				assert.True(t, cmp.Equal(got, productToCreate, cmpopts.IgnoreTypes(time.Time{})))
				assert.NotZero(t, got.ProductID)
				assert.NotZero(t, got.CreatedAt)

				// database assertion
				productCreated := &entity.Product{}
				tx.Take(productCreated, got)
				assert.True(t, cmp.Equal(got, productCreated, cmpopts.IgnoreTypes(time.Time{})))

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Should fail on", func(t *testing.T) {
				t.Run("Product validation fails", func(t *testing.T) {
					// Fixture
					invalidProductToCreate := &entity.Product{}
					copier.Copy(invalidProductToCreate, &products[0])
					invalidProductToCreate.Name = ""
					invalidProductToCreate.Description = helpers.PointerToString(strings.Repeat("a", 501))
					invalidProductToCreate.Slug = strings.Repeat("a", 306)

					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rProduct := NewProductGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Product{}) // cleaning the table

					got, err := rProduct.Create(invalidProductToCreate)

					// data assertion
					assert.Error(t, err)
					assert.Nil(t, got)

					t.Cleanup(func() {
						tx.Rollback()
					})
				})
				t.Run("Table doesn't exists", func(t *testing.T) {
					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rProduct := NewProductGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Product{}) // cleaning the table
					tx.Migrator().DropTable(&entity.Product{})

					got, err := rProduct.Create(productToCreate)

					// data assertion
					assert.Nil(t, got)
					assert.Error(t, err)

					t.Cleanup(func() {
						tx.Rollback()
					})
				})
			})
		})

		t.Run("Delete", func(t *testing.T) {
			t.Run("Should success on", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				addFixtures(tx)
				rProduct := NewProductGormRepo(tx)

				err := rProduct.Delete(products[0].ProductID)
				// data assertion
				assert.NoError(t, err)

				// database assertion
				productDeleted := &entity.Product{}
				tx.Unscoped().Take(productDeleted, &entity.Product{ProductID: products[0].ProductID})
				assert.True(t, productDeleted.DeletedAt.Time.After(products[0].UpdatedAt))
				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Should fail on", func(t *testing.T) {
				t.Run("Product not found", func(t *testing.T) {
					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rProduct := NewProductGormRepo(tx)

					err := rProduct.Delete(999)

					// data assertion
					assert.EqualError(t, err, errors.ErrNotFound.Error())

					t.Cleanup(func() {
						tx.Rollback()
					})
				})
				t.Run("Table doesn't exists", func(t *testing.T) {
					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rProduct := NewProductGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Product{}) // cleaning the table
					tx.Migrator().DropTable(&entity.Product{})

					err := rProduct.Delete(products[0].ProductID)

					// data assertion
					assert.Error(t, err)

					t.Cleanup(func() {
						tx.Rollback()
					})
				})
			})
		})

		t.Run("FindByID", func(t *testing.T) {
			//Fixture
			id := products[0].ProductID
			t.Run("Should success on", func(t *testing.T) {
				tx := database.GetTrainingGormConnection().Begin()
				addFixtures(tx)
				rProduct := NewProductGormRepo(tx)

				got, err := rProduct.FindByID(id)

				//Data Assertion
				assert.NoError(t, err)
				assert.True(t, cmp.Equal(got, &products[0], cmpopts.IgnoreFields(entity.Product{}, "UpdatedAt", "CreatedAt")))

				t.Cleanup(func() {
					tx.Rollback()
				})
			})

			t.Run("Should fail on", func(t *testing.T) {
				t.Run("Product not found", func(t *testing.T) {
					//Fixture
					id := 9999999
					tx := database.GetTrainingGormConnection().Begin()
					addFixtures(tx)

					rProduct := NewProductGormRepo(tx)

					got, err := rProduct.FindByID(id)

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
					rProduct := NewProductGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Product{}) //Cleaning products
					tx.Migrator().DropTable(&entity.Product{})

					got, err := rProduct.FindByID(id)

					//Data Assertion
					assert.Nil(t, got)
					assert.Error(t, err)

					t.Cleanup(func() {
						tx.Rollback()
					})
				})
			})
		})

		t.Run("CountBySlug", func(t *testing.T) {
			t.Run("Should success on", func(t *testing.T) {
				tx := database.GetTrainingGormConnection().Begin()
				addFixtures(tx)
				rProduct := NewProductGormRepo(tx)

				got, err := rProduct.CountBySlug(products[0].Slug)

				//Data Assertion
				assert.NoError(t, err)
				assert.Equal(t, int64(1), got)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Should fail on", func(t *testing.T) {
				t.Run("Table doesn't exist", func(t *testing.T) {
					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rProduct := NewProductGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Product{}) //Cleaning products
					tx.Migrator().DropTable(&entity.Product{})

					got, err := rProduct.CountBySlug(products[0].Slug)

					//Data Assertion
					assert.Zero(t, got)
					assert.Error(t, err)

					t.Cleanup(func() {
						tx.Rollback()
					})
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
