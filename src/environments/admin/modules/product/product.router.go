package product

import (
	"net/http"
	"stori-service/src/environments/admin/resources/interfaces"
	"stori-service/src/utils/helpers"

	"github.com/gorilla/mux"
)

type productRouter struct {
	cProduct interfaces.IProductController
}

/*
NewProductRouter creates instances of repository, service and controller
then calls all functions for route versions
*/
func NewProductRouter(subRouter *mux.Router, cProduct interfaces.IProductController) {
	routerProduct := productRouter{cProduct}
	routerProduct.routes(subRouter)
}

/*
routes assigns controller function for routes
*/
func (r *productRouter) routes(subRouter *mux.Router) {
	subRouter.
		Path("").
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cProduct.Index),
		)).
		Methods(http.MethodGet)
	subRouter.
		Path(`/{id}`).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cProduct.FindByID),
		)).
		Methods(http.MethodGet)
	subRouter.
		Path(`/{id}`).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cProduct.Update),
		)).
		Methods(http.MethodPut)
	subRouter.
		Path(``).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cProduct.Create),
		)).
		Methods(http.MethodPost)
	subRouter.
		Path(`/{id}`).
		Handler(helpers.Middleware(
			http.HandlerFunc(r.cProduct.Delete),
		)).
		Methods(http.MethodDelete)
}
