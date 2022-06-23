package mock

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
)

/*
AdminProductRepository is a IProductRepository mock
*/
type AdminProductRepository struct {
	TransactionalRepository
}

// FindByID mock method
func (mock *AdminProductRepository) FindByID(id int) (*entity.Product, error) {
	args := mock.Called(id)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// Index mock method
func (mock *AdminProductRepository) Index(pagination *dto.Pagination) ([]entity.Product, error) {
	args := mock.Called(pagination)
	result := args.Get(0)
	if result != nil {
		return result.([]entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// Update mock method
func (mock *AdminProductRepository) Update(product *entity.Product) (*entity.Product, error) {
	args := mock.Called(product)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// Create mock method
func (mock *AdminProductRepository) Create(product *entity.Product) (*entity.Product, error) {
	args := mock.Called(product)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// Delete mock method
func (mock *AdminProductRepository) Delete(id int) error {
	args := mock.Called(id)
	return args.Error(0)
}

// CountBySlug mock method
func (mock *AdminProductRepository) CountBySlug(slug string) (int64, error) {
	args := mock.Called(slug)
	result := args.Get(0)
	if result != nil {
		return result.(int64), args.Error(1)
	}
	return 0, args.Error(1)
}
