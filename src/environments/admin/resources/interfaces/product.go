package interfaces

import (
	"net/http"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/dto"
)

/*
	IProductController methods to handle requests and responses
*/
type IProductController interface {
	Create(response http.ResponseWriter, request *http.Request)
	Update(response http.ResponseWriter, request *http.Request)
	Delete(response http.ResponseWriter, request *http.Request)
	Index(response http.ResponseWriter, request *http.Request)
	FindByID(response http.ResponseWriter, request *http.Request)
}

/*
	IProductService methods with bussiness logic
*/
type IProductService interface {
	Create(product *dto.CreateProduct) (*entity.Product, error)
	Update(product *dto.UpdateProduct) (*entity.Product, error)
	Index(pagination *dto.Pagination) ([]entity.Product, error)
	FindByID(id int) (*entity.Product, error)
	Delete(id int) error
}

/*
	IProductRepository to interact with entity and database
*/
type IProductRepository interface {
	interfaces.ITransactionalRepository
	Create(product *entity.Product) (*entity.Product, error)
	CountBySlug(slug string) (int64, error)
	Update(product *entity.Product) (*entity.Product, error)
	FindByID(id int) (*entity.Product, error)
	Index(pagination *dto.Pagination) ([]entity.Product, error)
	Delete(productId int) error
}
