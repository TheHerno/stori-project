package mock

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"

	"github.com/stretchr/testify/mock"
)

type ClientStockMovementService struct {
	mock.Mock
}

/*
Create mock method
*/
func (mock *ClientStockMovementService) Create(newStockMovement *dto.NewStockMovement) (*entity.StockMovement, error) {
	args := mock.Called(newStockMovement)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.StockMovement), nil
	}
	return nil, args.Error(1)
}
