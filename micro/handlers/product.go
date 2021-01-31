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
	"github.com/guacamole/microservices/grpc/protos/currency"
	"github.com/guacamole/microservices/micro/data"
	"github.com/hashicorp/go-hclog"
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

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

type Products struct {
	l hclog.Logger
	cc currency.CurrencyClient
	productDB data.ProductsDB
}

func NewProducts(l hclog.Logger, cc currency.CurrencyClient, pdb data.ProductsDB) *Products {

	return &Products{l,cc, pdb}
}
