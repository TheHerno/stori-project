package product

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
	"stori-service/src/utils/helpers"
	"stori-service/src/utils/test/mock"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProductController(t *testing.T) {
	// Fixture
	serviceErr := goerrors.New("service error")

	t.Run("Index", func(t *testing.T) {
		// Fixture
		pagination := dto.NewPagination(1, 20, 0)
		desc1 := helpers.PointerToString("Productito")
		desc2 := helpers.PointerToString("Shampoo")
		productFixture := []entity.Product{
			{
				ProductID:   11,
				Name:        "Product once",
				Slug:        "product-once",
				Description: desc1,
				Enabled:     &trueValue,
				CreatedAt:   time.Now(),
			},
			{
				ProductID:   22,
				Name:        "Product veintidos",
				Slug:        "product-veintidos",
				Description: desc2,
				Enabled:     &falseValue,
				CreatedAt:   time.Now(),
			},
		}
		expectedProducts := []entity.Product{
			{
				ProductID:   11,
				Name:        "Product once",
				Slug:        "product-once",
				Description: desc1,
				Enabled:     &trueValue,
			},
			{
				ProductID:   22,
				Name:        "Product veintidos",
				Slug:        "product-veintidos",
				Description: desc2,
				Enabled:     &falseValue,
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
					mockAdminProductService := new(mock.AdminProductService)
					cProduct := NewProductController(mockAdminProductService)

					// Expectations
					mockAdminProductService.On("Index", tC.ExpectedPagination).Return(productFixture, nil)

					// Actions
					resp := mock.MHTTPHandle(http.MethodGet, "/", cProduct.Index, "", tC.Querystring, nil)

					// Mock assertions
					mockAdminProductService.AssertExpectations(t)
					mockAdminProductService.AssertNumberOfCalls(t, "Index", 1)

					products := []entity.Product{}
					bodyResponse, _ := utils.GetBodyResponse(resp, &products)

					// Data assertions
					assert.EqualValues(t, strconv.Itoa(tC.ExpectedPagination.PageSize), resp.Header.Get("X-pagination-page-size"))
					assert.EqualValues(t, strconv.Itoa(tC.ExpectedPagination.Page), resp.Header.Get("X-pagination-current-page"))
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					assert.Empty(t, bodyResponse.Errors)
					assert.Equal(t, expectedProducts, products)
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
					mockAdminProductService := new(mock.AdminProductService)
					cProduct := NewProductController(mockAdminProductService)

					// Actions
					resp := mock.MHTTPHandle(http.MethodGet, "/", cProduct.Index, "", tC.Querystring, nil)

					// Mock assertions
					mockAdminProductService.AssertExpectations(t)
					mockAdminProductService.AssertNumberOfCalls(t, "Index", 0)

					products := []entity.Product{}
					bodyResponse, _ := utils.GetBodyResponse(resp, &products)

					// Data assertions
					assert.Equal(t, errors.GetStatusCode(tC.ExpectedError), resp.StatusCode)
					assert.NotEmpty(t, tC.ExpectedError.Error(), bodyResponse.Errors[0]["error"])
					assert.Empty(t, products)
				})
			}
			t.Run("Index service fail", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				//Expectation
				mockAdminProductService.On("Index", pagination).Return(nil, serviceErr)

				//Action
				resp := mock.MHTTPHandle(http.MethodGet, "/", cProduct.Index, "", url.Values{}, nil)

				//Mock Assertion
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "Index", 1)

				products := []entity.Product{}
				bodyResponse, _ := utils.GetBodyResponse(resp, &products)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, products)
			})
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		// Fixture
		path := `/{id}`
		productId := 1212
		desc := helpers.PointerToString("Motherboard")
		productFound := &entity.Product{
			ProductID:   1212,
			Name:        "Gigabyte A320m",
			Slug:        "gigabyte-a320m",
			Description: desc,
			Enabled:     &trueValue,
			CreatedAt:   time.Now(),
		}
		productExpected := &entity.Product{
			ProductID:   1212,
			Name:        "Gigabyte A320m",
			Slug:        "gigabyte-a320m",
			Description: desc,
			Enabled:     &trueValue,
		}
		t.Run("Should success on", func(t *testing.T) {
			mockAdminProductService := new(mock.AdminProductService)
			cProduct := NewProductController(mockAdminProductService)

			// Expectations
			mockAdminProductService.On("FindByID", productId).Return(productFound, nil)

			// Action
			resp := mock.MHTTPHandle(http.MethodGet, path, cProduct.FindByID, strconv.Itoa(productId), url.Values{}, nil)

			// Mock assertions
			mockAdminProductService.AssertExpectations(t)
			mockAdminProductService.AssertNumberOfCalls(t, "FindByID", 1)

			product := &entity.Product{}
			bodyResponse, _ := utils.GetBodyResponse(resp, product)

			// Data assertions
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "PRODUCT.FOUND"}), bodyResponse.Message)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Empty(t, bodyResponse.Errors)
			assert.Equal(t, productExpected, product)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("FindByID service fails", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				// Expectations
				mockAdminProductService.On("FindByID", productId).Return(nil, serviceErr)

				// Action
				resp := mock.MHTTPHandle(http.MethodGet, path, cProduct.FindByID, strconv.Itoa(productId), url.Values{}, nil)

				// Mock assertions
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "FindByID", 1)

				product := &entity.Product{}
				bodyResponse, _ := utils.GetBodyResponse(resp, product)

				//Data Assertion
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Nil(t, bodyResponse.Data)
			})
			t.Run("Not numeric ID", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				// Action
				resp := mock.MHTTPHandle(http.MethodGet, `/{id}`, cProduct.FindByID, "abcd123", url.Values{}, nil)

				// Mock assertions
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "FindByID", 0)

				product := &entity.Product{}
				bodyResponse, _ := utils.GetBodyResponse(resp, product)

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
		product := products[0]
		urlvalues := url.Values{}
		desc := helpers.PointerToString("New Description")
		body := &dto.UpdateProduct{
			Name:        "New Name",
			Description: desc,
			Enabled:     &falseValue,
		}
		productToUpdate := &dto.UpdateProduct{
			ProductID:   product.ProductID,
			Name:        "New Name",
			Description: desc,
			Enabled:     &falseValue,
		}
		returnedProduct := &entity.Product{
			ProductID:   product.ProductID,
			Name:        productToUpdate.Name,
			Slug:        "new-name",
			Description: productToUpdate.Description,
			Enabled:     &falseValue,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
		expectedProduct := &entity.Product{
			ProductID:   product.ProductID,
			Description: productToUpdate.Description,
			Enabled:     &falseValue,
			Name:        productToUpdate.Name,
			Slug:        "new-name",
		}
		t.Run("Should success on", func(t *testing.T) {
			mockAdminProductService := new(mock.AdminProductService)
			cProduct := NewProductController(mockAdminProductService)

			//Expectation
			mockAdminProductService.On("Update", productToUpdate).Return(returnedProduct, nil)
			//Action
			resp := mock.MHTTPHandle(http.MethodPatch, path, cProduct.Update, strconv.Itoa(product.ProductID), urlvalues, body)

			//Mock Assertion
			mockAdminProductService.AssertExpectations(t)
			mockAdminProductService.AssertNumberOfCalls(t, "Update", 1)

			result := &entity.Product{}
			bodyResponse, _ := utils.GetBodyResponse(resp, &result)

			//Data Assertion
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "PRODUCT.UPDATED"}), bodyResponse.Message)
			assert.Empty(t, bodyResponse.Errors)
			assert.Equal(t, expectedProduct, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Update service fails", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				// Expectations
				mockAdminProductService.On("Update", productToUpdate).Return(nil, serviceErr)

				// Action
				resp := mock.MHTTPHandle(http.MethodPatch, path, cProduct.Update, strconv.Itoa(product.ProductID), url.Values{}, productToUpdate)

				// Mock assertions
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "Update", 1)

				product := &entity.Product{}
				bodyResponse, _ := utils.GetBodyResponse(resp, product)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Nil(t, bodyResponse.Data)
			})
			t.Run("Not numeric ID", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				// Action
				resp := mock.MHTTPHandle(http.MethodPatch, `/{id}`, cProduct.Update, "abcd123", url.Values{}, productToUpdate)

				// Mock assertions
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "Update", 0)

				product := &entity.Product{}
				bodyResponse, _ := utils.GetBodyResponse(resp, product)

				//Data Assertion
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				assert.Equal(t, myErrors.ErrIDNotNumeric.Error(), bodyResponse.Errors[0]["error"])
				assert.Nil(t, bodyResponse.Data)
			})
			t.Run("Missing Body", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				// Action
				resp := mock.MHTTPHandle(http.MethodPatch, path, cProduct.Update, strconv.Itoa(product.ProductID), urlvalues, nil)

				// Mock assertions
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "Update", 0)

				result := &entity.Product{}
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
		product := &products[0]
		urlvalues := url.Values{}
		productToCreate := &dto.CreateProduct{
			Name:        product.Name,
			Description: product.Description,
			Enabled:     &trueValue,
		}
		expectedProduct := &entity.Product{
			ProductID:   product.ProductID,
			Name:        product.Name,
			Slug:        product.Slug,
			Description: product.Description,
			Enabled:     &trueValue,
		}
		t.Run("Should success on", func(t *testing.T) {
			mockAdminProductService := new(mock.AdminProductService)
			cProduct := NewProductController(mockAdminProductService)

			//Expectation
			mockAdminProductService.On("Create", productToCreate).Return(product, nil)
			//Action
			resp := mock.MHTTPHandle(http.MethodPost, "/", cProduct.Create, "", urlvalues, productToCreate)

			//Mock Assertion
			mockAdminProductService.AssertExpectations(t)
			mockAdminProductService.AssertNumberOfCalls(t, "Create", 1)

			result := &entity.Product{}
			bodyResponse, _ := utils.GetBodyResponse(resp, &result)

			//Data Assertion
			assert.Equal(t, 201, resp.StatusCode)
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "PRODUCT.CREATED"}), bodyResponse.Message)
			assert.Empty(t, bodyResponse.Errors)
			assert.Equal(t, expectedProduct, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Missing body", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				//Action
				resp := mock.MHTTPHandle(http.MethodPost, "/", cProduct.Create, "", urlvalues, nil)

				//Mock Assertion
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "Create", 0)

				result := &entity.Product{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
			t.Run("Service error", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				//Expectation
				mockAdminProductService.On("Create", productToCreate).Return(nil, serviceErr)

				//Action
				resp := mock.MHTTPHandle(http.MethodPost, "/", cProduct.Create, "", urlvalues, productToCreate)

				//Mock Assertion
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "Create", 1)

				result := &entity.Product{}
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
		product := &products[0]
		urlvalues := url.Values{}
		path := `/{id}`
		t.Run("Should success on", func(t *testing.T) {
			mockAdminProductService := new(mock.AdminProductService)
			cProduct := NewProductController(mockAdminProductService)

			//Expectation
			mockAdminProductService.On("Delete", product.ProductID).Return(nil)

			//Action
			resp := mock.MHTTPHandle(http.MethodDelete, path, cProduct.Delete, strconv.Itoa(product.ProductID), urlvalues, nil)

			//Mock Assertion
			mockAdminProductService.AssertExpectations(t)
			mockAdminProductService.AssertNumberOfCalls(t, "Delete", 1)

			result := &entity.Product{}
			bodyResponse, _ := utils.GetBodyResponse(resp, result)

			//Data Assertion
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, i18n.T(i18n.Message{MessageID: "PRODUCT.DELETED"}), bodyResponse.Message)
			assert.Empty(t, bodyResponse.Errors)
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Service error", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				//Expectation
				mockAdminProductService.On("Delete", product.ProductID).Return(serviceErr)

				//Action
				resp := mock.MHTTPHandle(http.MethodDelete, path, cProduct.Delete, strconv.Itoa(product.ProductID), urlvalues, nil)

				//Mock Assertion
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "Delete", 1)

				result := &entity.Product{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
			t.Run("Not numeric ID", func(t *testing.T) {
				mockAdminProductService := new(mock.AdminProductService)
				cProduct := NewProductController(mockAdminProductService)

				//Action
				resp := mock.MHTTPHandle(http.MethodDelete, `/{id}`, cProduct.Delete, "abcd123", urlvalues, nil)

				//Mock Assertion
				mockAdminProductService.AssertExpectations(t)
				mockAdminProductService.AssertNumberOfCalls(t, "Delete", 0)

				result := &entity.Product{}
				bodyResponse, _ := utils.GetBodyResponse(resp, result)

				//Data Assertion
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				assert.Equal(t, myErrors.ErrIDNotNumeric.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, result)
			})
		})
	})
}
