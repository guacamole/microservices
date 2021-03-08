package data

import (
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {

	er := &ExchangeRates{l, map[string]float64{}}

	err := er.getRates()
	return er, err
}

func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}

	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", dest)
	}

	return dr / br, nil

}
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {

	ret := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)

		for {
			select {
			case <-ticker.C:
				for k, v := range e.rates {
					change := (rand.Float64() / 10)
					// is this a postive or negative change
					direction := rand.Intn(1)

					if direction == 0 {
						// new value with be min 90% of old
						change = 1 - change
					} else {
						// new value will be 110% of old
						change = 1 + change
					}

					// modify the rate
					e.rates[k] = v * change
				}
				ret <- struct{}{}
			}
		}
	}()
	return ret
}



func (e *ExchangeRates) getRates() error {

	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("expected error code 200, got: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := Cubes{}

	byts, err := ioutil.ReadAll(resp.Body)
	err = xml.Unmarshal(byts, &md)
	if err != nil {
		e.log.Error("error unmarshaling xml", err)
		return err
	}

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		e.rates[c.Currency] = r
	}
	e.rates["EUR"] = 1

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
