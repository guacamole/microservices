package server

import (
	"github.com/guacamole/microservices/grpc/data"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"

	//"google.golang.org/grpc"
	"context"
	"github.com/guacamole/microservices/grpc/protos/currency"
)

type Currency struct {
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {

	c := &Currency{r, l, make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest)}
	go c.handleUpdates()

	return c
}

func (c *Currency) handleUpdates() {
	ur := c.rates.MonitorRates(10 * time.Second)

	for range ur {

		c.log.Info("got updated rates")

		//loop ovre subscription
		for k, v := range c.subscriptions {
			//loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDest().String())

				if err != nil {
					c.log.Error("unable to get updated rates", "base", rr.GetBase().String(), "destination", rr.GetDest().String())
				}

				err = k.Send(&currency.RateResponse{Base: rr.Base, Dest: rr.Dest, Rate: r})

				if err != nil {
					c.log.Error("unable to send updated rates", "base", rr.GetBase().String(), "destination", rr.GetDest().String())

				}
			}

		}
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {

	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDest())

	if rr.Base == rr.Dest {

		err := status.Newf(
			codes.InvalidArgument,
			"base %s and destination %s cannot be same",
			rr.Base.String(),
			rr.Dest.String(),
		)

		err, wde := err.WithDetails(rr)
		if wde != nil {
			return nil, wde
		}
		return nil, err.Err()
	}

	rate, err := c.rates.GetRate(rr.Base.String(), rr.Dest.String())

	if err != nil {

		return nil, err
	}
	return &currency.RateResponse{Base: rr.Base, Dest: rr.Dest, Rate: rate}, nil
}

func (c *Currency) SubscribeRates(cur currency.Currency_SubscribeRatesServer) error {

	for {
		rr, err := cur.Recv()
		if err == io.EOF {
			c.log.Info("client closed connection", err)
			break
		}
		if err != nil {
			c.log.Error("error reciving from client", err)
			return err
		}
		c.log.Info("handling client request", "base", rr.GetBase(), "destination", rr.GetDest())

		rss, ok := c.subscriptions[cur]
		if !ok {
			rss = []*currency.RateRequest{}
		}
		rss = append(rss, rr)
		c.subscriptions[cur] = rss
	}

	return nil
}
