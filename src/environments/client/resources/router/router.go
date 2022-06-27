package router

import (
	"stori-service/src/environments/client/modules/customer"
	movement "stori-service/src/environments/client/modules/movement"
	"stori-service/src/libs/database"

	"github.com/gorilla/mux"
)

/*
SetupClientRoutes creates all instances for client enviroment and calls each router
*/
func SetupClientRoutes(subRouter *mux.Router) {
	movementRoutes(subRouter.PathPrefix("/client-movements").Subrouter())
}

/*
movementRoutes creates the router for movement module
*/
func movementRoutes(subRouter *mux.Router) {
	connection := database.GetStoriGormConnection()
	rMovement := movement.NewMovementGormRepo(connection)
	rCustomer := customer.NewCustomerGormRepo(connection)
	sMovement := movement.NewMovementService(rMovement, rCustomer)
	cMovement := movement.NewMovementController(sMovement)
	movement.NewMovementRouter(subRouter, cMovement)
}
