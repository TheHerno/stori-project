package product

import (
	goerrors "errors"
	"fmt"
	"stori-service/src/environments/admin/resources/interfaces"
	"stori-service/src/environments/common/resources/entity"
	commonInterfaces "stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"

	"github.com/gosimple/slug"
)

/*
	Struct that implements the IProductService interface
*/
type productService struct {
	rProduct       interfaces.IProductRepository
	rCProduct      commonInterfaces.IProductRepository
	rStockMovement interfaces.IStockMovementRepository
}

/*
	NewProductService creates a new service, receives repository by dependency injection
	and returns IRepositoryService, so it needs to implement all its methods
*/
func NewProductService(rProduct interfaces.IProductRepository, rCProduct commonInterfaces.IProductRepository, rStockMovement interfaces.IStockMovementRepository) interfaces.IProductService {
	return &productService{rProduct, rCProduct, rStockMovement}
}

/*
	Index receives pagination and calls Index from repository
*/
func (s *productService) Index(pagination *dto.Pagination) ([]entity.Product, error) {
	return s.rProduct.Index(pagination)
}

/*
createSlug recieves a string and returns a slug
Also checks the count of the slug in the database
*/
func (s *productService) createSlug(name string) (string, error) {
	slug := slug.Make(name)
	count, err := s.rProduct.CountBySlug(slug)
	if err != nil {
		return "", err
	}
	if count > 0 {
		slug = fmt.Sprintf("%s-%d", slug, count)
	}
	return slug, nil
}

/*
	Update receives an UpdateProduct dto and validates it
	Finds Product by ID, replace new values and calls Update from repository
	If there is an error, returns it as a second result
*/
func (s *productService) Update(updateProduct *dto.UpdateProduct) (*entity.Product, error) {
	if err := updateProduct.Validate(); err != nil {
		return nil, err
	}

	product, err := s.FindByID(updateProduct.ProductID)
	if err != nil {
		return nil, err
	}
	if product.Name != updateProduct.Name {
		product.Name = updateProduct.Name
		slug, err := s.createSlug(updateProduct.Name)
		if err != nil {
			return nil, err
		}
		product.Slug = slug
	}
	product.Enabled = updateProduct.Enabled
	product.Description = updateProduct.Description
	return s.rProduct.Update(product)
}

/*
	Create receives a product, validates and creates it
*/
func (s *productService) Create(createProduct *dto.CreateProduct) (*entity.Product, error) {
	if err := createProduct.Validate(); err != nil {
		return nil, err
	}
	product := createProduct.ParseToProduct()
	slug, err := s.createSlug(product.Name)
	if err != nil {
		return nil, err
	}
	product.Slug = slug
	return s.rProduct.Create(product)
}

/*
	Delete receives an ID, check if it has movements and calls Delete from repository
*/
func (s *productService) Delete(id int) error {
	if id < 1 {
		return errors.ErrNotFound
	}
	rStockMovement := s.rStockMovement.Clone().(interfaces.IStockMovementRepository)
	rProduct := s.rProduct.Clone().(interfaces.IProductRepository)
	rCProduct := s.rCProduct.Clone().(commonInterfaces.IProductRepository)
	tx := rProduct.Begin(nil)
	rStockMovement.Begin(tx)
	rCProduct.Begin(tx)
	defer rProduct.Rollback()

	_, err := rCProduct.FindAndLockByID(id)
	if err != nil {
		return err
	}
	movement, err := rStockMovement.FindLastMovementByProductID(id)
	if movement != nil || !goerrors.Is(err, errors.ErrNotFound) { // si hay movimiento no se borra
		return errors.ErrCantDeleteProduct
	}
	err = rProduct.Delete(id)
	if err != nil {
		return err
	}
	err = rProduct.Commit()
	if err != nil {
		return err
	}
	return nil
}

/*
FindByID receives an ID and calls FindByID from repository
*/
func (s *productService) FindByID(id int) (*entity.Product, error) {
	if id < 1 {
		return nil, errors.ErrNotFound
	}
	return s.rProduct.FindByID(id)
}
