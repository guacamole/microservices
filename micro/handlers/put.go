package handlers

import (
	"github.com/gorilla/mux"
	"github.com/guacamole/microservices/micro/data"
	"net/http"
	"strconv"
)

/*
	swagger:route PUT /Products/{id} updateProduct
	updates an existing product in the list of products
	Responses:
	201: noContent
	404: errorResponse
	422: errorValidation
*/
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Debug("handling PUT method")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err := p.productDB.UpdateProduct(id,prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product  not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product  not found", http.StatusInternalServerError)
		return
	}

}
