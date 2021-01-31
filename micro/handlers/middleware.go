package handlers

import (
	"context"
	"github.com/guacamole/microservices/micro/data"
	"net/http"
)

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {

	sig := func(rw http.ResponseWriter, r *http.Request) {

		prod := &data.Product{}

		err := data.FromJSON(prod,r.Body)
		if err != nil {
			http.Error(rw, "error while unmarshalling", http.StatusBadRequest)
		}
		//validate the product

		err = prod.Validate()
		if err != nil {
			http.Error(rw, "error while validating", http.StatusBadRequest)
			p.l.Error("error while validating", err)
			return
		}
		//add product to context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(sig)
}
