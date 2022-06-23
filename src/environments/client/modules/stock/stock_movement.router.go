package stock

import (
	"net/http"
	"stori-service/src/environments/client/resources/interfaces"
	"stori-service/src/utils/helpers"

	"github.com/gorilla/mux"
)

type stockMovementRouter struct {
	cStockMovement interfaces.IStockMovementController
}

/*
NewPStockMovementRouter creates instances of repository, service and controller
then calls all functions for route versions
*/
func NewStockMovementRouter(subRouter *mux.Router, cStockMovement interfaces.IStockMovementController) {
	routerStockMovement := stockMovementRouter{cStockMovement}
	routerStockMovement.routes(subRouter)
}

/*
routes assigns controller function for routes
*/
func (r *stockMovementRouter) routes(subRouter *mux.Router) {
	subRouter.
		Path(``).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cStockMovement.Create),
		)).
		Methods(http.MethodPost)
}
