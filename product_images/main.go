package main

import (
	"context"
	"github.com/gorilla/mux"
	gohandlers "github.com/gorilla/handlers"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"os"
	"os/signal"
	"github.com/guacamole/microservices/product_images/files"
	"github.com/guacamole/microservices/product_images/handlers"
	"time"
)

func main(){

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:              "product-images",
			Level:             hclog.Debug,
		},
		)

	//create a logger for server from default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})
	//create a storage, use local storage
	//max size 5 MB
	stor,err := files.NewLocal("./imageStore",1024*1000*5)

	if err != nil{
		l.Error("unable to generate storage",err)
		os.Exit(1)
	}

	//create handlers
	fh := handlers.NewFiles(l,stor)
	//mw := handlers.GzipHandler{}

	sm := mux.NewRouter()
	mw := handlers.GzipHandler{}

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	//upload files
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadREST)
	ph.HandleFunc("/",fh.UploadMultipart)

	//get files

	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir("./imageStore"))),
	)
	gh.Use(mw.GzipMiddleware)

	s := http.Server{
		Addr:         ":8888",
		Handler:      ch(sm),
		TLSConfig:    nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     sl,
	}
	// start the server
	go func() {
		l.Info("Starting server", "bind_address",":8888")

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()


	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal,1)
	signal.Notify(c,os.Kill)
	signal.Notify(c,os.Interrupt)

	//block until signal is received

	sig := <- c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete

	ctx,_ := context.WithTimeout(context.Background(),30* time.Second)
	s.Shutdown(ctx)



}