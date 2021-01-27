package main

import (
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc/protos/currency"
	"grpc/server"
	"net"
	"os"
)

func main(){

	log := hclog.Default()

	gs:= grpc.NewServer()
	cs := server.NewCurrency(log)

	currency.RegisterCurrencyServer(gs,cs)

	reflection.Register(gs)

	l,err := net.Listen("tcp",":9999")
	if err != nil{
		log.Error("error listening","error:",err)
		os.Exit(1 )
	}

	log.Info("Starting server at Localhost :9999")

	err = gs.Serve(l)
	if err != nil{
		log.Error("unable to start server","error",err)
	}
}


