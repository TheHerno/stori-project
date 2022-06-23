package mock

import (
	"stori-service/src/libs/dto"

	"github.com/stretchr/testify/mock"
)

type ClientProductService struct {
	mock.Mock
}

// GetStockList mock method
func (mock *ClientProductService) GetStockList(userID int, pagination *dto.Pagination) ([]dto.ProductWithStock, error) {
	args := mock.Called(userID, pagination)
	result := args.Get(0)
	if result != nil {
		return result.([]dto.ProductWithStock), args.Error(1)
	}
	return nil, args.Error(1)
}
