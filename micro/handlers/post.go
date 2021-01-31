package handlers

import (
	"github.com/guacamole/microservices/micro/data"
	"net/http"
)

/*
	swagger:route POST /Products/ addProduct
	adds a product to the list of products
	Responses:
	200: productsResponse
	422: errorValidation
	501: errorResponse
*/
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {

	p.l.Debug("handling POST method")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Info("Prod: %#v", prod)

	p.productDB.AddProduct(prod)
}
