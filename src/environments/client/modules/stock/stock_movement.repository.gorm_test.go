package stock

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/utils/constant"
	"stori-service/src/utils/helpers"
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

var (
	trueValue  = true
	falseValue = false
)

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
	{
		ProductID:   5,
		Name:        "Product 5",
		Slug:        "product-5",
		Description: helpers.PointerToString("Descripción 5"),
		Enabled:     &trueValue,
	},
}

var warehouses = []entity.Warehouse{
	{
		WarehouseID: 1,
		Name:        "Warehouse 1",
		Address:     "Address 1",
		UserID:      1,
	},
	{
		WarehouseID: 2,
		Name:        "Warehouse 2",
		Address:     "Address 2",
		UserID:      2,
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

var stockMovements = []entity.StockMovement{
	{
		StockMovementID: 1,
		ProductID:       1,
		WarehouseID:     1,
		Quantity:        10,
		Available:       10,
		Concept:         "Ingreso",
		Type:            constant.IncomeType,
	},
	{
		StockMovementID: 2,
		ProductID:       1,
		WarehouseID:     1,
		Quantity:        5,
		Available:       5,
		Concept:         "Venta",
		Type:            constant.OutcomeType,
	},
	{
		StockMovementID: 3,
		ProductID:       1,
		WarehouseID:     1,
		Quantity:        10,
		Available:       15,
		Concept:         "Ingreso",
		Type:            constant.IncomeType,
	},
	{
		StockMovementID: 4,
		ProductID:       2,
		WarehouseID:     4,
		Quantity:        17,
		Available:       17,
		Concept:         "Ingreso",
		Type:            constant.IncomeType,
	},
	{
		StockMovementID: 5,
		ProductID:       2,
		WarehouseID:     1,
		Quantity:        20,
		Available:       20,
		Concept:         "Ingreso",
		Type:            constant.IncomeType,
	},
	{
		StockMovementID: 6,
		ProductID:       3,
		WarehouseID:     1,
		Quantity:        100,
		Available:       100,
		Concept:         "Ingreso",
		Type:            constant.IncomeType,
	},
	{
		StockMovementID: 7,
		ProductID:       3,
		WarehouseID:     1,
		Quantity:        75,
		Available:       25,
		Concept:         "Venta",
		Type:            constant.OutcomeType,
	},
	{
		StockMovementID: 8,
		ProductID:       5,
		WarehouseID:     1,
		Quantity:        12,
		Available:       12,
		Concept:         "Ingreso",
		Type:            constant.IncomeType,
	},
}

/*
	Foreign Fixures
*/
func addForeignFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.Product{})   // cleaning products
	tx.Unscoped().Where("1=1").Delete(&entity.Warehouse{}) // cleaning warehouses
	tx.Create(products)
	tx.Create(warehouses)
}

/*
	Fixtures: four movements
*/
func addFixtures(tx *gorm.DB) {
	addForeignFixtures(tx)
	tx.Unscoped().Where("1=1").Delete(&entity.StockMovement{}) // cleaning movements
	tx.Create(stockMovements)
}

func TestGormRepository(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		// fixture
		movementToCreate := &entity.StockMovement{}
		copier.Copy(movementToCreate, stockMovements[0])
		t.Run("Should success on", func(t *testing.T) {
			connection := database.GetTrainingGormConnection()
			tx := connection.Begin()
			addForeignFixtures(tx)
			rStockMovement := NewStockMovementGormRepo(tx)
			tx.Unscoped().Where("1=1").Delete(&entity.StockMovement{}) // cleaning the table

			got, err := rStockMovement.Create(movementToCreate)

			// data assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, movementToCreate, cmpopts.IgnoreTypes(time.Time{})))
			assert.NotZero(t, got.StockMovementID)
			assert.NotZero(t, got.CreatedAt)

			// database assertion
			movementCreated := &entity.StockMovement{}
			tx.Take(movementCreated, got)
			assert.True(t, cmp.Equal(got, movementCreated, cmpopts.IgnoreTypes(time.Time{})))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Movement validation fails", func(t *testing.T) {
				// Fixture
				invalidMovementToCreate := &entity.StockMovement{}
				copier.Copy(invalidMovementToCreate, &stockMovements[0])
				invalidMovementToCreate.ProductID = 0
				invalidMovementToCreate.WarehouseID = 0
				invalidMovementToCreate.Quantity = 0
				invalidMovementToCreate.Available = -1
				invalidMovementToCreate.Concept = ""
				invalidMovementToCreate.Type = 0

				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rMovement := NewStockMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.StockMovement{}) // cleaning the table

				got, err := rMovement.Create(invalidMovementToCreate)

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
				rStockMovement := NewStockMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.StockMovement{}) // cleaning the table
				tx.Migrator().DropTable(&entity.StockMovement{})

				got, err := rStockMovement.Create(movementToCreate)

				// data assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("FindLastStockMovement", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			// fixture
			productIDToFind := products[1].ProductID
			warehouseIDToFind := warehouses[3].WarehouseID

			connection := database.GetTrainingGormConnection()
			tx := connection.Begin()
			addFixtures(tx)
			rStockMovement := NewStockMovementGormRepo(tx)

			got, err := rStockMovement.FindLastStockMovement(warehouseIDToFind, productIDToFind)

			// assertions
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, &stockMovements[3], cmpopts.IgnoreTypes(time.Time{})))

			t.Cleanup(func() {
				tx.Rollback()
				connection.Unscoped().Where("1=1").Delete(&entity.StockMovement{})
				connection.Unscoped().Where("1=1").Delete(&entity.Warehouse{})
				connection.Unscoped().Where("1=1").Delete(&entity.Product{})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Movement doesn't exists", func(t *testing.T) {
				// fixture
				productIDToFind := 78
				warehouseIDToFind := 87
				connection := database.GetTrainingGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rStockMovement := NewStockMovementGormRepo(tx)

				// action
				got, err := rStockMovement.FindLastStockMovement(warehouseIDToFind, productIDToFind)

				// assertions
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.EqualError(t, err, errors.ErrNotFound.Error())

				t.Cleanup(func() {
					tx.Rollback()
					connection.Unscoped().Where("1=1").Delete(&entity.StockMovement{})
					connection.Unscoped().Where("1=1").Delete(&entity.Warehouse{})
					connection.Unscoped().Where("1=1").Delete(&entity.Product{})
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				rStockMovement := NewStockMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.StockMovement{}) //Cleaning movements
				tx.Migrator().DropTable(&entity.StockMovement{})

				got, err := rStockMovement.FindLastStockMovement(1, 1)

				//Data Assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("FindStocksByWarehouse", func(t *testing.T) {
		// fixture
		warehouseIDToFind := warehouses[0].WarehouseID
		t.Run("Should success on", func(t *testing.T) {
			// fixture
			productsWithStockToFind := []dto.ProductWithStock{
				{
					ProductID:   products[0].ProductID,
					Description: products[0].Description,
					Name:        products[0].Name,
					Slug:        products[0].Slug,
					Stock:       stockMovements[2].Available,
				},
				{
					ProductID:   products[1].ProductID,
					Description: products[1].Description,
					Name:        products[1].Name,
					Slug:        products[1].Slug,
					Stock:       stockMovements[4].Available,
				},
				{
					ProductID:   products[2].ProductID,
					Description: products[2].Description,
					Name:        products[2].Name,
					Slug:        products[2].Slug,
					Stock:       stockMovements[6].Available,
				},
				{
					ProductID:   products[4].ProductID,
					Description: products[4].Description,
					Name:        products[4].Name,
					Slug:        products[4].Slug,
					Stock:       stockMovements[7].Available,
				},
			}

			connection := database.GetTrainingGormConnection()
			tx := connection.Begin()
			addFixtures(tx)
			rStockMovement := NewStockMovementGormRepo(tx)

			testCases := []struct {
				TestName          string
				Pagination        *dto.Pagination
				Expected          []dto.ProductWithStock
				ExpectedPageCount int
			}{
				{
					TestName:          "All in one page",
					Pagination:        dto.NewPagination(1, 20, 0),
					Expected:          productsWithStockToFind,
					ExpectedPageCount: 1,
				},
				{
					TestName:          "With offset",
					Pagination:        dto.NewPagination(2, 2, 0),
					Expected:          productsWithStockToFind[2:],
					ExpectedPageCount: 2,
				},
				{
					TestName:          "With small page_size",
					Pagination:        dto.NewPagination(1, 1, 0),
					Expected:          productsWithStockToFind[:1],
					ExpectedPageCount: 4,
				},
				{
					TestName:          "With small page_size and second page",
					Pagination:        dto.NewPagination(2, 1, 0),
					Expected:          productsWithStockToFind[1:2],
					ExpectedPageCount: 4,
				},
			}

			for _, tC := range testCases {
				t.Run(tC.TestName, func(t *testing.T) {
					got, err := rStockMovement.FindStocksByWarehouse(warehouseIDToFind, tC.Pagination)

					// assertions
					assert.NoError(t, err)
					assert.Len(t, got, len(tC.Expected))
					assert.Equal(t, tC.Expected, got)
					assert.Equal(t, int64(4), tC.Pagination.TotalCount)
					assert.Equal(t, tC.ExpectedPageCount, tC.Pagination.PageCount())
				})
			}

			t.Cleanup(func() {
				tx.Rollback()
				connection.Unscoped().Where("1=1").Delete(&entity.StockMovement{})
				connection.Unscoped().Where("1=1").Delete(&entity.Warehouse{})
				connection.Unscoped().Where("1=1").Delete(&entity.Product{})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetTrainingGormConnection()
				tx := connection.Begin()
				pagination := dto.NewPagination(1, 20, 0)
				rStockMovement := NewStockMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.StockMovement{}) //Cleaning movements
				tx.Migrator().DropTable(&entity.StockMovement{})

				got, err := rStockMovement.FindStocksByWarehouse(1, pagination)

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
		rStockMovement := NewStockMovementGormRepo(db)

		clone := rStockMovement.Clone()

		assert.NotNil(t, clone)
		assert.Equal(t, rStockMovement, clone)
	})
}
