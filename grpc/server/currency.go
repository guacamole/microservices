package server

import (
	"github.com/hashicorp/go-hclog"
	//"google.golang.org/grpc"
	"context"
	"grpc/protos/currency"
)

type Currency struct{

	log hclog.Logger
}

func NewCurrency(l hclog.Logger) *Currency{

	return &Currency{l}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {

	c.log.Info("Handle GetRate","base",rr.GetBase(),"destination",rr.GetDest())
	return &currency.RateResponse{Rate: 0.5},nil

}