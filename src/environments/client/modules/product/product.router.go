package product

import (
	"net/http"
	"stori-service/src/environments/client/resources/interfaces"
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
			http.HandlerFunc(r.cProduct.GetStocks),
		)).
		Methods(http.MethodGet)
}
