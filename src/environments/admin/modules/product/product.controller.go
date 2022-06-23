package product

import (
	"net/http"
	"stori-service/src/environments/admin/resources/controller"
	"stori-service/src/environments/admin/resources/interfaces"
	"stori-service/src/libs/dto"
	myErrors "stori-service/src/libs/errors"
	"stori-service/src/libs/i18n"
	"stori-service/src/utils"
	"stori-service/src/utils/helpers"
	"stori-service/src/utils/pagination"
)

// struct that implements IProductController
type productController struct {
	controller.AdminController
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
Index extracts pagination from query and calls Index service
*/
func (c *productController) Index(response http.ResponseWriter, request *http.Request) {
	page, err := pagination.GetPaginationFromQuery(request.URL.Query())
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	products, err := c.sProduct.Index(page)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakePaginateResponse(response, products, http.StatusOK, page)
}

/*
FindByID extracts id from params and calls FindByID service
*/
func (c *productController) FindByID(response http.ResponseWriter, request *http.Request) {
	productId, err := helpers.IDFromRequestToInt(request)
	if err != nil {
		c.MakeErrorResponse(response, myErrors.ErrIDNotNumeric)
		return
	}
	product, err := c.sProduct.FindByID(productId)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	c.MakeSuccessResponse(response, product, http.StatusOK, i18n.T(i18n.Message{MessageID: "PRODUCT.FOUND"}))
}

/*
Update gets the body and code from params, then calls Update service
*/
func (c *productController) Update(response http.ResponseWriter, request *http.Request) {
	productToUpdate := &dto.UpdateProduct{}
	if err := utils.GetBodyRequest(request, productToUpdate); err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	productId, err := helpers.IDFromRequestToInt(request)
	if err != nil {
		c.MakeErrorResponse(response, myErrors.ErrIDNotNumeric)
		return
	}
	productToUpdate.ProductID = productId
	productUpdated, err := c.sProduct.Update(productToUpdate)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, productUpdated, http.StatusOK, i18n.T(i18n.Message{MessageID: "PRODUCT.UPDATED"}))
}

/*
Create gets the body, then calls Create service
*/
func (c *productController) Create(response http.ResponseWriter, request *http.Request) {
	productToCreate := &dto.CreateProduct{}
	if err := utils.GetBodyRequest(request, productToCreate); err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	productCreated, err := c.sProduct.Create(productToCreate)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, productCreated, http.StatusCreated, i18n.T(i18n.Message{MessageID: "PRODUCT.CREATED"}))
}

/*
Delete gets the id from params, then calls Delete service
*/
func (c *productController) Delete(response http.ResponseWriter, request *http.Request) {
	productId, err := helpers.IDFromRequestToInt(request)
	if err != nil {
		c.MakeErrorResponse(response, myErrors.ErrIDNotNumeric)
		return
	}
	err = c.sProduct.Delete(productId)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, nil, http.StatusOK, i18n.T(i18n.Message{MessageID: "PRODUCT.DELETED"}))
}
