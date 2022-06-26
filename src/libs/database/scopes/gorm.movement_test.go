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
	database.SetupStoriGormDB()
	code := m.Run() //run tests
	os.Exit(code)
}

func TestMovementScope(t *testing.T) {
	db := database.GetStoriGormConnection().Session(&gorm.Session{DryRun: true})
	t.Run("MovementByCustomerID", func(t *testing.T) {
		subQuery := db.Scopes(MovementByCustomerID(1)).Find(nil).Statement

		//Data Assertion: query
		assert.Contains(t, subQuery.SQL.String(), `SELECT * FROM "movement"`)
		assert.Contains(t, subQuery.SQL.String(), "\"customer_id\" = $1")

		//Data Assertion: inteporlated values
		assert.Equal(t, 1, subQuery.Vars[0])
	})
	t.Run("StocksByCustomerID", func(t *testing.T) {
		pagination := dto.NewPagination(2, 20, 0)
		subQuery := db.Scopes(StocksByCustomerID(1, pagination)).Find(nil).Statement

		//Data Assertion: query
		assert.Contains(t, subQuery.SQL.String(), `SELECT DISTINCT ON (movement.product_id)`)
		assert.Contains(t, subQuery.SQL.String(), `product.product_id, movement.available as stock, product.name, product.slug, product.description`)
		assert.Contains(t, subQuery.SQL.String(), `JOIN product ON product.product_id = movement.product_id`)
		assert.Contains(t, subQuery.SQL.String(), `ORDER BY movement.product_id ASC, movement_id DESC`)
		assert.Contains(t, subQuery.SQL.String(), "\"customer_id\" = $1")
		assert.Contains(t, subQuery.SQL.String(), "enabled = true")
		assert.Contains(t, subQuery.SQL.String(), "LIMIT "+strconv.Itoa(pagination.PageSize))
		assert.Contains(t, subQuery.SQL.String(), "OFFSET "+strconv.Itoa(pagination.Offset()))
		//Data Assertion: inteporlated values
		assert.Equal(t, 1, subQuery.Vars[0])
	})
}
