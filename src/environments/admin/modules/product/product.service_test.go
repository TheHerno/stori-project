package product

import (
	goerrors "errors"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/utils/helpers"
	customMock "stori-service/src/utils/test/mock"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService(t *testing.T) {
	// Fixture
	repositoryErr := goerrors.New("repository error")
	commitErr := goerrors.New("Commit Error")
	t.Run("Index", func(t *testing.T) {
		// fixture
		pagination := dto.NewPagination(1, 20, 0)
		expectedProducts := []entity.Product{
			products[0],
			products[1],
		}

		t.Run("Should success on", func(t *testing.T) {
			mockProductRepo := new(customMock.AdminProductRepository)
			mockCommonProductRepo := new(customMock.CommonProductRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)

			// expectations
			mockProductRepo.On("Index", pagination).Return(expectedProducts, nil)

			// action
			result, err := sProduct.Index(pagination)

			// mock assertion
			mockProductRepo.AssertExpectations(t)
			mockCommonProductRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockProductRepo.AssertNumberOfCalls(t, "Index", 1)

			// data assertion
			assert.Equal(t, expectedProducts, result)
			assert.NoError(t, err)
		})

		t.Run("Should fail on", func(t *testing.T) {
			mockProductRepo := new(customMock.AdminProductRepository)
			mockCommonProductRepo := new(customMock.CommonProductRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)

			// Expectations
			mockProductRepo.On("Index", pagination).Return(nil, repositoryErr)

			// action
			result, err := sProduct.Index(pagination)

			// mock assertion
			mockProductRepo.AssertExpectations(t)
			mockCommonProductRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockProductRepo.AssertNumberOfCalls(t, "Index", 1)

			// data assertion
			assert.Nil(t, result)
			assert.Error(t, err, repositoryErr.Error())
		})
	})

	t.Run("Update", func(t *testing.T) {
		// Fixture
		product := &products[0]
		productToUpdate := &dto.UpdateProduct{
			Name:        "New Product Name",
			ProductID:   1,
			Description: helpers.PointerToString("New Description"),
			Enabled:     &falseValue,
		}
		productUpdated := &entity.Product{}
		copier.Copy(productUpdated, product)
		productUpdated.Description = productToUpdate.Description
		productUpdated.Name = productToUpdate.Name
		productUpdated.Slug = "new-product-name"
		productUpdated.Enabled = productToUpdate.Enabled

		repSlug := "new-edited-product-name"
		productToUpdateRepSlug := &dto.UpdateProduct{
			Name:        "New Edited Product Name",
			ProductID:   1,
			Description: helpers.PointerToString("New Description"),
			Enabled:     &falseValue,
		}
		productUpdatedSlugRep := &entity.Product{}
		copier.Copy(productUpdatedSlugRep, product)
		productUpdatedSlugRep.Name = productToUpdateRepSlug.Name
		productUpdatedSlugRep.Description = productToUpdateRepSlug.Description
		productUpdatedSlugRep.Enabled = productToUpdateRepSlug.Enabled
		productUpdatedSlugRep.Slug = "new-edited-product-name-1"
		t.Run("Should success on", func(t *testing.T) {

			testCases := []struct {
				name         string
				input        *dto.UpdateProduct
				mock         func(*customMock.AdminProductRepository)
				assertMock   func(*customMock.AdminProductRepository)
				assertResult func(*testing.T, *entity.Product)
			}{
				{
					name:  "With no repeated slug",
					input: productToUpdate,
					mock: func(r *customMock.AdminProductRepository) {
						// fixture
						_product := &entity.Product{} // using _product prevents the original products[0] from being modified
						copier.Copy(_product, product)
						// Expectations
						r.On("FindByID", product.ProductID).Return(_product, nil)
						r.On("CountBySlug", productUpdated.Slug).Return(int64(0), nil)
						r.On("Update", productUpdated).Return(productUpdated, nil)
					},
					assertMock: func(r *customMock.AdminProductRepository) {
						r.AssertNumberOfCalls(t, "FindByID", 1)
						r.AssertNumberOfCalls(t, "CountBySlug", 1)
						r.AssertNumberOfCalls(t, "Update", 1)
					},
					assertResult: func(t *testing.T, result *entity.Product) {
						assert.Equal(t, productUpdated, result)
					},
				},
				{
					name:  "With repeated slug",
					input: productToUpdateRepSlug,
					mock: func(r *customMock.AdminProductRepository) {
						// fixture
						_product := &entity.Product{} // using _product prevents the original products[0] from being modified
						copier.Copy(_product, product)
						// Expectations
						r.On("FindByID", product.ProductID).Return(_product, nil)
						r.On("CountBySlug", repSlug).Return(int64(1), nil)
						r.On("Update", productUpdatedSlugRep).Return(productUpdatedSlugRep, nil)
					},
					assertMock: func(r *customMock.AdminProductRepository) {
						r.AssertNumberOfCalls(t, "FindByID", 1)
						r.AssertNumberOfCalls(t, "CountBySlug", 1)
						r.AssertNumberOfCalls(t, "Update", 1)
					},
					assertResult: func(t *testing.T, result *entity.Product) {
						assert.Equal(t, productUpdatedSlugRep, result)
					},
				},
			}

			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockProductRepo := new(customMock.AdminProductRepository)
					mockCommonProductRepo := new(customMock.CommonProductRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)
					tC.mock(mockProductRepo)
					// action
					result, err := sProduct.Update(tC.input)

					// mock assertion
					mockProductRepo.AssertExpectations(t)
					mockCommonProductRepo.AssertExpectations(t)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(mockProductRepo)

					// data assertion
					tC.assertResult(t, result)
					assert.NoError(t, err)
				})
			}

		})
		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      *dto.UpdateProduct
				mock       func(*customMock.AdminProductRepository)
				assertMock func(*customMock.AdminProductRepository)
			}{
				{
					name:  "Invalid update data",
					input: &dto.UpdateProduct{ProductID: 1},
					mock:  func(*customMock.AdminProductRepository) {},
					assertMock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "FindByID", 0)
						mockProductR.AssertNumberOfCalls(t, "Update", 0)
					},
				},
				{
					name:  "Invalid ProductID",
					input: &dto.UpdateProduct{ProductID: 0, Name: "New Name", Description: helpers.PointerToString("New Description")},
					mock:  func(*customMock.AdminProductRepository) {},
					assertMock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "FindByID", 0)
						mockProductR.AssertNumberOfCalls(t, "Update", 0)
					},
				},
				{
					name:  "Fail Finding ID",
					input: productToUpdate,
					mock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.On("FindByID", product.ProductID).Return(nil, repositoryErr)
					},
					assertMock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "Update", 0)
					},
				},
				{
					name:  "Fail to countBySlug",
					input: productToUpdateRepSlug,
					mock: func(r *customMock.AdminProductRepository) {
						// Expectations
						r.On("FindByID", product.ProductID).Return(product, nil)
						r.On("CountBySlug", repSlug).Return(int64(0), repositoryErr)
					},
					assertMock: func(r *customMock.AdminProductRepository) {
						r.AssertNumberOfCalls(t, "FindByID", 1)
						r.AssertNumberOfCalls(t, "CountBySlug", 1)
						r.AssertNumberOfCalls(t, "Update", 0)
					},
				},
				{
					name:  "Fail to Update because repository fails",
					input: productToUpdate,
					mock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.On("FindByID", product.ProductID).Return(product, nil)
						mockProductR.On("CountBySlug", productUpdated.Slug).Return(int64(0), nil)
						mockProductR.On("Update", productUpdated).Return(nil, repositoryErr)
					},
					assertMock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "CountBySlug", 1)
						mockProductR.AssertNumberOfCalls(t, "FindByID", 1)
						mockProductR.AssertNumberOfCalls(t, "Update", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockProductRepo := new(customMock.AdminProductRepository)
					mockCommonProductRepo := new(customMock.CommonProductRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)

					// Expectations
					tC.mock(mockProductRepo)

					// Action
					result, err := sProduct.Update(tC.input)

					// Mock Assertion
					mockProductRepo.AssertExpectations(t)
					mockCommonProductRepo.AssertExpectations(t)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(mockProductRepo)

					// Data Assertion
					assert.Error(t, err)
					assert.Nil(t, result)
				})
			}
		})
	})

	t.Run("Create", func(t *testing.T) {
		// Fixture
		product := &products[0]
		dtoProduct := &dto.CreateProduct{
			Name:        product.Name,
			Description: product.Description,
			Enabled:     &trueValue,
		}
		productToCreate := &entity.Product{
			Name:        product.Name,
			Description: product.Description,
			Enabled:     &trueValue,
			Slug:        product.Slug,
		}
		productCreated := &entity.Product{}
		productCreated.ProductID = 0
		copier.Copy(productCreated, product)
		t.Run("Should success on", func(t *testing.T) {
			mockProductRepo := new(customMock.AdminProductRepository)
			mockCommonProductRepo := new(customMock.CommonProductRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)
			// expectations
			mockProductRepo.On("CountBySlug", productToCreate.Slug).Return(int64(0), nil)
			mockProductRepo.On("Create", productToCreate).Return(productCreated, nil)

			// actions
			result, err := sProduct.Create(dtoProduct)

			// mock assertion
			mockProductRepo.AssertExpectations(t)
			mockCommonProductRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)

			// data assertion
			assert.NoError(t, err)
			assert.Equal(t, productCreated, result)
		})

		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      *dto.CreateProduct
				mock       func(*customMock.AdminProductRepository)
				assertMock func(*testing.T, *customMock.AdminProductRepository)
			}{
				{
					name:  "Invalid create data",
					input: &dto.CreateProduct{},
					mock:  func(*customMock.AdminProductRepository) {},
					assertMock: func(t *testing.T, mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "Create", 0)
					},
				},
				{
					name:  "Create",
					input: dtoProduct,
					mock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.On("CountBySlug", productToCreate.Slug).Return(int64(0), nil)
						mockProductR.On("Create", productToCreate).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "Create", 1)
						mockProductR.AssertNumberOfCalls(t, "CountBySlug", 1)
					},
				},
				{
					name:  "Fail to count slug",
					input: dtoProduct,
					mock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.On("CountBySlug", productToCreate.Slug).Return(int64(0), repositoryErr)
					},
					assertMock: func(t *testing.T, mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "Create", 0)
						mockProductR.AssertNumberOfCalls(t, "CountBySlug", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockProductRepo := new(customMock.AdminProductRepository)
					mockCommonProductRepo := new(customMock.CommonProductRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)

					// Expectations
					tC.mock(mockProductRepo)

					// Action
					result, err := sProduct.Create(tC.input)

					// Mock Assertion
					mockProductRepo.AssertExpectations(t)
					mockCommonProductRepo.AssertExpectations(t)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(t, mockProductRepo)

					// Data assertion
					assert.Error(t, err)
					assert.Nil(t, result)
				})
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		// Fixture
		idToDelete := products[0].ProductID
		t.Run("Should success on", func(t *testing.T) {
			mockProductRepo := new(customMock.AdminProductRepository)
			mockCommonProductRepo := new(customMock.CommonProductRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)
			// expectations
			mockProductRepo.On("Clone").Return(mockProductRepo)
			mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
			mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
			mockProductRepo.On("Begin", nil).Return(nil)
			mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
			mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
			mockProductRepo.On("Rollback").Return(nil)
			mockProductRepo.On("Commit").Return(nil)
			mockCommonProductRepo.On("FindAndLockByID", idToDelete).Return(&products[0], nil)
			mockStockMovementRepo.On("FindLastMovementByProductID", idToDelete).Return(nil, errors.ErrNotFound)
			mockProductRepo.On("Delete", idToDelete).Return(nil)

			// actions
			err := sProduct.Delete(idToDelete)

			// mock assertion
			mockProductRepo.AssertExpectations(t)
			mockCommonProductRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockProductRepo.AssertNumberOfCalls(t, "Clone", 1)
			mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
			mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
			mockProductRepo.AssertNumberOfCalls(t, "Begin", 1)
			mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
			mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
			mockProductRepo.AssertNumberOfCalls(t, "Rollback", 1)
			mockProductRepo.AssertNumberOfCalls(t, "Commit", 1)
			mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
			mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByProductID", 1)
			mockProductRepo.AssertNumberOfCalls(t, "Delete", 1)

			// data assertion
			assert.NoError(t, err)
		})
		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      int
				mock       func(*customMock.AdminProductRepository, *customMock.CommonProductRepository, *customMock.AdminStockMovementRepository)
				assertMock func(*testing.T, *customMock.AdminProductRepository, *customMock.CommonProductRepository, *customMock.AdminStockMovementRepository)
			}{
				{
					name:  "Invalid id",
					input: 0,
					mock: func(mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
					},
					assertMock: func(t *testing.T, mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 0)
						mockProductRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 0)
						mockProductRepo.AssertNumberOfCalls(t, "Rollback", 0)
						mockProductRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 0)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByProductID", 0)
						mockProductRepo.AssertNumberOfCalls(t, "Delete", 0)
					},
				},
				{
					name:  "Delete",
					input: idToDelete,
					mock: func(mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.On("Clone").Return(mockProductRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockProductRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockProductRepo.On("Rollback").Return(nil)
						mockCommonProductRepo.On("FindAndLockByID", idToDelete).Return(&products[0], nil)
						mockStockMovementRepo.On("FindLastMovementByProductID", idToDelete).Return(nil, errors.ErrNotFound)
						mockProductRepo.On("Delete", idToDelete).Return(repositoryErr)
					},
					assertMock: func(t *testing.T, mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByProductID", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Delete", 1)
					},
				},
				{
					name:  "FindAndLockByID",
					input: idToDelete,
					mock: func(mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.On("Clone").Return(mockProductRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockProductRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockProductRepo.On("Rollback").Return(nil)
						mockCommonProductRepo.On("FindAndLockByID", idToDelete).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByProductID", 0)
						mockProductRepo.AssertNumberOfCalls(t, "Delete", 0)
					},
				},
				{
					name:  "There is a last movement",
					input: idToDelete,
					mock: func(mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.On("Clone").Return(mockProductRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockProductRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockProductRepo.On("Rollback").Return(nil)
						mockCommonProductRepo.On("FindAndLockByID", idToDelete).Return(&products[0], nil)
						mockStockMovementRepo.On("FindLastMovementByProductID", idToDelete).Return(&entity.StockMovement{}, nil)
					},
					assertMock: func(t *testing.T, mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByProductID", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Delete", 0)
					},
				},
				{
					name:  "Getting last movement",
					input: idToDelete,
					mock: func(mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.On("Clone").Return(mockProductRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockProductRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockProductRepo.On("Rollback").Return(nil)
						mockCommonProductRepo.On("FindAndLockByID", idToDelete).Return(&products[0], nil)
						mockStockMovementRepo.On("FindLastMovementByProductID", idToDelete).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Commit", 0)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByProductID", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Delete", 0)
					},
				},
				{
					name:  "Commit",
					input: idToDelete,
					mock: func(mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.On("Clone").Return(mockProductRepo)
						mockCommonProductRepo.On("Clone").Return(mockCommonProductRepo)
						mockStockMovementRepo.On("Clone").Return(mockStockMovementRepo)
						mockProductRepo.On("Begin", nil).Return(nil)
						mockCommonProductRepo.On("Begin", mock.Anything).Return(nil)
						mockStockMovementRepo.On("Begin", mock.Anything).Return(nil)
						mockProductRepo.On("Rollback").Return(nil)
						mockCommonProductRepo.On("FindAndLockByID", idToDelete).Return(&products[0], nil)
						mockStockMovementRepo.On("FindLastMovementByProductID", idToDelete).Return(nil, errors.ErrNotFound)
						mockProductRepo.On("Delete", idToDelete).Return(nil)
						mockProductRepo.On("Commit").Return(commitErr)
					},
					assertMock: func(t *testing.T, mockProductRepo *customMock.AdminProductRepository, mockCommonProductRepo *customMock.CommonProductRepository, mockStockMovementRepo *customMock.AdminStockMovementRepository) {
						mockProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Clone", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "Begin", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Rollback", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Commit", 1)
						mockCommonProductRepo.AssertNumberOfCalls(t, "FindAndLockByID", 1)
						mockStockMovementRepo.AssertNumberOfCalls(t, "FindLastMovementByProductID", 1)
						mockProductRepo.AssertNumberOfCalls(t, "Delete", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockProductRepo := new(customMock.AdminProductRepository)
					mockCommonProductRepo := new(customMock.CommonProductRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)
					// Expectations
					tC.mock(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)
					// Action
					err := sProduct.Delete(tC.input)
					// mock assertion
					mockProductRepo.AssertExpectations(t)
					mockCommonProductRepo.AssertExpectations(t)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(t, mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)
					// Data assertion
					assert.Error(t, err)
				})
			}
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		// Fixture
		productToFind := &products[0]
		idToFind := products[0].ProductID
		t.Run("Should success on", func(t *testing.T) {
			mockProductRepo := new(customMock.AdminProductRepository)
			mockCommonProductRepo := new(customMock.CommonProductRepository)
			mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
			sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)
			// expectations
			mockProductRepo.On("FindByID", idToFind).Return(productToFind, nil)

			// actions
			result, err := sProduct.FindByID(idToFind)

			// mock assertion
			mockProductRepo.AssertExpectations(t)
			mockCommonProductRepo.AssertExpectations(t)
			mockStockMovementRepo.AssertExpectations(t)
			mockProductRepo.AssertNumberOfCalls(t, "FindByID", 1)

			// data assertion
			assert.NoError(t, err)
			assert.Equal(t, productToFind, result)
		})
		t.Run("Should fail on", func(t *testing.T) {
			testCases := []struct {
				name       string
				input      int
				mock       func(*customMock.AdminProductRepository)
				assertMock func(*testing.T, *customMock.AdminProductRepository)
			}{
				{
					name:  "Invalid id",
					input: 0,
					mock:  func(*customMock.AdminProductRepository) {},
					assertMock: func(t *testing.T, mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "FindByID", 0)
					},
				},
				{
					name:  "FindByID",
					input: idToFind,
					mock: func(mockProductR *customMock.AdminProductRepository) {
						mockProductR.On("FindByID", idToFind).Return(nil, repositoryErr)
					},
					assertMock: func(t *testing.T, mockProductR *customMock.AdminProductRepository) {
						mockProductR.AssertNumberOfCalls(t, "FindByID", 1)
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.name, func(t *testing.T) {
					mockProductRepo := new(customMock.AdminProductRepository)
					mockCommonProductRepo := new(customMock.CommonProductRepository)
					mockStockMovementRepo := new(customMock.AdminStockMovementRepository)
					sProduct := NewProductService(mockProductRepo, mockCommonProductRepo, mockStockMovementRepo)
					// Expectations
					tC.mock(mockProductRepo)
					// Action
					result, err := sProduct.FindByID(tC.input)
					// mock assertion
					mockProductRepo.AssertExpectations(t)
					mockCommonProductRepo.AssertExpectations(t)
					mockStockMovementRepo.AssertExpectations(t)
					tC.assertMock(t, mockProductRepo)
					// Data assertion
					assert.Error(t, err)
					assert.Nil(t, result)
				})
			}
		})
	})
}
