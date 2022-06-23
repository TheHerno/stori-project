package stock

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/libs/errors"
	"stori-service/src/utils/helpers"
	"testing"

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
		Type:            1,
	},
	{
		StockMovementID: 2,
		ProductID:       1,
		WarehouseID:     1,
		Quantity:        5,
		Available:       5,
		Concept:         "Venta",
		Type:            -1,
	},
	{
		StockMovementID: 3,
		ProductID:       1,
		WarehouseID:     1,
		Quantity:        10,
		Available:       15,
		Concept:         "Ingreso",
		Type:            1,
	},
	{
		StockMovementID: 5,
		ProductID:       2,
		WarehouseID:     4,
		Quantity:        17,
		Available:       17,
		Concept:         "Ingreso",
		Type:            1,
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
	t.Run("FindLastMovementByProductID", func(t *testing.T) {
		//Fixture
		productID := stockMovements[3].ProductID
		t.Run("Should success on", func(t *testing.T) {
			tx := database.GetTrainingGormConnection().Begin()
			addFixtures(tx)
			rStockMovement := NewStockMovementGormRepo(tx)

			got, err := rStockMovement.FindLastMovementByProductID(productID)

			//Data Assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, &stockMovements[3], cmpopts.IgnoreFields(entity.StockMovement{}, "UpdatedAt", "CreatedAt")))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})

		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Invalid ID", func(t *testing.T) {
				//Fixture
				var productID int
				tx := database.GetTrainingGormConnection().Begin()
				addFixtures(tx)

				rStockMovement := NewStockMovementGormRepo(tx)

				got, err := rStockMovement.FindLastMovementByProductID(productID)

				//Data Assertion
				assert.EqualError(t, err, errors.ErrNotFound.Error())
				assert.Nil(t, got)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Movement not found", func(t *testing.T) {
				//Fixture
				productID := 9999999
				tx := database.GetTrainingGormConnection().Begin()
				addFixtures(tx)

				rStockMovement := NewStockMovementGormRepo(tx)

				got, err := rStockMovement.FindLastMovementByProductID(productID)

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
				rStockMovement := NewStockMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.StockMovement{}) //Cleaning products
				tx.Migrator().DropTable(&entity.StockMovement{})

				got, err := rStockMovement.FindLastMovementByProductID(productID)

				//Data Assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})

	t.Run("FindLastMovementByWarehouseID", func(t *testing.T) {
		//Fixture
		warehouseID := stockMovements[3].WarehouseID
		t.Run("Should success on", func(t *testing.T) {
			tx := database.GetTrainingGormConnection().Begin()
			addFixtures(tx)
			rStockMovement := NewStockMovementGormRepo(tx)

			got, err := rStockMovement.FindLastMovementByWarehouseID(warehouseID)

			//Data Assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, &stockMovements[3], cmpopts.IgnoreFields(entity.StockMovement{}, "UpdatedAt", "CreatedAt")))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})

		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Invalid ID", func(t *testing.T) {
				//Fixture
				var warehouseID int
				tx := database.GetTrainingGormConnection().Begin()
				addFixtures(tx)

				rStockMovement := NewStockMovementGormRepo(tx)

				got, err := rStockMovement.FindLastMovementByWarehouseID(warehouseID)

				//Data Assertion
				assert.EqualError(t, err, errors.ErrNotFound.Error())
				assert.Nil(t, got)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Movement not found", func(t *testing.T) {
				//Fixture
				warehouseID := 9999999
				tx := database.GetTrainingGormConnection().Begin()
				addFixtures(tx)

				rStockMovement := NewStockMovementGormRepo(tx)

				got, err := rStockMovement.FindLastMovementByWarehouseID(warehouseID)

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
				rStockMovement := NewStockMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.StockMovement{}) //Cleaning products
				tx.Migrator().DropTable(&entity.StockMovement{})

				got, err := rStockMovement.FindLastMovementByWarehouseID(warehouseID)

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
