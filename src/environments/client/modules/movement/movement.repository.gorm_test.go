package movement

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/utils/constant"
)

/*
func TestMain(m *testing.M) {
	// setup
	database.SetupStoriGormDB()
	code := m.Run()
	os.Exit(code)
}

var (
	trueValue  = true
	falseValue = false
)
*/
var users = []entity.Customer{
	{
		Customerid: 1,
		Name:       "User 1",
	},
	{
		Customerid: 2,
		Name:       "User 2",
	},
	{
		Customerid: 3,
		Name:       "User 3",
	},
	{
		Customerid: 4,
		Name:       "User 4",
	},
}

var movements = []entity.Movement{
	{
		MovementID: 1,
		Customerid: 1,
		Quantity:   10,
		Available:  10,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 2,
		Customerid: 1,
		Quantity:   5,
		Available:  5,
		Type:       constant.OutcomeType,
	},
	{
		MovementID: 3,
		Customerid: 1,
		Quantity:   10,
		Available:  15,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 4,
		Customerid: 4,
		Quantity:   17,
		Available:  17,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 5,
		Customerid: 1,
		Quantity:   20,
		Available:  20,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 6,
		Customerid: 1,
		Quantity:   100,
		Available:  100,
		Type:       constant.IncomeType,
	},
	{
		MovementID: 7,
		Customerid: 1,
		Quantity:   75,
		Available:  25,
		Type:       constant.OutcomeType,
	},
	{
		MovementID: 8,
		Customerid: 1,
		Quantity:   12,
		Available:  12,
		Type:       constant.IncomeType,
	},
}

/*
	Foreign Fixures
*/
/*
func addForeignFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.User{}) // cleaning users
	tx.Create(users)
}
*/
/*
	Fixtures: four movements
*/
/*
func addFixtures(tx *gorm.DB) {
	addForeignFixtures(tx)
	tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning movements
	tx.Create(movements)
}
*/
/*
func TestGormRepository(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		// fixture
		movementToCreate := &entity.Movement{}
		copier.Copy(movementToCreate, movements[0])
		t.Run("Should success on", func(t *testing.T) {
			connection := database.GetStoriGormConnection()
			tx := connection.Begin()
			addForeignFixtures(tx)
			rMovement := NewMovementGormRepo(tx)
			tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning the table

			got, err := rMovement.Create(movementToCreate)

			// data assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, movementToCreate, cmpopts.IgnoreTypes(time.Time{})))
			assert.NotZero(t, got.MovementID)
			assert.NotZero(t, got.CreatedAt)

			// database assertion
			movementCreated := &entity.Movement{}
			tx.Take(movementCreated, got)
			assert.True(t, cmp.Equal(got, movementCreated, cmpopts.IgnoreTypes(time.Time{})))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Movement validation fails", func(t *testing.T) {
				// Fixture
				invalidMovementToCreate := &entity.Movement{}
				copier.Copy(invalidMovementToCreate, &movements[0])
				invalidMovementToCreate.Customerid = 0
				invalidMovementToCreate.Quantity = 0
				invalidMovementToCreate.Available = -1
				invalidMovementToCreate.Type = 0

				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				rMovement := NewMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning the table

				got, err := rMovement.Create(invalidMovementToCreate)

				// data assertion
				assert.Error(t, err)
				assert.Nil(t, got)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Table doesn't exists", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				rMovement := NewMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) // cleaning the table
				tx.Migrator().DropTable(&entity.Movement{})

				got, err := rMovement.Create(movementToCreate)

				// data assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("FindStocksByUser", func(t *testing.T) {
		// fixture
		customeridToFind := users[0].Customerid
		t.Run("Should success on", func(t *testing.T) {
			// fixture
			productsWithStockToFind := []dto.ProductWithStock{
				{
					ProductID:   products[0].ProductID,
					Description: products[0].Description,
					Name:        products[0].Name,
					Slug:        products[0].Slug,
					Stock:       movements[2].Available,
				},
				{
					ProductID:   products[1].ProductID,
					Description: products[1].Description,
					Name:        products[1].Name,
					Slug:        products[1].Slug,
					Stock:       movements[4].Available,
				},
				{
					ProductID:   products[2].ProductID,
					Description: products[2].Description,
					Name:        products[2].Name,
					Slug:        products[2].Slug,
					Stock:       movements[6].Available,
				},
				{
					ProductID:   products[4].ProductID,
					Description: products[4].Description,
					Name:        products[4].Name,
					Slug:        products[4].Slug,
					Stock:       movements[7].Available,
				},
			}

			connection := database.GetStoriGormConnection()
			tx := connection.Begin()
			addFixtures(tx)
			rMovement := NewMovementGormRepo(tx)

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
					got, err := rMovement.FindStocksByUser(customeridToFind, tC.Pagination)

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
				connection.Unscoped().Where("1=1").Delete(&entity.Movement{})
				connection.Unscoped().Where("1=1").Delete(&entity.User{})
				connection.Unscoped().Where("1=1").Delete(&entity.Product{})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				pagination := dto.NewPagination(1, 20, 0)
				rMovement := NewMovementGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.Movement{}) //Cleaning movements
				tx.Migrator().DropTable(&entity.Movement{})

				got, err := rMovement.FindStocksByUser(1, pagination)

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
*/
