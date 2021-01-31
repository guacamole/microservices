package data

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/guacamole/microservices/grpc/protos/currency"
	"github.com/hashicorp/go-hclog"
	"regexp"
	"time"
)

// swagger:model
type Product struct {

	// the id of product
	//
	// required: true
	// min: 1
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product


type ProductsDB struct {
	c currency.CurrencyClient
	l hclog.Logger
}


func NewProductsDB(c currency.CurrencyClient,l hclog.Logger) *ProductsDB {

	return &ProductsDB{c,l}
}



func NewProductStore() Products {
	ps := productList
	return ps
}

func NewProduct() *Product{
	return  &Product{}
}


func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", ValidateSKU)
	return validate.Struct(p)

}
func ValidateSKU(fl validator.FieldLevel) bool {
	//sku =ghfyg-jdgyf-hjdgk

	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}
	return true
}

func (p *ProductsDB) GetProducts(currency string) (Products,error) {

	if currency == "" {
		return productList,nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.l.Error("unable to get rate","currency",currency)
		return nil, err
	}

	pr := Products{}

	for _,p := range productList{
		np := *p
		np.Price = np.Price *  rate
		pr = append(pr,&np)

	}
	return pr,nil
}

func (p *ProductsDB) AddProduct(pr Product) {

	maxID := productList[len(productList)-1].ID
	pr.ID = maxID + 1
	productList = append(productList, &pr)
}

func (p *ProductsDB) GetProductByID(id int,currency string) (*Product, error) {
	if id == -1 {
		return nil, fmt.Errorf("invalid id %s", id)
	}

	if currency == "" {
		return productList[id], nil
	}

	rate, err := p.getRate(currency)
	if err !=nil {
		p.l.Error("unable to get rate")
		return nil, err
	}

	np := *productList[id]
	np.Price = np.Price * rate

	return &np, nil

}

func (p *ProductsDB) UpdateProduct(id int, pr *Product) error {

	_, pos, err := findProduct(id)

	if err != nil {
		return err
	}

	pr.ID = id
	productList[pos] = pr
	return nil
}

func (p *ProductsDB) DeleteProduct(id int) error {

	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	productList[pos] = productList[len(productList)-1]
	productList[len(productList)-1] = nil
	productList = productList[:len(productList)-1]

	return nil
}

var ErrProductNotFound = fmt.Errorf("product not found")

func findProduct(id int) (*Product, int, error) {

	for i, p := range productList {

		if p.ID == id {
			return p, i, nil
		}

	}
	return nil, 0, ErrProductNotFound
}

func generateID() int {

	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func (p *ProductsDB) getRate(destination string) (float64,error) {
	rr := &currency.RateRequest{
		Base: currency.Currencies(currency.Currencies_value["EUR"]),
		Dest: currency.Currencies(currency.Currencies_value[destination]),
	}

	resp, err := p.c.GetRate(context.Background(),rr)
	return resp.Rate, err
}

var productList = Products{

	&Product{
		ID:          1,
		Name:        "Cafe Latte",
		Description: "Frothy milk coffee",
		Price:       2.45,
		SKU:         "abc232",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},

	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Coffee with no milk",
		Price:       3.45,
		SKU:         "why232",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
