package server

import (
	"github.com/guacamole/microservices/grpc/data"
	"github.com/hashicorp/go-hclog"
	//"google.golang.org/grpc"
	"context"
	"github.com/guacamole/microservices/grpc/protos/currency"
)

type Currency struct{

	rates *data.ExchangeRates
	log hclog.Logger
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency{

	return &Currency{r,l}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {

	c.log.Info("Handle GetRate","base",rr.GetBase(),"destination",rr.GetDest())

	rate,err := c.rates.GetRate(rr.Base.String(),rr.Dest.String())

	if err != nil {
		return nil, err
	}
	return &currency.RateResponse{Rate: rate}, nil
}
