package interfaces

import (
	"net/http"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/dto"
)

/*
	IWarehouseController methods to handle requests and responses
*/
type IWarehouseController interface {
	Create(response http.ResponseWriter, request *http.Request)
	Update(response http.ResponseWriter, request *http.Request)
	Delete(response http.ResponseWriter, request *http.Request)
	Index(response http.ResponseWriter, request *http.Request)
	FindByID(response http.ResponseWriter, request *http.Request)
}

/*
	IWarehouseService methods with bussiness logic
*/
type IWarehouseService interface {
	Create(warehouse *dto.CreateWarehouse) (*entity.Warehouse, error)
	Update(warehouse *dto.UpdateWarehouse) (*entity.Warehouse, error)
	Index(pagination *dto.Pagination) (*[]entity.Warehouse, error)
	FindByID(id int) (*entity.Warehouse, error)
	Delete(id int) error
}

/*
	IWarehouseRepository to interact with entity and database
*/
type IWarehouseRepository interface {
	interfaces.ITransactionalRepository
	Create(warehouse *entity.Warehouse) (*entity.Warehouse, error)
	Update(warehouse *entity.Warehouse) (*entity.Warehouse, error)
	Index(pagination *dto.Pagination) (*[]entity.Warehouse, error)
	Delete(warehouseId int) error
}
