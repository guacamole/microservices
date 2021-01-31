package handlers

import (
	"github.com/gorilla/mux"
	"github.com/guacamole/microservices/micro/data"
	"net/http"
	"strconv"
)

/*
	swagger:route DELETE /Products/{id} deleteProduct
	Deletes a product from the database
	Responses:
	201: noContent
	404: errorNotFound
*/
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Debug("handling DELETE method")

	err := p.productDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "product not found", http.StatusBadRequest)
		return
	}
}
