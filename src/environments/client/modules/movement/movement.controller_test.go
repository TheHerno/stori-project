package movement

import (
	goErrors "errors"
	"net/http"
	"net/url"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/libs/i18n"
	"stori-service/src/utils"
	"stori-service/src/utils/constant"
	"stori-service/src/utils/test/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMovementController(t *testing.T) {
	urlvalues := url.Values{}
	serviceErr := goErrors.New("service error")
	path := `/{id}`
	expectedMovementList := &dto.MovementList{
		Customer: &entity.Customer{
			CustomerID: 1,
			Name:       "pepe",
			Email:      "pepepe@hotmail.com",
		},
		Movements: []entity.Movement{
			{
				MovementID: 1,
				CustomerID: 1,
				Available:  100.00,
				Quantity:   100.00,
				Type:       constant.IncomeType,
				Date:       time.Now(),
			},
		},
	}
	t.Run("ProcessFile", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Processing file", func(t *testing.T) {
				// fixture
				mockMovementService := new(mock.ClientMovementService)
				movementControler := NewMovementController(mockMovementService)

				// mock expectations
				mockMovementService.On("ProcessFile", 1).Return(expectedMovementList, nil)

				//Action
				resp := mock.MHTTPHandle(http.MethodGet, path, movementControler.ProcessFile, "1", urlvalues, nil)

				//Mock Assertion
				mockMovementService.AssertExpectations(t)
				mockMovementService.AssertNumberOfCalls(t, "ProcessFile", 1)

				result := &dto.MovementList{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, i18n.T(i18n.Message{MessageID: "MOVEMENT_LIST.CREATED"}), bodyResponse.Message)
				assert.Empty(t, bodyResponse.Errors)
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Invalid id", func(t *testing.T) {
				// fixture
				movementControler := NewMovementController(nil)

				//Action
				resp := mock.MHTTPHandle(http.MethodGet, path, movementControler.ProcessFile, "asd", urlvalues, nil)

				result := &dto.MovementList{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
			t.Run("Service fails processing file", func(t *testing.T) {
				// fixture
				mockMovementService := new(mock.ClientMovementService)
				movementControler := NewMovementController(mockMovementService)

				// mock expectations
				mockMovementService.On("ProcessFile", 1).Return(nil, serviceErr)

				//Action
				resp := mock.MHTTPHandle(http.MethodGet, path, movementControler.ProcessFile, "1", urlvalues, nil)

				//Mock Assertion
				mockMovementService.AssertExpectations(t)
				mockMovementService.AssertNumberOfCalls(t, "ProcessFile", 1)

				result := &dto.MovementList{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
		})
	})
}
