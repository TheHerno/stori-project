package product

import (
	"net/http"
	"stori-service/src/environments/client/resources/controller"
	"stori-service/src/environments/client/resources/interfaces"
	"stori-service/src/libs/middleware"
	"stori-service/src/utils/pagination"
)

// struct that implements IProductController
type productController struct {
	controller.ClientController
	sProduct interfaces.IProductService
}

/*
NewProductController creates a new controller, receives service by dependency injection
and returns IProductController, so needs to implement all its methods
*/
func NewProductController(sProduct interfaces.IProductService) interfaces.IProductController {
	return &productController{sProduct: sProduct}
}

/*
GetStocks takes the UserID from the request and returns a list of products from his warehouse with stock
*/
func (c *productController) GetStocks(response http.ResponseWriter, request *http.Request) {
	page, err := pagination.GetPaginationFromQuery(request.URL.Query())
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	user := request.Context().Value(middleware.ContextKeyUser).(*middleware.User)
	stocks, err := c.sProduct.GetStockList(user.UserID, page)

	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakePaginateResponse(response, stocks, http.StatusOK, page)
}
