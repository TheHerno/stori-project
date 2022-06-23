package warehouse

import (
	goerrors "errors"
	"net/http"
	"net/url"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	myErrors "stori-service/src/libs/errors"
	"stori-service/src/libs/i18n"
	"stori-service/src/utils"
	"stori-service/src/utils/test/mock"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWarehouseController(t *testing.T) {
	// Fixture
	serviceErr := goerrors.New("service error")

	t.Run("Index", func(t *testing.T) {
		// Fixture
		pagination := dto.NewPagination(1, 20, 0)
		warehouseFixture := []entity.Warehouse{
			{
				WarehouseID: 11,
				Name:        "Warehouse once",
				UserID:      1,
				Address:     "Tucuman 123",
				CreatedAt:   time.Now(),
			},
			{
				WarehouseID: 22,
				Name:        "Warehouse veintidos",
				UserID:      11,
				Address:     "Tucuman 321",
				CreatedAt:   time.Now(),
			},
		}
		expectedWarehouses := []entity.Warehouse{
			{
				WarehouseID: 11,
				Name:        "Warehouse once",
				UserID:      1,
				Address:     "Tucuman 123",
			},
			{
				WarehouseID: 22,
				Name:        "Warehouse veintidos",
				UserID:      11,
				Address:     "Tucuman 321",
			},
		}

		t.Run("Should success on", func(t *testing.T) {
			testCases := []struct {
				TestName           string
				Querystring        url.Values
				ExpectedPagination *dto.Pagination
			}{
				{
					TestName:           "Without query string, using default values",
					ExpectedPagination: dto.NewPagination(1, 20, 0),
				},
				{
					TestName: "Valid pagination data",
					Querystring: url.Values{
						"page_size": []string{"10"},
						"page":      []string{"1"},
					},
					ExpectedPagination: dto.NewPagination(1, 10, 0),
				},
			}
			for _, tC := range testCases {
				t.Run(tC.TestName, func(t *testing.T) {
					mockAdminWarehouseService := new(mock.AdminWarehouseService)
					cWarehouse := NewWarehouseController(mockAdminWarehouseService)

					// Expectations
					mockAdminWarehouseService.On("Index", tC.ExpectedPagination).Return(&warehouseFixture, nil)

					// Actions
					resp := mock.MHTTPHandle(http.MethodGet, "/", cWarehouse.Index, "", tC.Querystring, nil)

					// Mock assertions
					mockAdminWarehouseService.AssertExpectations(t)
					mockAdminWarehouseService.AssertNumberOfCalls(t, "Index", 1)

					warehouses := []entity.Warehouse{}
					bodyResponse, _ := utils.GetBodyResponse(resp, &warehouses)

					// Data assertions
					assert.EqualValues(t, strconv.Itoa(tC.ExpectedPagination.PageSize), resp.Header.Get("X-pagination-page-size"))
					assert.EqualValues(t, strconv.Itoa(tC.ExpectedPagination.Page), resp.Header.Get("X-pagination-current-page"))
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					assert.Empty(t, bodyResponse.Errors)
					assert.Equal(t, expectedWarehouses, warehouses)
				})
			}
		})

		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				TestName      string
				Querystring   url.Values
				ExpectedError error
			}{
				{
					TestName: "Bad Page Size",
					Querystring: url.Values{
						"page_size": []string{"101"},
						"page":      []string{"1"},
					},
					ExpectedError: errors.ErrPageSizeTooHigh,
				},
				{
					TestName: "Bad Page",
					Querystring: url.Values{
						"page_size": []string{"20"},
						"page":      []string{"101"},
					},
					ExpectedError: errors.ErrPageTooHigh,
				},
			}
			for _, tC := range testCases {
				t.Run(tC.TestName, func(t *testing.T) {
					mockAdminWarehouseService := new(mock.AdminWarehouseService)
					cWarehouse := NewWarehouseController(mockAdminWarehouseService)

					// Actions
					resp := mock.MHTTPHandle(http.MethodGet, "/", cWarehouse.Index, "", tC.Querystring, nil)

					// Mock assertions
					mockAdminWarehouseService.AssertExpectations(t)
					mockAdminWarehouseService.AssertNumberOfCalls(t, "Index", 0)

					warehouses := []entity.Warehouse{}
					bodyResponse, _ := utils.GetBodyResponse(resp, &warehouses)

					// Data assertions
					assert.Equal(t, errors.GetStatusCode(tC.ExpectedError), resp.StatusCode)
					assert.NotEmpty(t, tC.ExpectedError.Error(), bodyResponse.Errors[0]["error"])
					assert.Empty(t, warehouses)
				})
			}
			t.Run("Index service fail", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				//Expectation
				mockAdminWarehouseService.On("Index", pagination).Return(nil, serviceErr)

				//Action
				resp := mock.MHTTPHandle(http.MethodGet, "/", cWarehouse.Index, "", url.Values{}, nil)

				//Mock Assertion
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "Index", 1)

				warehouses := []entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, &warehouses)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, warehouses)
			})
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		// Fixture
		path := `/{id}`
		warehouseId := 1212
		warehouseFound := &entity.Warehouse{
			WarehouseID: 1212,
			Name:        "Warehouse del centro",
			Address:     "San Juan 2020",
			CreatedAt:   time.Now(),
		}
		warehouseExpected := &entity.Warehouse{
			WarehouseID: 1212,
			Name:        "Warehouse del centro",
			Address:     "San Juan 2020",
		}
		t.Run("Should success on", func(t *testing.T) {
			mockAdminWarehouseService := new(mock.AdminWarehouseService)
			cWarehouse := NewWarehouseController(mockAdminWarehouseService)

			// Expectations
			mockAdminWarehouseService.On("FindByID", warehouseId).Return(warehouseFound, nil)

			// Action
			resp := mock.MHTTPHandle(http.MethodGet, path, cWarehouse.FindByID, strconv.Itoa(warehouseId), url.Values{}, nil)

			// Mock assertions
			mockAdminWarehouseService.AssertExpectations(t)
			mockAdminWarehouseService.AssertNumberOfCalls(t, "FindByID", 1)

			warehouse := &entity.Warehouse{}
			bodyResponse, _ := utils.GetBodyResponse(resp, warehouse)

			// Data assertions
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "WAREHOUSE.FOUND"}), bodyResponse.Message)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Empty(t, bodyResponse.Errors)
			assert.Equal(t, warehouseExpected, warehouse)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("FindByID service fails", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				// Expectations
				mockAdminWarehouseService.On("FindByID", warehouseId).Return(nil, serviceErr)

				// Action
				resp := mock.MHTTPHandle(http.MethodGet, path, cWarehouse.FindByID, strconv.Itoa(warehouseId), url.Values{}, nil)

				// Mock assertions
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "FindByID", 1)

				warehouse := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, warehouse)

				//Data Assertion
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Nil(t, bodyResponse.Data)
			})
			t.Run("Not numeric ID", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				// Action
				resp := mock.MHTTPHandle(http.MethodGet, `/{id}`, cWarehouse.FindByID, "abcd123", url.Values{}, nil)

				// Mock assertions
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "FindByID", 0)

				warehouse := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, warehouse)

				//Data Assertion
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				assert.Equal(t, myErrors.ErrIDNotNumeric.Error(), bodyResponse.Errors[0]["error"])
				assert.Nil(t, bodyResponse.Data)
			})
		})
	})

	t.Run("Update", func(t *testing.T) {
		// Fixture
		path := `/{id}`
		warehouse := warehouses[0]
		urlvalues := url.Values{}
		body := &dto.UpdateWarehouse{
			Name:    "New Name",
			Address: "New Address",
		}
		warehouseToUpdate := &dto.UpdateWarehouse{
			WarehouseID: warehouse.WarehouseID,
			Name:        "New Name",
			Address:     "New Address",
		}
		returnedWarehouse := &entity.Warehouse{
			WarehouseID: warehouse.WarehouseID,
			Name:        warehouseToUpdate.Name,
			Address:     warehouseToUpdate.Address,
			CreatedAt:   warehouse.CreatedAt,
			UpdatedAt:   warehouse.UpdatedAt,
		}
		expectedWarehouse := &entity.Warehouse{
			WarehouseID: warehouse.WarehouseID,
			Address:     warehouseToUpdate.Address,
			Name:        warehouseToUpdate.Name,
		}
		t.Run("Should success on", func(t *testing.T) {
			mockAdminWarehouseService := new(mock.AdminWarehouseService)
			cWarehouse := NewWarehouseController(mockAdminWarehouseService)

			//Expectation
			mockAdminWarehouseService.On("Update", warehouseToUpdate).Return(returnedWarehouse, nil)
			//Action
			resp := mock.MHTTPHandle(http.MethodPatch, path, cWarehouse.Update, strconv.Itoa(warehouse.WarehouseID), urlvalues, body)

			//Mock Assertion
			mockAdminWarehouseService.AssertExpectations(t)
			mockAdminWarehouseService.AssertNumberOfCalls(t, "Update", 1)

			result := &entity.Warehouse{}
			bodyResponse, _ := utils.GetBodyResponse(resp, &result)

			//Data Assertion
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "WAREHOUSE.UPDATED"}), bodyResponse.Message)
			assert.Empty(t, bodyResponse.Errors)
			assert.Equal(t, expectedWarehouse, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Update service fails", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				// Expectations
				mockAdminWarehouseService.On("Update", warehouseToUpdate).Return(nil, serviceErr)

				// Action
				resp := mock.MHTTPHandle(http.MethodPatch, path, cWarehouse.Update, strconv.Itoa(warehouse.WarehouseID), url.Values{}, warehouseToUpdate)

				// Mock assertions
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "Update", 1)

				warehouse := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, warehouse)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Nil(t, bodyResponse.Data)
			})
			t.Run("Not numeric ID", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				// Action
				resp := mock.MHTTPHandle(http.MethodPatch, `/{id}`, cWarehouse.Update, "abcd123", url.Values{}, warehouseToUpdate)

				// Mock assertions
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "Update", 0)

				warehouse := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, warehouse)

				//Data Assertion
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				assert.Equal(t, myErrors.ErrIDNotNumeric.Error(), bodyResponse.Errors[0]["error"])
				assert.Nil(t, bodyResponse.Data)
			})
			t.Run("Missing Body", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				// Action
				resp := mock.MHTTPHandle(http.MethodPatch, path, cWarehouse.Update, strconv.Itoa(warehouse.WarehouseID), urlvalues, nil)

				// Mock assertions
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "Update", 0)

				result := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
		})
	})
	t.Run("Create", func(t *testing.T) {
		// Fixture
		warehouse := &warehouses[0]
		urlvalues := url.Values{}
		warehouseToCreate := &dto.CreateWarehouse{
			Name:    warehouse.Name,
			Address: warehouse.Address,
			UserID:  2,
		}
		expectedWarehouse := &entity.Warehouse{
			WarehouseID: warehouse.WarehouseID,
			Name:        warehouse.Name,
			UserID:      2,
			Address:     warehouse.Address,
		}
		t.Run("Should success on", func(t *testing.T) {
			mockAdminWarehouseService := new(mock.AdminWarehouseService)
			cWarehouse := NewWarehouseController(mockAdminWarehouseService)
			//Expectation
			mockAdminWarehouseService.On("Create", warehouseToCreate).Return(warehouse, nil)
			//Action
			resp := mock.MHTTPHandle(http.MethodPost, "/", cWarehouse.Create, "", urlvalues, warehouseToCreate)

			//Mock Assertion
			mockAdminWarehouseService.AssertExpectations(t)
			mockAdminWarehouseService.AssertNumberOfCalls(t, "Create", 1)

			result := &entity.Warehouse{}
			bodyResponse, _ := utils.GetBodyResponse(resp, &result)

			//Data Assertion
			assert.Equal(t, 201, resp.StatusCode)
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "WAREHOUSE.CREATED"}), bodyResponse.Message)
			assert.Empty(t, bodyResponse.Errors)
			assert.Equal(t, expectedWarehouse, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Missing body", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				//Action
				resp := mock.MHTTPHandle(http.MethodPost, "/", cWarehouse.Create, "", urlvalues, nil)
				//Mock Assertion
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "Create", 0)

				result := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
			t.Run("Service error", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				//Expectation
				mockAdminWarehouseService.On("Create", warehouseToCreate).Return(nil, serviceErr)

				//Action
				resp := mock.MHTTPHandle(http.MethodPost, "/", cWarehouse.Create, "", urlvalues, warehouseToCreate)

				//Mock Assertion
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "Create", 1)

				result := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
		})
	})
	t.Run("Delete", func(t *testing.T) {
		// Fixture
		warehouse := &warehouses[0]
		urlvalues := url.Values{}
		path := `/{id}`
		t.Run("Should success on", func(t *testing.T) {
			mockAdminWarehouseService := new(mock.AdminWarehouseService)
			cWarehouse := NewWarehouseController(mockAdminWarehouseService)

			//Expectation
			mockAdminWarehouseService.On("Delete", warehouse.WarehouseID).Return(nil)

			//Action
			resp := mock.MHTTPHandle(http.MethodDelete, path, cWarehouse.Delete, strconv.Itoa(warehouse.WarehouseID), urlvalues, nil)

			//Mock Assertion
			mockAdminWarehouseService.AssertExpectations(t)
			mockAdminWarehouseService.AssertNumberOfCalls(t, "Delete", 1)

			result := &entity.Warehouse{}
			bodyResponse, _ := utils.GetBodyResponse(resp, result)

			//Data Assertion
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "WAREHOUSE.DELETED"}), bodyResponse.Message)
			assert.Empty(t, bodyResponse.Errors)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Service error", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				//Expectation
				mockAdminWarehouseService.On("Delete", warehouse.WarehouseID).Return(serviceErr)

				//Action
				resp := mock.MHTTPHandle(http.MethodDelete, path, cWarehouse.Delete, strconv.Itoa(warehouse.WarehouseID), urlvalues, nil)

				//Mock Assertion
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "Delete", 1)

				result := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
			t.Run("Not numeric ID", func(t *testing.T) {
				mockAdminWarehouseService := new(mock.AdminWarehouseService)
				cWarehouse := NewWarehouseController(mockAdminWarehouseService)

				//Action
				resp := mock.MHTTPHandle(http.MethodDelete, `/{id}`, cWarehouse.Delete, "abcd123", urlvalues, nil)

				//Mock Assertion
				mockAdminWarehouseService.AssertExpectations(t)
				mockAdminWarehouseService.AssertNumberOfCalls(t, "Delete", 0)

				result := &entity.Warehouse{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				assert.Equal(t, myErrors.ErrIDNotNumeric.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
		})
	})
}
