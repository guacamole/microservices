package main

import (
	"github.com/guacamole/microservices/grpc/data"
	"github.com/guacamole/microservices/grpc/protos/currency"
	"github.com/guacamole/microservices/grpc/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {

	log := hclog.Default()

	gs := grpc.NewServer()

	rates, err := data.NewRates(log)

	if err != nil {
		log.Error("couldn't generate rates", "error", err)
		os.Exit(1)
	}

	cs := server.NewCurrency(rates, log)

	currency.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Error("error listening", "error:", err)
		os.Exit(1)
	}

	log.Info("Starting server at Localhost :8888")

	err = gs.Serve(l)
	if err != nil {
		log.Error("unable to start server", "error", err)
	}
}
