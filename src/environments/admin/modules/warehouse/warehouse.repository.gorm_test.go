package warehouse

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
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
	t.Run("Index", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			connection := database.GetTrainingGormConnection()
			tx := connection.Begin()
			addFixtures(tx)
			rWarehouse := NewWarehouseGormRepo(tx)

			testCases := []struct {
				TestName          string
				Pagination        *dto.Pagination
				Expected          []entity.Warehouse
				ExpectedPageCount int
			}{
				{
					TestName:          "All in one page",
					Pagination:        dto.NewPagination(1, 20, 0),
					Expected:          warehouses,
					ExpectedPageCount: 1,
				},
				{
					TestName:          "With offset",
					Pagination:        dto.NewPagination(2, 2, 0),
					Expected:          warehouses[2:],
					ExpectedPageCount: 2,
				},
				{
					TestName:          "With small page_size",
					Pagination:        dto.NewPagination(1, 1, 0),
					Expected:          warehouses[:1],
					ExpectedPageCount: 4,
				},
				{
					TestName:          "With small page_size and second page",
					Pagination:        dto.NewPagination(2, 1, 0),
					Expected:          warehouses[1:2],
					ExpectedPageCount: 4,
				},
			}

			for _, testCase := range testCases {
				t.Run(testCase.TestName, func(t *testing.T) {
					got, err := rWarehouse.Index(testCase.Pagination)

					// data assertion
					assert.NoError(t, err)
					assert.Len(t, *got, len(testCase.Expected))
					assert.True(t, cmp.Equal(testCase.Expected, *got, cmpopts.IgnoreTypes(time.Time{})))
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
				rWarehouse := NewWarehouseGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning the table
				tx.Migrator().DropTable(&entity.Warehouse{})           // droping table

				pagination := dto.NewPagination(1, 2, 0)
				got, err := rWarehouse.Index(pagination)

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
		newName := "Warehouse 1 edited"
		newAddress := "Address editada 123"
		warehouseToUpdate := &entity.Warehouse{}
		copier.Copy(warehouseToUpdate, &warehouses[0])
		warehouseToUpdate.Name = newName
		warehouseToUpdate.Address = newAddress

		t.Run("Should success on", func(t *testing.T) {
			connection := database.GetTrainingGormConnection()
			tx := connection.Begin()
			addFixtures(tx)
			rWarehouse := NewWarehouseGormRepo(tx)

			got, err := rWarehouse.Update(warehouseToUpdate)
			// data assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, warehouseToUpdate, cmpopts.IgnoreTypes(time.Time{})))
			assert.Equal(t, newAddress, got.Address)
			assert.Equal(t, newName, got.Name)
			assert.True(t, got.UpdatedAt.After(warehouses[0].UpdatedAt))

			// database assertion
			warehouseUpdated := &entity.Warehouse{}
			tx.Take(warehouseUpdated, got)
			assert.True(t, cmp.Equal(got, warehouseUpdated, cmpopts.IgnoreTypes(time.Time{})))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})

		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Warehouse validation fails", func(t *testing.T) {
				// fixture

				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				addFixtures(tx)
				rWarehouse := NewWarehouseGormRepo(tx)

				// Fixture
				invalidWarehouseToUpdate := &entity.Warehouse{}
				copier.Copy(invalidWarehouseToUpdate, &warehouses[0])
				invalidWarehouseToUpdate.Name = ""
				invalidWarehouseToUpdate.Address = strings.Repeat("a", 301)

				got, err := rWarehouse.Update(invalidWarehouseToUpdate)

				// data assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Warehouse not found", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rWarehouse := NewWarehouseGormRepo(tx)

				got, err := rWarehouse.Update(&entity.Warehouse{
					WarehouseID: 999,
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
				rWarehouse := NewWarehouseGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning the table
				tx.Migrator().DropTable(&entity.Warehouse{})

				got, err := rWarehouse.Update(warehouseToUpdate)

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
			warehouseToCreate := &entity.Warehouse{}
			copier.Copy(warehouseToCreate, &warehouses[0])

			t.Run("Should success on", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rWarehouse := NewWarehouseGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning the table

				got, err := rWarehouse.Create(warehouseToCreate)

				// data assertion
				assert.NoError(t, err)
				assert.True(t, cmp.Equal(got, warehouseToCreate, cmpopts.IgnoreTypes(time.Time{})))
				assert.NotZero(t, got.WarehouseID)
				assert.NotZero(t, got.CreatedAt)

				// database assertion
				warehouseCreated := &entity.Warehouse{}
				tx.Take(warehouseCreated, got)
				assert.True(t, cmp.Equal(got, warehouseCreated, cmpopts.IgnoreTypes(time.Time{})))

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Should fail on", func(t *testing.T) {
				t.Run("Repeated User ID", func(t *testing.T) {
					// fixture
					warehouseInvalidID := &entity.Warehouse{}
					copier.Copy(warehouseInvalidID, &warehouses[0])
					warehouseInvalidID.WarehouseID = 8
					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rWarehouse := NewWarehouseGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning the table
					addFixtures(tx)

					got, err := rWarehouse.Create(warehouseInvalidID)

					// data assertion
					assert.Nil(t, got)
					assert.Error(t, err)

					t.Cleanup(func() {
						tx.Rollback()
					})
				})
				t.Run("Warehouse validation fails", func(t *testing.T) {
					// Fixture
					invalidWarehouseToCreate := &entity.Warehouse{}
					copier.Copy(invalidWarehouseToCreate, &warehouses[0])
					invalidWarehouseToCreate.Name = ""
					invalidWarehouseToCreate.Address = strings.Repeat("a", 301)

					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rWarehouse := NewWarehouseGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning the table

					got, err := rWarehouse.Create(invalidWarehouseToCreate)

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
					rWarehouse := NewWarehouseGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning the table
					tx.Migrator().DropTable(&entity.Warehouse{})

					got, err := rWarehouse.Create(warehouseToCreate)

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
				rWarehouse := NewWarehouseGormRepo(tx)

				err := rWarehouse.Delete(warehouses[0].WarehouseID)
				// data assertion
				assert.NoError(t, err)

				// database assertion
				warehouseDeleted := &entity.Warehouse{}
				tx.Unscoped().Take(warehouseDeleted, &entity.Warehouse{WarehouseID: warehouses[0].WarehouseID})
				assert.True(t, warehouseDeleted.DeletedAt.Time.After(warehouses[0].UpdatedAt))
				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Should fail on", func(t *testing.T) {
				t.Run("Warehouse not found", func(t *testing.T) {
					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rWarehouse := NewWarehouseGormRepo(tx)

					err := rWarehouse.Delete(999)

					// data assertion
					assert.EqualError(t, err, errors.ErrNotFound.Error())

					t.Cleanup(func() {
						tx.Rollback()
					})
				})
				t.Run("Table doesn't exists", func(t *testing.T) {
					connection := database.GetTrainingGormConnection()
					tx := connection.Begin()
					rWarehouse := NewWarehouseGormRepo(tx)
					tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning the table
					tx.Migrator().DropTable(&entity.Warehouse{})

					err := rWarehouse.Delete(warehouses[0].WarehouseID)

					// data assertion
					assert.Error(t, err)

					t.Cleanup(func() {
						tx.Rollback()
					})
				})
			})
		})

		t.Run("Clone", func(t *testing.T) {
			db := database.GetTrainingGormConnection()
			rProduct := NewWarehouseGormRepo(db)

			clone := rProduct.Clone()

			assert.NotNil(t, clone)
			assert.Equal(t, rProduct, clone)
		})
	})
}
