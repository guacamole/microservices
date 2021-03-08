package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/guacamole/microservices/grpc/protos/currency"
	"github.com/guacamole/microservices/micro/data"
	"github.com/guacamole/microservices/micro/handlers"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main(){

	l := hclog.Default()

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	//create client
	cc := currency.NewCurrencyClient(conn)
	ps := data.NewProductsDB(cc,l)
	ph := handlers.NewProducts(l,cc, *ps)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/",ph.GetProducts)
	getRouter.HandleFunc("/",ph.GetProducts).Queries("currency","{[A-Z]{3}}")

	getRouter.HandleFunc("/product/{id:[0-9]+}",ph.GetSingleProduct)
	getRouter.HandleFunc("/product/{id:[0-9]+}",ph.GetSingleProduct).Queries("currency","{[A-Z]{3}}")

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/",ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}",ph.DeleteProduct)

	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops,nil)

	getRouter.Handle("/swagger.yaml",http.FileServer(http.Dir("./")))
	getRouter.Handle("/docs",sh)

	s := http.Server{
		Addr:              ":9009",
		Handler:           sm,
		TLSConfig:         nil,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		ErrorLog: l.StandardLogger(&hclog.StandardLoggerOptions{}),
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			l.Error("error connecting to the server", err)
			return
		}
	}()
	l.Info("started on port ",s.Addr)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan,os.Interrupt)
	signal.Notify(sigChan,os.Kill)

	sig := <-sigChan
	l.Info("received terminate,graceful shutdown ",sig)

	tc, cancel := context.WithTimeout(context.Background(), 30 *time.Second)

	defer cancel()

	if err := s.Shutdown(tc); err != nil {
		l.Error("error shutting down: ", err)
	}


}