package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/guacamole/microservices/grpc/protos/currency"
	"log"
	"github.com/guacamole/microservices/micro/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main(){

	l := log.New(os.Stdout,"product-api",log.LstdFlags)
	ph := handlers.NewProducts(l)


	currency.NewCurrencyClient()
 
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/",ph.GetProducts)

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
		Addr:              ":8888",
		Handler:           sm,
		TLSConfig:         nil,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}


	go func() {
		if err := s.ListenAndServe(); err != nil {
			l.Fatal("error connecting to the server", err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan,os.Interrupt)
	signal.Notify(sigChan,os.Kill)

	sig := <-sigChan
	l.Println("received terminate,graceful shutdown ",sig)

	tc, cancel := context.WithTimeout(context.Background(), 30 *time.Second)

	defer cancel()

	if err := s.Shutdown(tc); err != nil {
		l.Println("error shutting down: ", err)
	}


}