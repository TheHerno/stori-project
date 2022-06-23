package product

import (
	goerrors "errors"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/utils/helpers"
	"stori-service/src/utils/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductService(t *testing.T) {
	// Fixture
	repositoryErr := goerrors.New("repository error")
	t.Run("GetStockList", func(t *testing.T) {
		// fixture
		pagination := dto.NewPagination(1, 20, 0)
		userIDToFind := 1
		warehouseFound := &entity.Warehouse{
			WarehouseID: 1,
			Name:        "Warehouse 1",
			Address:     "Address 1",
			UserID:      1,
		}
		expectedStocks := []dto.ProductWithStock{
			{
				ProductID:   1,
				Name:        "Auto",
				Description: helpers.PointerToString("Ferrari"),
				Slug:        "auto",
				Stock:       15,
			},
			{
				ProductID:   2,
				Name:        "Motocicleta",
				Description: helpers.PointerToString("Zanelita 50"),
				Slug:        "motocicleta",
				Stock:       9,
			},
		}
		t.Run("Should success on", func(t *testing.T) {
			mockWarehouseRepo := new(mock.ClientWarehouseRepository)
			mockStockMovementRepo := new(mock.ClientStockMovementRepository)
			sProduct := NewProductService(mockWarehouseRepo, mockStockMovementRepo)

			// expectations
			mockWarehouseRepo.On("FindByUserID", userIDToFind).Return(warehouseFound, nil)
			mockStockMovementRepo.On("FindStocksByWarehouse", warehouseFound.WarehouseID, pagination).Return(expectedStocks, nil)

			// action
			stocks, err := sProduct.GetStockList(userIDToFind, pagination)

			// mock assertions
			mockWarehouseRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockWarehouseRepo.AssertNumberOfCalls(t, "FindByUserID", 1)
			mockStockMovementRepo.AssertNumberOfCalls(t, "FindStocksByWarehouse", 1)

			// data assertions
			assert.NoError(t, err)
			assert.Equal(t, expectedStocks, stocks)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("FindByUserID", func(t *testing.T) {
				mockWarehouseRepo := new(mock.ClientWarehouseRepository)
				mockStockMovementRepo := new(mock.ClientStockMovementRepository)
				sProduct := NewProductService(mockWarehouseRepo, mockStockMovementRepo)

				// expectations
				mockWarehouseRepo.On("FindByUserID", userIDToFind).Return(nil, repositoryErr)

				// action
				stocks, err := sProduct.GetStockList(userIDToFind, pagination)

				// mock assertions
				mockWarehouseRepo.AssertExpectations(t)
				mockStockMovementRepo.AssertExpectations(t)
				mockWarehouseRepo.AssertNumberOfCalls(t, "FindByUserID", 1)
				mockStockMovementRepo.AssertNumberOfCalls(t, "FindStocksByWarehouse", 0)

				// data assertions
				assert.Error(t, err)
				assert.Nil(t, stocks)
			})
			t.Run("FindStocksByWarehouse", func(t *testing.T) {
				mockWarehouseRepo := new(mock.ClientWarehouseRepository)
				mockStockMovementRepo := new(mock.ClientStockMovementRepository)
				sProduct := NewProductService(mockWarehouseRepo, mockStockMovementRepo)

				// expectations
				mockWarehouseRepo.On("FindByUserID", userIDToFind).Return(warehouseFound, nil)
				mockStockMovementRepo.On("FindStocksByWarehouse", warehouseFound.WarehouseID, pagination).Return(nil, repositoryErr)

				// action
				stocks, err := sProduct.GetStockList(userIDToFind, pagination)

				// mock assertions
				mockWarehouseRepo.AssertExpectations(t)
				mockStockMovementRepo.AssertExpectations(t)
				mockWarehouseRepo.AssertNumberOfCalls(t, "FindByUserID", 1)
				mockStockMovementRepo.AssertNumberOfCalls(t, "FindStocksByWarehouse", 1)

				// data assertions
				assert.Error(t, err)
				assert.Nil(t, stocks)
			})
		})
	})
}
