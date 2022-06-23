package interfaces

import (
	"net/http"
	"stori-service/src/libs/dto"
)

/*
	IProductService methods with bussiness logic
*/
type IProductService interface {
	GetStockList(userID int, pagination *dto.Pagination) ([]dto.ProductWithStock, error)
}

/*
	IProductController methods to handle requests and responses
*/
type IProductController interface {
	GetStocks(response http.ResponseWriter, request *http.Request)
}
