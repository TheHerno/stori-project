package router

import (
	movement "stori-service/src/environments/client/modules/movement"
	"stori-service/src/environments/client/modules/user"
	"stori-service/src/libs/database"

	"github.com/gorilla/mux"
)

/*
SetupClientRoutes creates all instances for client enviroment and calls each router
*/
func SetupClientRoutes(subRouter *mux.Router) {
	movementRoutes(subRouter.PathPrefix("/movement").Subrouter())
}

/*
movementRoutes creates the router for movement module
*/
func movementRoutes(subRouter *mux.Router) {
	connection := database.GetStoriGormConnection()
	rMovement := movement.NewMovementGormRepo(connection)
	rUser := user.NewUserGormRepo(connection)
	sMovement := movement.NewMovementService(rMovement, rUser)
	cMovement := movement.NewMovementController(sMovement)
	movement.NewMovementRouter(subRouter, cMovement)
}
