package warehouse

import (
	goerrors "errors"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	customMock "stori-service/src/utils/test/mock"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWarehouseService(t *testing.T) {
	// Fixture
	repositoryErr := goerrors.New("repository error")
	commitErr := goerrors.New("commit error")
	t.Run("Index", func(t *testing.T) {
		// fixture
		pagination := dto.NewPagination(1, 20, 0)
		expectedWarehouses := []entity.Warehouse{
			warehouses[0],
			warehouses[1],
		}

		t.Run("Should success on", func(t *testing.T) {
			mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
			mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)

			// expectations
			mockWarehouseRepo.On("Index", pagination).Return(&expectedWarehouses, nil)

			// action
			result, err := sWarehouse.Index(pagination)

			// mock assertion
			mockWarehouseRepo.AssertExpectations(t)
			mockCWarehouseRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockWarehouseRepo.AssertNumberOfCalls(t, "Index", 1)

			// data assertion
			assert.Equal(t, &expectedWarehouses, result)
			assert.NoError(t, err)
		})

		t.Run("Should fail on", func(t *testing.T) {
			mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
			mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)

			// Expectations
			mockWarehouseRepo.On("Index", pagination).Return(nil, repositoryErr)

			// action
			result, err := sWarehouse.Index(pagination)

			// mock assertion
			mockWarehouseRepo.AssertExpectations(t)
			mockCWarehouseRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockWarehouseRepo.AssertNumberOfCalls(t, "Index", 1)

			// data assertion
			assert.Nil(t, result)
			assert.Error(t, err, repositoryErr.Error())
		})
	})

	t.Run("Update", func(t *testing.T) {
		// Fixture
		warehouse := &warehouses[0]
		warehouseToUpdate := &dto.UpdateWarehouse{
			Name:        "New Name",
			WarehouseID: 1,
			Address:     "New Address",
		}
		warehouseUpdated := &entity.Warehouse{}
		copier.Copy(warehouseUpdated, warehouse)
		warehouseUpdated.Address = warehouseToUpdate.Address
		warehouseUpdated.Name = warehouseToUpdate.Name

		t.Run("Should success on", func(t *testing.T) {
			mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
			mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)

			// Expectations
			mockCWarehouseRepo.On("FindByID", warehouse.WarehouseID).Return(warehouse, nil)
			mockWarehouseRepo.On("Update", warehouseUpdated).Return(warehouseUpdated, nil)

			// Action
			result, err := sWarehouse.Update(warehouseToUpdate)

			// Mock Assertion
			mockWarehouseRepo.AssertExpectations(t)
			mockCWarehouseRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockCWarehouseRepo.AssertNumberOfCalls(t, "FindByID", 1)
			mockWarehouseRepo.AssertNumberOfCalls(t, "Update", 1)

			// Data Assertion
			assert.NoError(t, err)
			assert.Equal(t, warehouseUpdated, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      *dto.UpdateWarehouse
				mock       func(*customMock.AdminWarehouseRepository, *customMock.CommonWarehouseRepository)
				assertMock func(*customMock.AdminWarehouseRepository, *customMock.CommonWarehouseRepository)
			}{
				{
					name:  "Invalid update data",
					input: &dto.UpdateWarehouse{WarehouseID: 1},
					mock:  func(*customMock.AdminWarehouseRepository, *customMock.CommonWarehouseRepository) {},
					assertMock: func(mockWarehouseR *customMock.AdminWarehouseRepository, mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.AssertNumberOfCalls(t, "FindByID", 0)
						mockWarehouseR.AssertNumberOfCalls(t, "Update", 0)
					},
				},
				{
					name:  "Invalid WarehouseID",
					input: &dto.UpdateWarehouse{WarehouseID: 0, Name: "New Name", Address: "New Address"},
					mock:  func(*customMock.AdminWarehouseRepository, *customMock.CommonWarehouseRepository) {},
					assertMock: func(mockWarehouseR *customMock.AdminWarehouseRepository, mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.AssertNumberOfCalls(t, "FindByID", 0)
						mockWarehouseR.AssertNumberOfCalls(t, "Update", 0)
					},
				},
				{
					name:  "Fail Finding ID",
					input: warehouseToUpdate,
					mock: func(mockWarehouseR *customMock.AdminWarehouseRepository, mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.On("FindByID", warehouse.WarehouseID).Return(nil, repositoryErr)
					},
					assertMock: func(mockWarehouseR *customMock.AdminWarehouseRepository, mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.AssertNumberOfCalls(t, "FindByID", 1)
						mockWarehouseR.AssertNumberOfCalls(t, "Update", 0)
					},
				},
				{
					name:  "Fail to Update",
					input: warehouseToUpdate,
					mock: func(mockWarehouseR *customMock.AdminWarehouseRepository, mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.On("FindByID", warehouse.WarehouseID).Return(warehouse, nil)
						mockWarehouseR.On("Update", warehouse).Return(nil, repositoryErr)
					},
					assertMock: func(mockWarehouseR *customMock.AdminWarehouseRepository, mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.AssertNumberOfCalls(t, "FindByID", 1)
						mockWarehouseR.AssertNumberOfCalls(t, "Update", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
					mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)

					// Expectations
					tC.mock(mockWarehouseRepo, mockCWarehouseRepo)

					// Action
					result, err := sWarehouse.Update(tC.input)

					// Mock Assertion
					mockWarehouseRepo.AssertExpectations(t)
					tC.assertMock(mockWarehouseRepo, mockCWarehouseRepo)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(mockWarehouseRepo, mockCWarehouseRepo)

					// Data Assertion
					assert.Error(t, err)
					assert.Nil(t, result)
				})
			}
		})
	})

	t.Run("Create", func(t *testing.T) {
		// Fixture
		warehouse := &warehouses[0]
		dtoWarehouse := &dto.CreateWarehouse{
			Name:    warehouse.Name,
			Address: warehouse.Address,
			UserID:  1,
		}
		wareHouseToCreate := &entity.Warehouse{
			Name:    warehouse.Name,
			Address: warehouse.Address,
			UserID:  1,
		}
		warehouseCreated := &entity.Warehouse{}
		warehouseCreated.WarehouseID = 0
		copier.Copy(warehouseCreated, warehouse)
		t.Run("Should success on", func(t *testing.T) {
			mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
			mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)
			// expectations
			mockWarehouseRepo.On("Create", wareHouseToCreate).Return(warehouseCreated, nil)

			// actions
			result, err := sWarehouse.Create(dtoWarehouse)

			// mock assertion
			mockWarehouseRepo.AssertExpectations(t)
			mockCWarehouseRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)

			// data assertion
			assert.NoError(t, err)
			assert.Equal(t, warehouseCreated, result)
		})

		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      *dto.CreateWarehouse
				mock       func(*customMock.AdminWarehouseRepository)
				assertMock func(*testing.T, *customMock.AdminWarehouseRepository)
			}{
				{
					name:  "Invalid create data",
					input: &dto.CreateWarehouse{},
					mock:  func(*customMock.AdminWarehouseRepository) {},
					assertMock: func(t *testing.T, mockWarehouseR *customMock.AdminWarehouseRepository) {
						mockWarehouseR.AssertNumberOfCalls(t, "Create", 0)
					},
				},
				{
					name:  "Create",
					input: dtoWarehouse,
					mock: func(mockWarehouseR *customMock.AdminWarehouseRepository) {
						mockWarehouseR.On("Create", wareHouseToCreate).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockWarehouseR *customMock.AdminWarehouseRepository) {
						mockWarehouseR.AssertNumberOfCalls(t, "Create", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
					mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)

					// Expectations
					tC.mock(mockWarehouseRepo)

					// Action
					result, err := sWarehouse.Create(tC.input)

					// Mock Assertion
					mockWarehouseRepo.AssertExpectations(t)
					tC.assertMock(t, mockWarehouseRepo)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(t, mockWarehouseRepo)

					// Data assertion
					assert.Error(t, err)
					assert.Nil(t, result)
				})
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		// Fixture
		idToDelete := warehouses[0].WarehouseID
		t.Run("Should success on", func(t *testing.T) {
			mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
			mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)
			// expectations
			mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
			mockCWarehouseRepo.On("Clone").Return(mockCWarehouseRepo)
			mockWarehouseRepo.On("Clone").Return(mockWarehouseRepo)
			mockWarehouseRepo.On("Begin", nil).Return(nil)
			mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
			mockCWarehouseRepo.On("Begin", mock.Anything).Return(nil)
			mockWarehouseRepo.On("Rollback").Return(nil)
			mockCWarehouseRepo.On("FindAndLockByID", idToDelete).Return(&warehouses[0], nil)
			mockStockMovementRepo.On("FindLastMovementByWarehouseID", idToDelete).Return(nil, errors.ErrNotFound)
			mockWarehouseRepo.On("Commit").Return(nil)
			mockWarehouseRepo.On("Delete", idToDelete).Return(nil)

			// actions
			err := sWarehouse.Delete(idToDelete)

			// mock assertion
			mockWarehouseRepo.AssertExpectations(t)
			mockCWarehouseRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
			mockWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
			mockCWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
			mockWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
			mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
			mockCWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
			mockWarehouseRepo.AssertNumberOfCalls(t, "Rollback", 1)
			mockCWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
			mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByWarehouseID", 1)
			mockWarehouseRepo.AssertNumberOfCalls(t, "Commit", 1)
			mockWarehouseRepo.AssertNumberOfCalls(t, "Delete", 1)

			// data assertion
			assert.NoError(t, err)
		})
		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      int
				mock       func(*customMock.AdminWarehouseRepository, *customMock.CommonWarehouseRepository, *customMock.AdminStockMovementRepository)
				assertMock func(*testing.T, *customMock.AdminWarehouseRepository, *customMock.CommonWarehouseRepository, *customMock.AdminStockMovementRepository)
			}{
				{
					name:  "Invalid id",
					input: 0,
					mock: func(mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
					},
					assertMock: func(t *testing.T, mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Rollback", 0)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByWarehouseID", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Delete", 0)
					},
				},
				{
					name:  "Delete",
					input: idToDelete,
					mock: func(mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockWarehouseRepo.On("Clone").Return(mockWarehouseRepo)
						mockCWarehouseRepo.On("Clone").Return(mockCWarehouseRepo)
						mockWarehouseRepo.On("Begin", nil).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockCWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockWarehouseRepo.On("Rollback").Return(nil)
						mockCWarehouseRepo.On("FindAndLockByID", idToDelete).Return(&warehouses[0], nil)
						mockStockMovementRepo.On("FindLastMovementByWarehouseID", idToDelete).Return(nil, errors.ErrNotFound)
						mockWarehouseRepo.On("Delete", idToDelete).Return(repositoryErr)
					},
					assertMock: func(t *testing.T, mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByWarehouseID", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Delete", 1)
					},
				},
				{
					name:  "FindAndLockByID",
					input: idToDelete,
					mock: func(mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockWarehouseRepo.On("Clone").Return(mockWarehouseRepo)
						mockCWarehouseRepo.On("Clone").Return(mockCWarehouseRepo)
						mockWarehouseRepo.On("Begin", nil).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockCWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockWarehouseRepo.On("Rollback").Return(nil)
						mockCWarehouseRepo.On("FindAndLockByID", idToDelete).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByWarehouseID", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Delete", 0)
					},
				},
				{
					name:  "It has stock movement",
					input: idToDelete,
					mock: func(mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockWarehouseRepo.On("Clone").Return(mockWarehouseRepo)
						mockCWarehouseRepo.On("Clone").Return(mockCWarehouseRepo)
						mockWarehouseRepo.On("Begin", nil).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockCWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockWarehouseRepo.On("Rollback").Return(nil)
						mockCWarehouseRepo.On("FindAndLockByID", idToDelete).Return(&warehouses[0], nil)
						mockStockMovementRepo.On("FindLastMovementByWarehouseID", idToDelete).Return(&entity.StockMovement{}, nil)
					},
					assertMock: func(t *testing.T, mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByWarehouseID", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Delete", 0)
					},
				},
				{
					name:  "Error getting stock movement",
					input: idToDelete,
					mock: func(mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockWarehouseRepo.On("Clone").Return(mockWarehouseRepo)
						mockCWarehouseRepo.On("Clone").Return(mockCWarehouseRepo)
						mockWarehouseRepo.On("Begin", nil).Return(nil)
						mockCWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockWarehouseRepo.On("Rollback").Return(nil)
						mockCWarehouseRepo.On("FindAndLockByID", idToDelete).Return(&warehouses[0], nil)
						mockStockMovementRepo.On("FindLastMovementByWarehouseID", idToDelete).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByWarehouseID", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Delete", 0)
					},
				},
				{
					name:  "Commit",
					input: idToDelete,
					mock: func(mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockWarehouseRepo.On("Clone").Return(mockWarehouseRepo)
						mockCWarehouseRepo.On("Clone").Return(mockCWarehouseRepo)
						mockWarehouseRepo.On("Begin", nil).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockCWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockWarehouseRepo.On("Rollback").Return(nil)
						mockCWarehouseRepo.On("FindAndLockByID", idToDelete).Return(&warehouses[0], nil)
						mockStockMovementRepo.On("FindLastMovementByWarehouseID", idToDelete).Return(nil, errors.ErrNotFound)
						mockWarehouseRepo.On("Delete", idToDelete).Return(nil)
						mockWarehouseRepo.On("Commit").Return(commitErr)
					},
					assertMock: func(t *testing.T, mockWarehouseRepo *customMock.AdminWarehouseRepository, mockCWarehouseRepo *customMock.CommonWarehouseRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByWarehouseID", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Commit", 1)
						mockWarehouseRepo.AssertNumberOfCalls(t, "Delete", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
					mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)
					// Expectations
					tC.mock(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)
					// Action
					err := sWarehouse.Delete(tC.input)
					// mock assertion
					mockWarehouseRepo.AssertExpectations(t)
					mockCWarehouseRepo.AssertExpectations(t)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(t, mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)
					// Data assertion
					assert.Error(t, err)
				})
			}
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		// Fixture
		warehouseToFind := &warehouses[0]
		idToFind := warehouses[0].WarehouseID
		t.Run("Should success on", func(t *testing.T) {
			mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
			mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)
			// expectations
			mockCWarehouseRepo.On("FindByID", idToFind).Return(warehouseToFind, nil)

			// actions
			result, err := sWarehouse.FindByID(idToFind)

			// mock assertion
			mockWarehouseRepo.AssertExpectations(t)
			mockCWarehouseRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockCWarehouseRepo.AssertNumberOfCalls(t, "FindByID", 1)

			// data assertion
			assert.NoError(t, err)
			assert.Equal(t, warehouseToFind, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      int
				mock       func(*customMock.CommonWarehouseRepository)
				assertMock func(*testing.T, *customMock.CommonWarehouseRepository)
			}{
				{
					name:  "Invalid id",
					input: 0,
					mock:  func(*customMock.CommonWarehouseRepository) {},
					assertMock: func(t *testing.T, mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.AssertNumberOfCalls(t, "FindByID", 0)
					},
				},
				{
					name:  "FindByID",
					input: idToFind,
					mock: func(mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.On("FindByID", idToFind).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockCWarehouseR *customMock.CommonWarehouseRepository) {
						mockCWarehouseR.AssertNumberOfCalls(t, "FindByID", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockWarehouseRepo := new(customMock.AdminWarehouseRepository)
					mockCWarehouseRepo := new(customMock.CommonWarehouseRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sWarehouse := NewWarehouseService(mockWarehouseRepo, mockCWarehouseRepo, mockStockMovementRepo)
					// Expectations
					tC.mock(mockCWarehouseRepo)
					// Action
					result, err := sWarehouse.FindByID(tC.input)
					// mock assertion
					mockWarehouseRepo.AssertExpectations(t)
					mockCWarehouseRepo.AssertExpectations(t)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(t, mockCWarehouseRepo)
					// Data assertion
					assert.Error(t, err)
					assert.Nil(t, result)
				})
			}
		})
	})
}
