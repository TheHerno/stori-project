package warehouse

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

// struct that implements IWarehouseController
type warehouseController struct {
	controller.AdminController
	sWarehouse interfaces.IWarehouseService
}

/*
NewWarehouseController creates a new controller, receives service by dependency injection
and returns IWarehouseController, so needs to implement all its methods
*/
func NewWarehouseController(sWarehouse interfaces.IWarehouseService) interfaces.IWarehouseController {
	return &warehouseController{sWarehouse: sWarehouse}
}

/*
Index extracts pagination from query and calls Index service
*/
func (c *warehouseController) Index(response http.ResponseWriter, request *http.Request) {
	page, err := pagination.GetPaginationFromQuery(request.URL.Query())
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	warehouses, err := c.sWarehouse.Index(page)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakePaginateResponse(response, warehouses, http.StatusOK, page)
}

/*
FindByID extracts id from params and calls FindByID service
*/
func (c *warehouseController) FindByID(response http.ResponseWriter, request *http.Request) {
	warehouseId, err := helpers.IDFromRequestToInt(request)
	if err != nil {
		c.MakeErrorResponse(response, myErrors.ErrIDNotNumeric)
		return
	}
	warehouse, err := c.sWarehouse.FindByID(warehouseId)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	c.MakeSuccessResponse(response, warehouse, http.StatusOK, i18n.T(i18n.Message{MessageID: "WAREHOUSE.FOUND"}))
}

/*
Update gets the body and code from params, then calls Update service
*/
func (c *warehouseController) Update(response http.ResponseWriter, request *http.Request) {
	warehouseToUpdate := &dto.UpdateWarehouse{}
	if err := utils.GetBodyRequest(request, warehouseToUpdate); err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	warehouseId, err := helpers.IDFromRequestToInt(request)
	if err != nil {
		c.MakeErrorResponse(response, myErrors.ErrIDNotNumeric)
		return
	}
	warehouseToUpdate.WarehouseID = warehouseId
	warehouseUpdated, err := c.sWarehouse.Update(warehouseToUpdate)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, warehouseUpdated, http.StatusOK, i18n.T(i18n.Message{MessageID: "WAREHOUSE.UPDATED"}))
}

/*
Create gets the body, then calls Create service
*/
func (c *warehouseController) Create(response http.ResponseWriter, request *http.Request) {
	warehouseToCreate := &dto.CreateWarehouse{}
	if err := utils.GetBodyRequest(request, warehouseToCreate); err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	warehouseCreated, err := c.sWarehouse.Create(warehouseToCreate)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, warehouseCreated, http.StatusCreated, i18n.T(i18n.Message{MessageID: "WAREHOUSE.CREATED"}))
}

/*
Delete gets the id from params, then calls Delete service
*/
func (c *warehouseController) Delete(response http.ResponseWriter, request *http.Request) {
	warehouseId, err := helpers.IDFromRequestToInt(request)
	if err != nil {
		c.MakeErrorResponse(response, myErrors.ErrIDNotNumeric)
		return
	}
	err = c.sWarehouse.Delete(warehouseId)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, nil, http.StatusOK, i18n.T(i18n.Message{MessageID: "WAREHOUSE.DELETED"}))
}
