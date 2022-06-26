package movement

import (
	"stori-service/src/environments/client/resources/interfaces"

	"github.com/gorilla/mux"
)

type movementRouter struct {
	cMovement interfaces.IMovementController
}

/*
NewPMovementRouter creates instances of repository, service and controller
then calls all functions for route versions
*/
func NewMovementRouter(subRouter *mux.Router, cMovement interfaces.IMovementController) {
	routerMovement := movementRouter{cMovement}
	routerMovement.routes(subRouter)
}

/*
routes assigns controller function for routes
*/
func (r *movementRouter) routes(subRouter *mux.Router) {
	/*subRouter.
	Path(``).
	Handler(helpers.Middleware(
		http.HandlerFunc(r.cMovement.Create),
	)).
	Methods(http.MethodPost)*/
}
