package mock

import (
	"stori-service/src/libs/dto"

	"github.com/stretchr/testify/mock"
)

type ClientMovementService struct {
	mock.Mock
}

// ProcessFile mock method
func (c *ClientMovementService) ProcessFile(customerID int) (*dto.MovementList, error) {
	args := c.Called(customerID)
	result := args.Get(0)
	if result != nil {
		return result.(*dto.MovementList), args.Error(1)
	}
	return nil, args.Error(1)
}
