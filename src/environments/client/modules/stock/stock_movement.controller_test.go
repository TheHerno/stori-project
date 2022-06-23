package stock

import (
	"context"
	goerrors "errors"
	"net/http"
	"net/url"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/libs/i18n"
	"stori-service/src/libs/middleware"
	"stori-service/src/utils"
	"stori-service/src/utils/test/mock"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func TestStockMovementController(t *testing.T) {
	// Fixture
	serviceErr := goerrors.New("service error")
	t.Run("Create", func(t *testing.T) {
		stockMovement := &stockMovements[0]
		urlvalues := url.Values{}
		user := &middleware.User{
			UserID: 1,
			Email:  "messi@gmail.com",
			Name:   "Lionel Messi",
			Role:   "client",
		}
		stockMovementToCreate := &dto.NewStockMovement{
			ProductID: stockMovement.ProductID,
			Quantity:  stockMovement.Quantity,
			Concept:   stockMovement.Concept,
			Type:      stockMovement.Type,
		}
		stockMovementToCreateWithUserID := &dto.NewStockMovement{}
		copier.Copy(stockMovementToCreateWithUserID, stockMovementToCreate)
		stockMovementToCreateWithUserID.UserID = 1
		expectedStockMovement := &entity.StockMovement{
			StockMovementID: stockMovement.StockMovementID,
			ProductID:       stockMovement.ProductID,
			WarehouseID:     stockMovement.WarehouseID,
			Quantity:        stockMovement.Quantity,
			Available:       stockMovement.Available,
			Concept:         stockMovement.Concept,
			Type:            stockMovement.Type,
		}
		t.Run("Should success on", func(t *testing.T) {
			mockClientStockMovementService := new(mock.ClientStockMovementService)
			cStockMovement := NewStockMovementController(mockClientStockMovementService)

			// expectations
			mockClientStockMovementService.On("Create", stockMovementToCreateWithUserID).Return(stockMovement, nil)

			// action
			resp := mock.MHTTPHandle(http.MethodPost, "/", func(res http.ResponseWriter, req *http.Request) {
				ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, user)

				cStockMovement.Create(res, req.WithContext(ctx))
			}, "", urlvalues, stockMovementToCreate)

			//Mock Assertion
			mockClientStockMovementService.AssertExpectations(t)
			mockClientStockMovementService.AssertNumberOfCalls(t, "Create", 1)

			result := &entity.StockMovement{}
			bodyResponse, _ := utils.GetBodyResponse(resp, &result)

			//Data Assertion
			assert.Equal(t, 201, resp.StatusCode)
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "STOCKMOVEMENT.CREATED"}), bodyResponse.Message)
			assert.Empty(t, bodyResponse.Errors)
			assert.Equal(t, expectedStockMovement, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Missing body", func(t *testing.T) {
				mockClientStockMovementService := new(mock.ClientStockMovementService)
				cStockMovement := NewStockMovementController(mockClientStockMovementService)

				// action
				resp := mock.MHTTPHandle(http.MethodPost, "/", func(res http.ResponseWriter, req *http.Request) {
					ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, user)

					cStockMovement.Create(res, req.WithContext(ctx))
				}, "", urlvalues, nil)

				//Mock Assertion
				mockClientStockMovementService.AssertExpectations(t)
				mockClientStockMovementService.AssertNumberOfCalls(t, "Create", 0)

				result := &entity.StockMovement{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
			t.Run("Service error", func(t *testing.T) {
				mockClientStockMovementService := new(mock.ClientStockMovementService)
				cStockMovement := NewStockMovementController(mockClientStockMovementService)

				// expectations
				mockClientStockMovementService.On("Create", stockMovementToCreateWithUserID).Return(nil, serviceErr)

				// action
				resp := mock.MHTTPHandle(http.MethodPost, "/", func(res http.ResponseWriter, req *http.Request) {
					ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, user)

					cStockMovement.Create(res, req.WithContext(ctx))
				}, "", urlvalues, stockMovementToCreate)

				//Mock Assertion
				mockClientStockMovementService.AssertExpectations(t)
				mockClientStockMovementService.AssertNumberOfCalls(t, "Create", 1)

				result := &entity.StockMovement{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
		})
	})
}
