package scopes

import (
	"os"
	"stori-service/src/libs/database"
	"stori-service/src/libs/dto"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	//setup
	database.SetupTrainingGormDB()
	code := m.Run() //run tests
	os.Exit(code)
}

func TestStockMovementScope(t *testing.T) {
	db := database.GetTrainingGormConnection().Session(&gorm.Session{DryRun: true})
	t.Run("StockMovementByWarehouseAndProductID", func(t *testing.T) {
		subQuery := db.Scopes(StockMovementByWarehouseAndProductID(1, 2)).Find(nil).Statement

		//Data Assertion: query
		assert.Contains(t, subQuery.SQL.String(), `SELECT * FROM "stock_movement"`)
		assert.Contains(t, subQuery.SQL.String(), "\"warehouse_id\" = $2")
		assert.Contains(t, subQuery.SQL.String(), "\"product_id\" = $1")

		//Data Assertion: inteporlated values
		assert.Equal(t, 1, subQuery.Vars[1])
		assert.Equal(t, 2, subQuery.Vars[0])
	})
	t.Run("StockMovementByProductID", func(t *testing.T) {
		subQuery := db.Scopes(StockMovementByProductID(2)).Find(nil).Statement

		//Data Assertion: query
		assert.Contains(t, subQuery.SQL.String(), `SELECT * FROM "stock_movement"`)
		assert.Contains(t, subQuery.SQL.String(), "\"product_id\" = $1")

		//Data Assertion: inteporlated values
		assert.Equal(t, 2, subQuery.Vars[0])
	})
	t.Run("StockMovementByWarehouseID", func(t *testing.T) {
		subQuery := db.Scopes(StockMovementByWarehouseID(1)).Find(nil).Statement

		//Data Assertion: query
		assert.Contains(t, subQuery.SQL.String(), `SELECT * FROM "stock_movement"`)
		assert.Contains(t, subQuery.SQL.String(), "\"warehouse_id\" = $1")

		//Data Assertion: inteporlated values
		assert.Equal(t, 1, subQuery.Vars[0])
	})
	t.Run("StocksByWarehouseID", func(t *testing.T) {
		pagination := dto.NewPagination(2, 20, 0)
		subQuery := db.Scopes(StocksByWarehouseID(1, pagination)).Find(nil).Statement

		//Data Assertion: query
		assert.Contains(t, subQuery.SQL.String(), `SELECT DISTINCT ON (stock_movement.product_id)`)
		assert.Contains(t, subQuery.SQL.String(), `product.product_id, stock_movement.available as stock, product.name, product.slug, product.description`)
		assert.Contains(t, subQuery.SQL.String(), `JOIN product ON product.product_id = stock_movement.product_id`)
		assert.Contains(t, subQuery.SQL.String(), `ORDER BY stock_movement.product_id ASC, stock_movement_id DESC`)
		assert.Contains(t, subQuery.SQL.String(), "\"warehouse_id\" = $1")
		assert.Contains(t, subQuery.SQL.String(), "enabled = true")
		assert.Contains(t, subQuery.SQL.String(), "LIMIT "+strconv.Itoa(pagination.PageSize))
		assert.Contains(t, subQuery.SQL.String(), "OFFSET "+strconv.Itoa(pagination.Offset()))
		//Data Assertion: inteporlated values
		assert.Equal(t, 1, subQuery.Vars[0])
	})
}
