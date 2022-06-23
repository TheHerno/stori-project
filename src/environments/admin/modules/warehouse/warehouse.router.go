package warehouse

import (
	"net/http"
	"stori-service/src/environments/admin/resources/interfaces"
	"stori-service/src/utils/helpers"

	"github.com/gorilla/mux"
)

type warehouseRouter struct {
	cWarehouse interfaces.IWarehouseController
}

/*
NewWarehouseRouter creates instances of repository, service and controller
then calls all functions for route versions
*/
func NewWarehouseRouter(subRouter *mux.Router, cWarehouse interfaces.IWarehouseController) {
	routerWarehouse := warehouseRouter{cWarehouse}
	routerWarehouse.routes(subRouter)
}

/*
routes assigns controller function for routes
*/
func (r *warehouseRouter) routes(subRouter *mux.Router) {
	subRouter.
		Path("").
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cWarehouse.Index),
		)).
		Methods(http.MethodGet)
	subRouter.
		Path(`/{id}`).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cWarehouse.FindByID),
		)).
		Methods(http.MethodGet)
	subRouter.
		Path(`/{id}`).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cWarehouse.Update),
		)).
		Methods(http.MethodPut)
	subRouter.
		Path(``).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cWarehouse.Create),
		)).
		Methods(http.MethodPost)
	subRouter.
		Path(`/{id}`).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cWarehouse.Delete),
		)).
		Methods(http.MethodDelete)
}
