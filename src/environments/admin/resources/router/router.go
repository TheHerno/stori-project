package router

import (
	"stori-service/src/environments/admin/modules/product"
	"stori-service/src/environments/admin/modules/stock"
	"stori-service/src/environments/admin/modules/warehouse"
	commonProduct "stori-service/src/environments/common/modules/product"
	commonWarehouse "stori-service/src/environments/common/modules/warehouse"
	"stori-service/src/libs/database"
	"stori-service/src/libs/middleware"

	"github.com/gorilla/mux"
)

/*
SetupAdminRoutes creates all instances for admin enviroment and calls each router
*/
func SetupAdminRoutes(subRouter *mux.Router) {
	subRouter.Use(middleware.NewAuthMiddleware().HandlerAdmin())
	warehouseRoutes(subRouter.PathPrefix("/warehouse").Subrouter())
	productRoutes(subRouter.PathPrefix("/product").Subrouter())
}

/*
warehouseRoutes creates the router for warehouse module
*/
func warehouseRoutes(subRouter *mux.Router) {
	connection := database.GetTrainingGormConnection()
	rWarehouse := warehouse.NewWarehouseGormRepo(connection)
	rCWarehouse := commonWarehouse.NewWarehouseGormRepo(connection)
	rStockMovement := stock.NewStockMovementGormRepo(connection)
	sWarehouse := warehouse.NewWarehouseService(rWarehouse, rCWarehouse, rStockMovement)
	cWarehouse := warehouse.NewWarehouseController(sWarehouse)
	warehouse.NewWarehouseRouter(subRouter, cWarehouse)
}

/*
productRoutes creates the router for product module
*/
func productRoutes(subRouter *mux.Router) {
	connection := database.GetTrainingGormConnection()
	rProduct := product.NewProductGormRepo(connection)
	rCProduct := commonProduct.NewProductGormRepo(connection)
	rStockMovement := stock.NewStockMovementGormRepo(connection)
	sProduct := product.NewProductService(rProduct, rCProduct, rStockMovement)
	cProduct := product.NewProductController(sProduct)
	product.NewProductRouter(subRouter, cProduct)
}
