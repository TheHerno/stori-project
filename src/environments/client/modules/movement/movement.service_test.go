package movement

import (
	goerrors "errors"
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/env"
	"stori-service/src/libs/errors"
	"stori-service/src/utils/constant"
	customMocks "stori-service/src/utils/test/mock"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMovementService(t *testing.T) {
	t.Run("parseLine", func(t *testing.T) {
		expectedID := 1
		currentYear := strconv.Itoa(time.Now().Year())
		expectedDate, _ := time.Parse("1/02/2006", "5/25/"+currentYear)
		expectedQuantity := 1.6
		expectedType := constant.OutcomeType
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Parsing a valid line", func(t *testing.T) {
				sMovement := &movementService{nil, nil}
				line := []string{
					"1",
					"5/25",
					"-1.6",
				}
				movement, err := sMovement.parseLine(line)

				assert.Nil(t, err)
				assert.NotNil(t, movement)
				assert.Equal(t, expectedID, movement.MovementID)
				assert.Equal(t, expectedDate, movement.Date)
				assert.Equal(t, expectedQuantity, movement.Quantity)
				assert.Equal(t, expectedType, movement.Type)
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name string
				line []string
			}{
				{
					name: "Parsing a line with less than 3 elements",
					line: []string{
						"1",
						"5/25",
					},
				},
				{
					name: "Parsing a line with more than 3 elements",
					line: []string{
						"1",
						"5/25",
						"-1.6",
						"1.6",
					},
				},
				{
					name: "Parsing a line with an invalid date",
					line: []string{
						"1",
						"255/25/2012",
						"-1.6",
					},
				},
				{
					name: "Parsing a line with an invalid quantity",
					line: []string{
						"1",
						"5/25",
						"a",
					},
				},
				{
					name: "Parsing a line with an invalid id",
					line: []string{
						"a",
						"5/25",
						"-1.6",
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					sMovement := &movementService{nil, nil}

					movement, err := sMovement.parseLine(tC.line)

					assert.Error(t, err)
					assert.Nil(t, movement)
				})
			}
		})
	})
	t.Run("getPath", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Getting the path", func(t *testing.T) {
				path := getPath(1)
				assert.Equal(t, env.FileRoute+"/customer_1.csv", path)
			})
		})
	})
	t.Run("ProcessFile", func(t *testing.T) {
		var path string
		validLine1 := "1,5/25,+3.5"
		validLine2 := "2,3/20,-1.6"
		validInput := strings.Join([]string{
			"linea 1 sin info",
			validLine1,
			validLine2,
		}, "\n")
		getPathBackup := getPath
		t.Run("Should success on", func(t *testing.T) {
			t.Run("Processing a valid file", func(t *testing.T) {
				currentYear := time.Now().Year()
				expectedMovements := []entity.Movement{
					{
						MovementID: 1,
						Date:       time.Date(currentYear, 5, 25, 0, 0, 0, 0, time.UTC),
						Quantity:   3.5,
						Type:       constant.IncomeType,
						CustomerID: 1,
						Available:  3.5,
					},
					{
						MovementID: 2,
						Date:       time.Date(currentYear, 3, 20, 0, 0, 0, 0, time.UTC),
						Quantity:   1.6,
						Type:       constant.OutcomeType,
						CustomerID: 1,
						Available:  1.9,
					},
				}
				getPath = func(id int) string {
					return "testfile.csv"
				}
				path = getPath(1)
				mockCustomerRepo := new(customMocks.ClientCustomerRepository)
				mockMovementRepo := new(customMocks.ClientMovementRepository)
				sMovement := NewMovementService(mockMovementRepo, mockCustomerRepo)

				// write a fake file
				file, _ := os.Create(path)
				defer file.Close()
				file.WriteString(validInput)

				// mock preparation
				mockCustomerRepo.On("Clone").Return(mockCustomerRepo, nil)
				mockMovementRepo.On("Clone").Return(mockMovementRepo, nil)
				mockMovementRepo.On("Begin", nil).Return(nil)
				mockCustomerRepo.On("Begin", mock.Anything).Return(nil)
				mockMovementRepo.On("Rollback").Return(nil)
				mockCustomerRepo.On("FindAndLockByCustomerID", 1).Return(&customers[0], nil)
				mockMovementRepo.On("GetLastMovementByCustomerID", 1).Return(&entity.Movement{Available: 0}, nil)
				mockMovementRepo.On("BulkCreate", expectedMovements).Return(nil)
				mockMovementRepo.On("Commit").Return(nil)

				// action
				movementList, err := sMovement.ProcessFile(1)

				// mock assertion
				mockCustomerRepo.AssertExpectations(t)
				mockMovementRepo.AssertExpectations(t)
				mockCustomerRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockCustomerRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "BulkCreate", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "GetLastMovementByCustomerID", 1)
				mockCustomerRepo.AssertNumberOfCalls(t, "FindAndLockByCustomerID", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Commit", 1)

				// assertion
				assert.Nil(t, err)
				assert.Equal(t, 2, len(movementList.Movements))

				t.Cleanup(func() {
					os.RemoveAll(path)
					getPath = getPathBackup
				})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("Invalid line", func(t *testing.T) {
				invalidInput := strings.Join([]string{
					"linea 1 sin info",
					"1,5/25,+3.5",
					"2,juan/20,-1.6",
					"a,3/20,-1.6",
				}, "\n")
				getPath = func(id int) string {
					return "testfile.csv"
				}
				path = getPath(1)
				mockCustomerRepo := new(customMocks.ClientCustomerRepository)
				mockMovementRepo := new(customMocks.ClientMovementRepository)
				sMovement := NewMovementService(mockMovementRepo, mockCustomerRepo)

				// write a fake file
				file, _ := os.Create(path)
				defer file.Close()
				file.WriteString(invalidInput)

				// mock preparation
				mockCustomerRepo.On("Clone").Return(mockCustomerRepo, nil)
				mockMovementRepo.On("Clone").Return(mockMovementRepo, nil)
				mockMovementRepo.On("Begin", nil).Return(nil)
				mockCustomerRepo.On("Begin", mock.Anything).Return(nil)
				mockMovementRepo.On("Rollback").Return(nil)
				mockCustomerRepo.On("FindAndLockByCustomerID", 1).Return(&customers[0], nil)
				mockMovementRepo.On("GetLastMovementByCustomerID", 1).Return(&entity.Movement{Available: 0}, nil)

				// action
				movementList, err := sMovement.ProcessFile(1)

				// mock assertion
				mockCustomerRepo.AssertExpectations(t)
				mockMovementRepo.AssertExpectations(t)
				mockCustomerRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockCustomerRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "BulkCreate", 0)
				mockMovementRepo.AssertNumberOfCalls(t, "GetLastMovementByCustomerID", 1)
				mockCustomerRepo.AssertNumberOfCalls(t, "FindAndLockByCustomerID", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Commit", 0)

				// assertion
				assert.Error(t, err)
				assert.Nil(t, movementList)

				t.Cleanup(func() {
					os.RemoveAll(path)
					getPath = getPathBackup
				})
			})
			t.Run("Fails opening file", func(t *testing.T) {
				getPath = func(id int) string {
					return "testfile.csv"
				}
				path = getPath(1)
				mockCustomerRepo := new(customMocks.ClientCustomerRepository)
				mockMovementRepo := new(customMocks.ClientMovementRepository)
				sMovement := NewMovementService(mockMovementRepo, mockCustomerRepo)

				mockCustomerRepo.On("Clone").Return(mockCustomerRepo, nil)
				mockMovementRepo.On("Clone").Return(mockMovementRepo, nil)
				mockMovementRepo.On("Begin", nil).Return(nil)
				mockCustomerRepo.On("Begin", mock.Anything).Return(nil)
				mockMovementRepo.On("Rollback").Return(nil)
				mockCustomerRepo.On("FindAndLockByCustomerID", 1).Return(&customers[0], nil)

				// action
				movementList, err := sMovement.ProcessFile(1)

				// mock assertion
				mockCustomerRepo.AssertExpectations(t)
				mockMovementRepo.AssertExpectations(t)
				mockCustomerRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockCustomerRepo.AssertNumberOfCalls(t, "Begin", 1)
				mockMovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
				mockCustomerRepo.AssertNumberOfCalls(t, "FindAndLockByCustomerID", 1)

				// assertion
				assert.Error(t, err)
				assert.Nil(t, movementList)

				t.Cleanup(func() {
					getPath = getPathBackup
				})
			})
			testCases := []struct {
				name        string
				prepareMock func(*customMocks.ClientMovementRepository, *customMocks.ClientCustomerRepository)
				assertMock  func(*customMocks.ClientMovementRepository, *customMocks.ClientCustomerRepository)
			}{
				{
					name: "Repository fails on FindAndLockByCustomerID",
					prepareMock: func(mockMovementRepo *customMocks.ClientMovementRepository, mockCustomerRepo *customMocks.ClientCustomerRepository) {
						mockCustomerRepo.On("Clone").Return(mockCustomerRepo, nil)
						mockMovementRepo.On("Clone").Return(mockMovementRepo, nil)
						mockMovementRepo.On("Begin", nil).Return(nil)
						mockCustomerRepo.On("Begin", mock.Anything).Return(nil)
						mockMovementRepo.On("Rollback").Return(nil)
						mockCustomerRepo.On("FindAndLockByCustomerID", 1).Return(nil, goerrors.New("repository error"))
					},
					assertMock: func(mockMovementRepo *customMocks.ClientMovementRepository, mockCustomerRepo *customMocks.ClientCustomerRepository) {
						mockCustomerRepo.AssertExpectations(t)
						mockMovementRepo.AssertExpectations(t)
						mockCustomerRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCustomerRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCustomerRepo.AssertNumberOfCalls(t, "FindAndLockByCustomerID", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "GetLastMovementByCustomerID", 0)
						mockMovementRepo.AssertNumberOfCalls(t, "BulkCreate", 0)
						mockMovementRepo.AssertNumberOfCalls(t, "Commit", 0)
					},
				},
				{
					name: "Repository fails on GetLastMovementByCustomerID",
					prepareMock: func(mockMovementRepo *customMocks.ClientMovementRepository, mockCustomerRepo *customMocks.ClientCustomerRepository) {
						mockCustomerRepo.On("Clone").Return(mockCustomerRepo, nil)
						mockMovementRepo.On("Clone").Return(mockMovementRepo, nil)
						mockMovementRepo.On("Begin", nil).Return(nil)
						mockCustomerRepo.On("Begin", mock.Anything).Return(nil)
						mockMovementRepo.On("Rollback").Return(nil)
						mockCustomerRepo.On("FindAndLockByCustomerID", 1).Return(&customers[0], nil)
						mockMovementRepo.On("GetLastMovementByCustomerID", 1).Return(nil, goerrors.New("repository error"))
					},
					assertMock: func(mockMovementRepo *customMocks.ClientMovementRepository, mockCustomerRepo *customMocks.ClientCustomerRepository) {
						mockCustomerRepo.AssertExpectations(t)
						mockMovementRepo.AssertExpectations(t)
						mockCustomerRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCustomerRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCustomerRepo.AssertNumberOfCalls(t, "FindAndLockByCustomerID", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "GetLastMovementByCustomerID", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "BulkCreate", 0)
						mockMovementRepo.AssertNumberOfCalls(t, "Commit", 0)
					},
				},
				{
					name: "Repository fails on BulkCreate",
					prepareMock: func(mockMovementRepo *customMocks.ClientMovementRepository, mockCustomerRepo *customMocks.ClientCustomerRepository) {
						mockCustomerRepo.On("Clone").Return(mockCustomerRepo, nil)
						mockMovementRepo.On("Clone").Return(mockMovementRepo, nil)
						mockMovementRepo.On("Begin", nil).Return(nil)
						mockCustomerRepo.On("Begin", mock.Anything).Return(nil)
						mockMovementRepo.On("Rollback").Return(nil)
						mockCustomerRepo.On("FindAndLockByCustomerID", 1).Return(&customers[0], nil)
						mockMovementRepo.On("GetLastMovementByCustomerID", 1).Return(nil, errors.ErrNotFound)
						mockMovementRepo.On("BulkCreate", mock.AnythingOfType("[]entity.Movement")).Return(goerrors.New("repository error"))
					},
					assertMock: func(mockMovementRepo *customMocks.ClientMovementRepository, mockCustomerRepo *customMocks.ClientCustomerRepository) {
						mockCustomerRepo.AssertExpectations(t)
						mockMovementRepo.AssertExpectations(t)
						mockCustomerRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCustomerRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCustomerRepo.AssertNumberOfCalls(t, "FindAndLockByCustomerID", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "GetLastMovementByCustomerID", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "BulkCreate", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Commit", 0)
					},
				},
				{
					name: "Commit fails",
					prepareMock: func(mockMovementRepo *customMocks.ClientMovementRepository, mockCustomerRepo *customMocks.ClientCustomerRepository) {
						mockCustomerRepo.On("Clone").Return(mockCustomerRepo, nil)
						mockMovementRepo.On("Clone").Return(mockMovementRepo, nil)
						mockMovementRepo.On("Begin", nil).Return(nil)
						mockCustomerRepo.On("Begin", mock.Anything).Return(nil)
						mockMovementRepo.On("Rollback").Return(nil)
						mockCustomerRepo.On("FindAndLockByCustomerID", 1).Return(&customers[0], nil)
						mockMovementRepo.On("GetLastMovementByCustomerID", 1).Return(nil, errors.ErrNotFound)
						mockMovementRepo.On("BulkCreate", mock.AnythingOfType("[]entity.Movement")).Return(nil)
						mockMovementRepo.On("Commit").Return(goerrors.New("commit error"))
					},
					assertMock: func(mockMovementRepo *customMocks.ClientMovementRepository, mockCustomerRepo *customMocks.ClientCustomerRepository) {
						mockCustomerRepo.AssertExpectations(t)
						mockMovementRepo.AssertExpectations(t)
						mockCustomerRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCustomerRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockCustomerRepo.AssertNumberOfCalls(t, "FindAndLockByCustomerID", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "GetLastMovementByCustomerID", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "BulkCreate", 1)
						mockMovementRepo.AssertNumberOfCalls(t, "Commit", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					getPath = func(id int) string {
						return "testfile.csv"
					}
					path = getPath(1)
					mockCustomerRepo := new(customMocks.ClientCustomerRepository)
					mockMovementRepo := new(customMocks.ClientMovementRepository)
					sMovement := NewMovementService(mockMovementRepo, mockCustomerRepo)

					// write a fake file
					file, _ := os.Create(path)
					defer file.Close()
					file.WriteString(validInput)

					// mock preparation
					tC.prepareMock(mockMovementRepo, mockCustomerRepo)

					// action
					movementList, err := sMovement.ProcessFile(1)

					// mock assertion
					mockCustomerRepo.AssertExpectations(t)
					mockMovementRepo.AssertExpectations(t)
					tC.assertMock(mockMovementRepo, mockCustomerRepo)

					// assertion
					assert.Nil(t, movementList)
					assert.Error(t, err)
					t.Cleanup(func() {
						os.RemoveAll(path)
						getPath = getPathBackup
					})
				})
			}
		})
	})
}
