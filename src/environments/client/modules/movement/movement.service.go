package movement

import (
	"bufio"
	goerrors "errors"
	"math"
	"os"
	"stori-service/src/environments/client/resources/interfaces"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/env"
	"stori-service/src/libs/errors"
	"strconv"
	"strings"
	"time"
)

var (
	getPath = func(customerID int) string { // declared here for easy testing with spy
		return env.FileRoute + "/customer_" + strconv.Itoa(customerID) + ".csv"
	}
)

/*
Struct that implements IMovementService
*/
type movementService struct {
	rMovement interfaces.IMovementRepository
	rCustomer interfaces.ICustomerRepository
}

/*
	NewMovementService creates a new service, receives repository by dependency injection
	and returns IRepositoryService, so it needs to implement all its methods
*/
func NewMovementService(rMovement interfaces.IMovementRepository, rCustomer interfaces.ICustomerRepository) interfaces.IMovementService {
	return &movementService{rMovement, rCustomer}
}

/*
ProcessFile takes a customerID, check if the customer exists and process that user file
*/
func (s *movementService) ProcessFile(customerID int) (*dto.MovementList, error) {
	rCustomer := s.rCustomer.Clone().(interfaces.ICustomerRepository)
	rMovement := s.rMovement.Clone().(interfaces.IMovementRepository)
	tx := rMovement.Begin(nil)
	rCustomer.Begin(tx)
	defer rMovement.Rollback()

	customer, err := rCustomer.FindAndLockByCustomerID(customerID) // first check that user exists
	if err != nil {
		return nil, err
	}
	path := getPath(customerID)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip the first line because it doesnt't have data
	var movementList dto.MovementList
	movementList.Customer = customer
	// declare variables outside for loop to avoid memory leaks
	var movement *entity.Movement
	var line []string
	var lastAvailable float64
	// get the last movement of the customer to calculate the new balance
	lastMovement, err := rMovement.GetLastMovementByCustomerID(customerID)
	if !goerrors.Is(err, errors.ErrNotFound) {
		if err != nil {
			return nil, err
		}
		lastAvailable = lastMovement.Available
	}
	for scanner.Scan() {
		// parse line to movement
		line = strings.Split(scanner.Text(), ",")
		// save last available for the next movement
		if movement != nil {
			lastAvailable = movement.Available
		}
		movement, err = s.parseLine(line)
		if err != nil {
			return nil, err
		}
		movement.CustomerID = customerID
		movement.Available = lastAvailable + (movement.Quantity * float64(movement.Type))
		// add movement to list
		movementList.Movements = append(movementList.Movements, *movement)
	}
	err = rMovement.BulkCreate(movementList.Movements)
	if err != nil {
		return nil, err
	}
	err = rMovement.Commit()
	if err != nil {
		return nil, err
	}
	return &movementList, nil
}

/*
parseLine takes a line of the file and returns a movement
*/
func (s *movementService) parseLine(line []string) (*entity.Movement, error) {
	if len(line) != 3 {
		return nil, errors.ErrInvalidFileLine
	}
	var movement entity.Movement
	movementID, err := strconv.Atoi(line[0])
	if err != nil {
		return nil, err
	}
	currentYear := strconv.Itoa(time.Now().Year())
	date, err := time.Parse("1/02/2006", line[1]+"/"+currentYear) // As in the date doesn't have a year, I add the current
	if err != nil {
		return nil, err
	}
	qty, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return nil, err
	}
	movement.Quantity = math.Abs(qty)
	movement.Type = int(qty / math.Abs(qty))
	movement.MovementID = movementID
	movement.Date = date
	return &movement, nil
}
