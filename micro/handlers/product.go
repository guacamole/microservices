/*
  Package classification Product API

  Documentation for Package

  Schemes: http
  BasePath:/
  Version: 1.0.0

  Consumes:
    - application/json

  Produces:
    - application/json

  swagger:meta
*/
package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"micro/data"
	"net/http"
	"strconv"
)

// List of the products returned in the response
// swagger:response productsResponse
type productsResponse struct {
	// list of all products
	// in: body
	Body []data.Product
}

// swagger:parameters deleteProduct updateProduct
type productIDParameterWrapper struct{
	// id of the product to be removed from the database
	// in: path
	// required: true
	ID int `json:"id"`

}

// represents single product
// swagger:parameters addProduct updateProduct
type productResponse struct{
	//newly created product
	//in: body
	Body data.Product
}

// swagger:response noContent
type productNoContentWrapper struct{

}
// swagger:response errorValidation
type errorValidationWrapper struct{
	// validation error
	// in: body
	Body ValidationError
}
// swagger:response errorResponse
type errorResponseWrapper struct{
	// description of error
	// in: body
	Body GenericError
}
// swagger: response errorNotFound
type errorNotFoundWrapper struct{
	// product not found
	// in: path
	ID int `json:"id"`
}


type GenericError struct{
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}


type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {

	return &Products{l}
}

/*
	swagger:route GET /Products listProducts
	Returns a list of products
	Responses:
	200: productsResponse
*/
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("handling GET method")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "error marshalling", http.StatusInternalServerError)
	}

}
/*
	swagger:route POST /Products/ addProduct
	adds a product to the list of products
	Responses:
	200: productsResponse
	422: errorValidation
	501: errorResponse
*/
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("handling POST method")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}


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

	p.l.Println("handling PUT method")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err := data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product  not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product  not found", http.StatusInternalServerError)
		return
	}

}
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

	p.l.Println("handling DELETE method")

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "product not found", http.StatusBadRequest)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {

	sig := func(rw http.ResponseWriter, r *http.Request) {

		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "error while unmarshalling", http.StatusBadRequest)
		}
		//validate the product

		err = prod.Validate()
		if err != nil {
			http.Error(rw, "error while validating", http.StatusBadRequest)
			p.l.Println("error while validating", err)
			return
		}
		//add product to context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(sig)
}
