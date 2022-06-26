package movement

import (
	"testing"
)

func TestMovementController(t *testing.T) {
	// Fixture
	/*serviceErr := goerrors.New("service error")
	t.Run("Create", func(t *testing.T) {
		movement := &movements[0]
		urlvalues := url.Values{}
		movementToCreate := &dto.NewMovement{
			Quantity:  movement.Quantity,
			Type:      movement.Type,
		}
		movementToCreateWithUserID := &dto.NewMovement{}
		copier.Copy(movementToCreateWithUserID, movementToCreate)
		movementToCreateWithUserID.UserID = 1
		expectedMovement := &entity.Movement{
			MovementID: movement.MovementID,
			UserID:     movement.UserID,
			Quantity:   movement.Quantity,
			Available:  movement.Available,
			Type:       movement.Type,
		}
		t.Run("Should success on", func(t *testing.T) {
			mockClientMovementService := new(mock.ClientMovementService)
			cMovement := NewMovementController(mockClientMovementService)

			// expectations
			mockClientMovementService.On("Create", movementToCreateWithUserID).Return(movement, nil)

			// action
			resp := mock.MHTTPHandle(http.MethodPost, "/", cMovement.Create, "", urlvalues, movementToCreate)

			//Mock Assertion
			mockClientMovementService.AssertExpectations(t)
			mockClientMovementService.AssertNumberOfCalls(t, "Create", 1)

			result := &entity.Movement{}
			bodyResponse, _ := utils.GetBodyResponse(resp, &result)

			//Data Assertion
			assert.Equal(t, 201, resp.StatusCode)
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "MOVEMENT.CREATED"}), bodyResponse.Message)
			assert.Empty(t, bodyResponse.Errors)
			assert.Equal(t, expectedMovement, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Missing body", func(t *testing.T) {
				mockClientMovementService := new(mock.ClientMovementService)
				cMovement := NewMovementController(mockClientMovementService)

				// action
				resp := mock.MHTTPHandle(http.MethodPost, "/", func(res http.ResponseWriter, req *http.Request) {
					ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, user)

					cMovement.Create(res, req.WithContext(ctx))
				}, "", urlvalues, nil)

				//Mock Assertion
				mockClientMovementService.AssertExpectations(t)
				mockClientMovementService.AssertNumberOfCalls(t, "Create", 0)

				result := &entity.Movement{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
			t.Run("Service error", func(t *testing.T) {
				mockClientMovementService := new(mock.ClientMovementService)
				cMovement := NewMovementController(mockClientMovementService)

				// expectations
				mockClientMovementService.On("Create", movementToCreateWithUserID).Return(nil, serviceErr)

				// action
				resp := mock.MHTTPHandle(http.MethodPost, "/", func(res http.ResponseWriter, req *http.Request) {
					ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, user)

					cMovement.Create(res, req.WithContext(ctx))
				}, "", urlvalues, movementToCreate)

				//Mock Assertion
				mockClientMovementService.AssertExpectations(t)
				mockClientMovementService.AssertNumberOfCalls(t, "Create", 1)

				result := &entity.Movement{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
		})
	})*/
}
