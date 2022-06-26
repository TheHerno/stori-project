package mock

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"

	"github.com/stretchr/testify/mock"
)

type ClientMovementService struct {
	mock.Mock
}

/*
Create mock method
*/
func (mock *ClientMovementService) Create(newMovement *dto.NewMovement) (*entity.Movement, error) {
	args := mock.Called(newMovement)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Movement), nil
	}
	return nil, args.Error(1)
}
