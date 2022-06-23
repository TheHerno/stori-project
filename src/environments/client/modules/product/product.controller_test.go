package product

import (
	"context"
	goerrors "errors"
	"net/http"
	"net/url"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/libs/middleware"
	"stori-service/src/utils"
	"stori-service/src/utils/helpers"
	"stori-service/src/utils/test/mock"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductController(t *testing.T) {
	// Fixture
	serviceErr := goerrors.New("service error")
	user := &middleware.User{
		UserID: 1,
		Email:  "messi@gmail.com",
		Name:   "Lionel Messi",
		Role:   "client",
	}

	t.Run("GetStocks", func(t *testing.T) {
		// fixture
		stocksServiceExpected := []dto.ProductWithStock{
			{
				ProductID:   1,
				Description: helpers.PointerToString("Producto número uno"),
				Name:        "Producto 1",
				Slug:        "product-1",
				Stock:       10,
			},
			{
				ProductID:   2,
				Description: helpers.PointerToString("Producto número dos"),
				Name:        "Producto 2",
				Slug:        "product-2",
				Stock:       150,
			},
			{
				ProductID:   3,
				Description: helpers.PointerToString("Producto número tres"),
				Name:        "Producto 3",
				Slug:        "product-3",
				Stock:       7,
			},
			{
				ProductID:   4,
				Description: helpers.PointerToString("Producto número cuatro"),
				Name:        "Producto 4",
				Slug:        "product-4",
				Stock:       70,
			},
		}
		stocksControllerExpected := []dto.ProductWithStock{
			{
				Description: helpers.PointerToString("Producto número uno"),
				Name:        "Producto 1",
				Slug:        "product-1",
				Stock:       10,
			},
			{
				Description: helpers.PointerToString("Producto número dos"),
				Name:        "Producto 2",
				Slug:        "product-2",
				Stock:       150,
			},
			{
				Description: helpers.PointerToString("Producto número tres"),
				Name:        "Producto 3",
				Slug:        "product-3",
				Stock:       7,
			},
			{
				Description: helpers.PointerToString("Producto número cuatro"),
				Name:        "Producto 4",
				Slug:        "product-4",
				Stock:       70,
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
					mockClientProductService := new(mock.ClientProductService)
					cProduct := NewProductController(mockClientProductService)
					// expectations
					mockClientProductService.On("GetStockList", user.UserID, tC.ExpectedPagination).Return(stocksServiceExpected, nil)

					// action
					resp := mock.MHTTPHandle(http.MethodGet, "/", func(res http.ResponseWriter, req *http.Request) {
						ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, user)
						cProduct.GetStocks(res, req.WithContext(ctx))
					}, "", tC.Querystring, nil)

					// mock assertions
					mockClientProductService.AssertExpectations(t)
					mockClientProductService.AssertNumberOfCalls(t, "GetStockList", 1)

					stocks := []dto.ProductWithStock{}
					bodyResponse, _ := utils.GetBodyResponse(resp, &stocks)

					// data assertions
					assert.EqualValues(t, strconv.Itoa(tC.ExpectedPagination.PageSize), resp.Header.Get("X-pagination-page-size"))
					assert.EqualValues(t, strconv.Itoa(tC.ExpectedPagination.Page), resp.Header.Get("X-pagination-current-page"))
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					assert.Empty(t, bodyResponse.Errors)
					assert.Equal(t, stocksControllerExpected, stocks)
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
					mockClientProductService := new(mock.ClientProductService)
					cProduct := NewProductController(mockClientProductService)

					// action
					resp := mock.MHTTPHandle(http.MethodGet, "/", func(res http.ResponseWriter, req *http.Request) {
						ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, user)
						cProduct.GetStocks(res, req.WithContext(ctx))
					}, "", tC.Querystring, nil)

					// Mock assertions
					mockClientProductService.AssertExpectations(t)
					mockClientProductService.AssertNumberOfCalls(t, "GetStockList", 0)

					stocks := []dto.ProductWithStock{}
					bodyResponse, _ := utils.GetBodyResponse(resp, &stocks)

					// Data assertions
					assert.Equal(t, errors.GetStatusCode(tC.ExpectedError), resp.StatusCode)
					assert.NotEmpty(t, tC.ExpectedError.Error(), bodyResponse.Errors[0]["error"])
					assert.Empty(t, stocks)
				})
			}
			t.Run("GetStockList service fail", func(t *testing.T) {
				pagination := dto.NewPagination(1, 20, 0)
				mockClientProductService := new(mock.ClientProductService)
				cProduct := NewProductController(mockClientProductService)

				//Expectation
				mockClientProductService.On("GetStockList", user.UserID, pagination).Return(nil, serviceErr)

				//Action
				resp := mock.MHTTPHandle(http.MethodGet, "/", func(res http.ResponseWriter, req *http.Request) {
					ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, user)
					cProduct.GetStocks(res, req.WithContext(ctx))
				}, "", url.Values{}, nil)

				//Mock Assertion
				mockClientProductService.AssertExpectations(t)
				mockClientProductService.AssertNumberOfCalls(t, "GetStockList", 1)

				stocks := []dto.ProductWithStock{}
				bodyResponse, _ := utils.GetBodyResponse(resp, &stocks)

				//Data Assertion
				assert.Equal(t, errors.GetStatusCode(serviceErr), resp.StatusCode)
				assert.Equal(t, errors.ErrInternalServer.Error(), bodyResponse.Errors[0]["error"])
				assert.Empty(t, stocks)
			})
		})
	})
}
