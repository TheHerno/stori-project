package stock

import (
	goerrors "errors"
	"stori-service/src/environments/client/modules/stock/providers/params"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/utils/constant"
	customMocks "stori-service/src/utils/test/mock"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStockMovementService(t *testing.T) {
	// Fixture
	repositoryErr := goerrors.New("repository error")
	paramServiceErr := goerrors.New("param service error")
	commitErr := goerrors.New("Commit Error")
	t.Run("Create", func(t *testing.T) {
		// fixture
		stockMovementToCreate := &dto.NewStockMovement{
			ProductID: stockMovements[0].ProductID,
			UserID:    warehouses[0].UserID,
			Quantity:  stockMovements[0].Quantity,
			Concept:   stockMovements[0].Concept,
			Type:      stockMovements[0].Type,
		}
		stockMovementTransferToCreate := &dto.NewStockMovement{
			ProductID:   stockMovements[1].ProductID,
			UserID:      warehouses[0].UserID,
			Quantity:    stockMovements[1].Quantity,
			Concept:     stockMovements[1].Concept,
			Type:        stockMovements[1].Type,
			WarehouseID: warehouses[1].WarehouseID,
		}
		createdStockMovement := &entity.StockMovement{
			ProductID:   stockMovementToCreate.ProductID,
			WarehouseID: warehouses[0].WarehouseID,
			Quantity:    stockMovementToCreate.Quantity,
			Available:   stockMovementToCreate.Quantity + stockMovements[0].Available,
			Type:        stockMovementToCreate.Type,
			Concept:     stockMovementToCreate.Concept,
		}
		createdReceiverStockMovement := &entity.StockMovement{
			ProductID:   stockMovementTransferToCreate.ProductID,
			WarehouseID: warehouses[1].WarehouseID,
			Quantity:    stockMovementTransferToCreate.Quantity,
			Available:   stockMovementTransferToCreate.Quantity,
			Type:        constant.IncomeType,
			Concept:     "Transfer from " + warehouses[0].Name,
		}
		createdTransferStockMovement := &entity.StockMovement{
			ProductID:   stockMovementTransferToCreate.ProductID,
			WarehouseID: warehouses[0].WarehouseID,
			Quantity:    stockMovementTransferToCreate.Quantity,
			Available:   stockMovementTransferToCreate.Quantity*stockMovementTransferToCreate.Type + stockMovements[0].Available,
			Type:        stockMovementTransferToCreate.Type,
			Concept:     stockMovementTransferToCreate.Concept,
		}
		createdStockMovementWithoutPrevious := &entity.StockMovement{
			ProductID:   stockMovementToCreate.ProductID,
			WarehouseID: warehouses[0].WarehouseID,
			Quantity:    stockMovementToCreate.Quantity,
			Available:   stockMovementToCreate.Quantity,
			Type:        stockMovementToCreate.Type,
			Concept:     stockMovementToCreate.Concept,
		}
		t.Run("Should success on", func(t *testing.T) {
			getParamInt = func(param string) (int, error) {
				return 1000, nil
			}
			testCases := []struct {
				name     string
				mock     func(*customMocks.ClientStockMovementRepository, *customMocks.ClientWarehouseRepository, *customMocks.CommonWarehouseRepository, *customMocks.CommonProductRepository)
				expected *entity.StockMovement
			}{
				{
					name: "With existent previous movement",
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Commit").Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementToCreate.ProductID).Return(&products[0], nil)
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementToCreate.ProductID).Return(&stockMovements[0], nil)
						mockClientStockmovementRepo.On("Create", createdStockMovement).Return(createdStockMovement, nil)
					},
					expected: createdStockMovement,
				},
				{
					name: "Without existent previous movement",
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Commit").Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementToCreate.ProductID).Return(&products[0], nil)
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementToCreate.ProductID).Return(nil, errors.ErrNotFound)
						mockClientStockmovementRepo.On("Create", createdStockMovementWithoutPrevious).Return(createdStockMovementWithoutPrevious, nil)
					},
					expected: createdStockMovementWithoutPrevious,
				},
			}

			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockClientStockmovementRepo := new(customMocks.ClientStockMovementRepository)
					mockClientWarehouseRepo := new(customMocks.ClientWarehouseRepository)
					mockCommonProductRepo := new(customMocks.CommonProductRepository)
					mockCommonWarehouseRepo := new(customMocks.CommonWarehouseRepository)
					tC.mock(mockClientStockmovementRepo, mockClientWarehouseRepo, mockCommonWarehouseRepo, mockCommonProductRepo)
					sStockMovement := NewStockMovementService(mockClientStockmovementRepo, mockClientWarehouseRepo, mockCommonWarehouseRepo, mockCommonProductRepo)

					result, err := sStockMovement.Create(stockMovementToCreate)

					mockClientStockmovementRepo.AssertExpectations(t)
					mockClientWarehouseRepo.AssertExpectations(t)
					mockCommonWarehouseRepo.AssertExpectations(t)
					mockCommonProductRepo.AssertExpectations(t)
					mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 1)
					mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
					mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
					mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 1)
					mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
					mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
					mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
					mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
					mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
					mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
					mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
					mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
					mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
					mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 1)

					// data assertion
					assert.Nil(t, err)
					assert.Equal(t, tC.expected, result)
				})
			}
			t.Run("Transfer to another Warehouse", func(t *testing.T) {

				mockClientStockmovementRepo := new(customMocks.ClientStockMovementRepository)
				mockClientWarehouseRepo := new(customMocks.ClientWarehouseRepository)
				mockCommonProductRepo := new(customMocks.CommonProductRepository)
				mockCommonWarehouseRepo := new(customMocks.CommonWarehouseRepository)

				mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
				mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
				mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
				mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
				mockClientStockmovementRepo.On("Begin", nil).Return(nil)
				mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
				mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
				mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
				mockClientStockmovementRepo.On("Commit").Return(nil)
				mockClientStockmovementRepo.On("Rollback").Return(nil)
				mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementTransferToCreate.UserID).Return(&warehouses[0], nil)
				mockCommonWarehouseRepo.On("FindAndLockByID", warehouses[1].WarehouseID).Return(&warehouses[1], nil)
				mockCommonProductRepo.On("FindAndLockByID", stockMovementTransferToCreate.ProductID).Return(&products[0], nil)
				mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementTransferToCreate.ProductID).Return(&stockMovements[0], nil).Once()
				mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[1].WarehouseID, stockMovementTransferToCreate.ProductID).Return(nil, errors.ErrNotFound).Once()
				mockClientStockmovementRepo.On("Create", createdTransferStockMovement).Return(createdTransferStockMovement, nil).Once()
				mockClientStockmovementRepo.On("Create", createdReceiverStockMovement).Return(createdReceiverStockMovement, nil).Once()

				sStockMovement := NewStockMovementService(mockClientStockmovementRepo, mockClientWarehouseRepo, mockCommonWarehouseRepo, mockCommonProductRepo)

				result, err := sStockMovement.Create(stockMovementTransferToCreate)

				mockClientStockmovementRepo.AssertExpectations(t)
				mockClientWarehouseRepo.AssertExpectations(t)
				mockCommonWarehouseRepo.AssertExpectations(t)
				mockCommonProductRepo.AssertExpectations(t)
				mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 2)
				mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
				mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
				mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 2)
				mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
				mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
				mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 1)

				// data assertion
				assert.Nil(t, err)
				assert.Equal(t, createdTransferStockMovement, result)
			})
			t.Cleanup(func() {
				getParamInt = params.GetParamInt
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      *dto.NewStockMovement
				mock       func(*customMocks.ClientStockMovementRepository, *customMocks.ClientWarehouseRepository, *customMocks.CommonWarehouseRepository, *customMocks.CommonProductRepository)
				assertMock func(*testing.T, *customMocks.ClientStockMovementRepository, *customMocks.ClientWarehouseRepository, *customMocks.CommonWarehouseRepository, *customMocks.CommonProductRepository)
			}{
				{
					name: "Invalid input",
					input: &dto.NewStockMovement{
						ProductID: 0,
						UserID:    0,
						Quantity:  0,
						Concept:   strings.Repeat("a", 101),
						Type:      0,
					},
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name:  "Fail finding warehouse",
					input: stockMovementToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name:  "Fail finding product",
					input: stockMovementToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementToCreate.ProductID).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name:  "Product not enabled",
					input: stockMovementToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementToCreate.ProductID).Return(&products[3], nil)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name:  "Fail finding last stock movement",
					input: stockMovementToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementToCreate.ProductID).Return(&products[0], nil)
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementToCreate.ProductID).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name: "Quantity to extract is more than stock",
					input: &dto.NewStockMovement{
						ProductID: stockMovementToCreate.ProductID,
						UserID:    stockMovementToCreate.UserID,
						Quantity:  100,
						Type:      constant.OutcomeType,
						Concept:   stockMovementToCreate.Concept,
					},
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementToCreate.ProductID).Return(&products[0], nil)
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementToCreate.ProductID).Return(&stockMovements[0], nil)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name:  "Create",
					input: stockMovementToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementToCreate.ProductID).Return(&products[0], nil)
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementToCreate.ProductID).Return(&stockMovements[0], nil)
						newStockMovement := &entity.StockMovement{
							ProductID:   stockMovementToCreate.ProductID,
							WarehouseID: warehouses[0].WarehouseID,
							Quantity:    stockMovementToCreate.Quantity,
							Available:   stockMovementToCreate.Quantity + stockMovements[0].Available,
							Type:        stockMovementToCreate.Type,
							Concept:     stockMovementToCreate.Concept,
						}
						mockClientStockmovementRepo.On("Create", newStockMovement).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name: "Invalid type to transfer",
					input: &dto.NewStockMovement{
						ProductID:   stockMovementTransferToCreate.ProductID,
						UserID:      stockMovementTransferToCreate.UserID,
						WarehouseID: warehouses[1].WarehouseID,
						Quantity:    stockMovementTransferToCreate.Quantity,
						Concept:     stockMovementTransferToCreate.Concept,
						Type:        constant.IncomeType,
					},
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name: "Same warehouse to transfer",
					input: &dto.NewStockMovement{
						ProductID:   stockMovementTransferToCreate.ProductID,
						UserID:      stockMovementTransferToCreate.UserID,
						WarehouseID: warehouses[0].WarehouseID,
						Quantity:    stockMovementTransferToCreate.Quantity,
						Concept:     stockMovementTransferToCreate.Concept,
						Type:        constant.OutcomeType,
					},
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
					},
				},
				{
					name:  "Fail getting receiver warehouse",
					input: stockMovementTransferToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementTransferToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonWarehouseRepo.On("FindAndLockByID", warehouses[1].WarehouseID).Return(nil, repositoryErr)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementTransferToCreate.ProductID).Return(&products[0], nil)
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementTransferToCreate.ProductID).Return(&stockMovements[0], nil)
						mockClientStockmovementRepo.On("Create", createdTransferStockMovement).Return(createdTransferStockMovement, nil)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
					},
				},
				{
					name:  "Fail creating receiver movement",
					input: stockMovementTransferToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementTransferToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonWarehouseRepo.On("FindAndLockByID", warehouses[1].WarehouseID).Return(&warehouses[1], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementTransferToCreate.ProductID).Return(&products[0], nil)
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementTransferToCreate.ProductID).Return(&stockMovements[0], nil).Once()
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[1].WarehouseID, stockMovementTransferToCreate.ProductID).Return(nil, errors.ErrNotFound).Once()
						mockClientStockmovementRepo.On("Create", createdReceiverStockMovement).Return(nil, repositoryErr).Once()
						mockClientStockmovementRepo.On("Create", createdTransferStockMovement).Return(createdTransferStockMovement, nil).Once()
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 2)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 2)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
					},
				},
				{
					name:  "Commit",
					input: stockMovementToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.On("Clone").Return(mockClientStockmovementRepo)
						mockClientWarehouseRepo.On("Clone").Return(mockClientWarehouseRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockCommonWarehouseRepo.On("Clone").Return(mockCommonWarehouseRepo)
						mockClientStockmovementRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockClientWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockCommonWarehouseRepo.On("Begin", mock.Anything).Return(nil)
						mockClientStockmovementRepo.On("Rollback").Return(nil)
						mockClientWarehouseRepo.On("FindAndLockByUserID", stockMovementToCreate.UserID).Return(&warehouses[0], nil)
						mockCommonProductRepo.On("FindAndLockByID", stockMovementToCreate.ProductID).Return(&products[0], nil)
						mockClientStockmovementRepo.On("FindLastStockMovement", warehouses[0].WarehouseID, stockMovementToCreate.ProductID).Return(&stockMovements[0], nil)
						newStockMovement := &entity.StockMovement{
							ProductID:   stockMovementToCreate.ProductID,
							WarehouseID: warehouses[0].WarehouseID,
							Quantity:    stockMovementToCreate.Quantity,
							Available:   stockMovementToCreate.Quantity + stockMovements[0].Available,
							Type:        stockMovementToCreate.Type,
							Concept:     stockMovementToCreate.Concept,
						}
						mockClientStockmovementRepo.On("Create", newStockMovement).Return(createdStockMovement, nil)
						mockClientStockmovementRepo.On("Commit").Return(commitErr)
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 1)
					},
				},
				{
					name:  "Error getting max from params",
					input: stockMovementToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						getParamInt = func(param string) (int, error) {
							return 0, paramServiceErr
						}
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
					},
				},
				{
					name:  "Quantity greater than max allowed",
					input: stockMovementToCreate,
					mock: func(mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						getParamInt = func(param string) (int, error) {
							return 1, nil
						}
					},
					assertMock: func(t *testing.T, mockClientStockmovementRepo *customMocks.ClientStockMovementRepository, mockClientWarehouseRepo *customMocks.ClientWarehouseRepository, mockCommonWarehouseRepo *customMocks.CommonWarehouseRepository, mockCommonProductRepo *customMocks.CommonProductRepository) {
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Create", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "FindAndLockByUserID", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "FindLastStockMovement", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockClientWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCommonWarehouseRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Rollback", 0)
						mockClientStockmovementRepo.AssertNumberOfCalls(t, "Commit", 0)
					},
				},
			}

			for _, testCase := range testCases {
				t.Run(testCase.name, func(t *testing.T) {
					mockClientStockmovementRepo := new(customMocks.ClientStockMovementRepository)
					mockClientWarehouseRepo := new(customMocks.ClientWarehouseRepository)
					mockCommonProductRepo := new(customMocks.CommonProductRepository)
					mockCommonWarehouseRepo := new(customMocks.CommonWarehouseRepository)
					sStockMovement := NewStockMovementService(mockClientStockmovementRepo, mockClientWarehouseRepo, mockCommonWarehouseRepo, mockCommonProductRepo)
					// inject spy
					getParamInt = func(param string) (int, error) {
						return 1000, nil
					}
					testCase.mock(mockClientStockmovementRepo, mockClientWarehouseRepo, mockCommonWarehouseRepo, mockCommonProductRepo)

					result, err := sStockMovement.Create(testCase.input)

					mockClientStockmovementRepo.AssertExpectations(t)
					mockClientWarehouseRepo.AssertExpectations(t)
					mockCommonWarehouseRepo.AssertExpectations(t)
					mockCommonProductRepo.AssertExpectations(t)
					testCase.assertMock(t, mockClientStockmovementRepo, mockClientWarehouseRepo, mockCommonWarehouseRepo, mockCommonProductRepo)
					assert.Error(t, err)
					assert.Nil(t, result)
					t.Cleanup(func() {
						getParamInt = params.GetParamInt
					})
				})
			}
		})
	})
}
