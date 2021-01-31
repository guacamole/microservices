package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/guacamole/microservices/grpc/protos/currency"
	"github.com/guacamole/microservices/micro/data"
	"net/http"
	"strconv"
)

/*
	swagger:route GET /Products listProducts
	Returns a list of products
	Responses:
	200: productsResponse
*/
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Debug("handling GET all method")
	rw.Header().Add("Content-type","application/json")

	lp,err := p.productDB.GetProducts("")
	if err != nil {

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = data.ToJSON(lp,rw)

	if err != nil {
		p.l.Error("unable to serialize the product","error",err)
		http.Error(rw, "error marshalling", http.StatusInternalServerError)
	}
	//p.cc.GetRate()

}

func (p *Products) GetSingleProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("handling GET single method")
	rw.Header().Add("Content-type","application/json")

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	sp,err := p.productDB.GetProductByID(id,"")

	switch err{
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("[error] unable to fetch product")
		rw.WriteHeader(http.StatusNotFound)
		return

	default:
		p.l.Error("[error] unable to fetch product")
		rw.WriteHeader(http.StatusNotFound)
		return

	}

	rr := &currency.RateRequest{
		Base: currency.Currencies(currency.Currencies_value["EUR"]),
		Dest: currency.Currencies(currency.Currencies_value["GBP"]),
	}
	resp, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Error("[error] error getting new rate","error",err)
		return
	}

	sp.Price = sp.Price * resp.Rate

	err = data.ToJSON(sp,rw)
	if err != nil {
		http.Error(rw, "error marshalling", http.StatusInternalServerError)
	}
}


