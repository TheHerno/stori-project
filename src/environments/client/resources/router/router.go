package router

import (
	"stori-service/src/environments/client/modules/product"
	"stori-service/src/environments/client/modules/stock"
	"stori-service/src/environments/client/modules/warehouse"
	commonProduct "stori-service/src/environments/common/modules/product"
	commonWarehouse "stori-service/src/environments/common/modules/warehouse"
	"stori-service/src/libs/database"
	"stori-service/src/libs/middleware"

	"github.com/gorilla/mux"
)

/*
SetupClientRoutes creates all instances for client enviroment and calls each router
*/
func SetupClientRoutes(subRouter *mux.Router) {
	subRouter.Use(middleware.NewAuthMiddleware().HandlerClient())
	stockMovementRoutes(subRouter.PathPrefix("/stock_movement").Subrouter())
	productRoutes(subRouter.PathPrefix("/product").Subrouter())
}

/*
productRoutes creates the router for product module
*/
func productRoutes(subRouter *mux.Router) {
	connection := database.GetTrainingGormConnection()
	rStockMovement := stock.NewStockMovementGormRepo(connection)
	rWarehouse := warehouse.NewWarehouseGormRepo(connection)
	sProduct := product.NewProductService(rWarehouse, rStockMovement)
	cProduct := product.NewProductController(sProduct)
	product.NewProductRouter(subRouter, cProduct)
}

/*
stockMovementRoutes creates the router for stock module
*/
func stockMovementRoutes(subRouter *mux.Router) {
	connection := database.GetTrainingGormConnection()
	rStockMovement := stock.NewStockMovementGormRepo(connection)
	rWarehouse := warehouse.NewWarehouseGormRepo(connection)
	rCWarehouse := commonWarehouse.NewWarehouseGormRepo(connection)
	rProduct := commonProduct.NewProductGormRepo(connection)
	sStockMovement := stock.NewStockMovementService(rStockMovement, rWarehouse, rCWarehouse, rProduct)
	cStockMovement := stock.NewStockMovementController(sStockMovement)
	stock.NewStockMovementRouter(subRouter, cStockMovement)
}
