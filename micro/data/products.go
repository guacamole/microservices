package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"
	"github.com/go-playground/validator"
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
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {

	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {

	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku",ValidateSKU)
	return validate.Struct(p)

}
func ValidateSKU(fl validator.FieldLevel) bool {
	//sku =ghfyg-jdgyf-hjdgk

	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches:= re.FindAllString(fl.Field().String(),-1)

	if len(matches) != 1 {
		return false
	}
	return true
}


func GetProducts() Products {

	return productList
}

func AddProduct(p *Product) {

	p.ID = generateID()
	productList = append(productList,p)

}

func UpdateProduct(id int, p *Product) error {

	_,pos,err := findProduct(id)

	if err != nil{
		return err
	}

	p.ID = id
	productList[pos] = p
	return nil
}

func DeleteProduct(id int) error {

	_,pos,err := findProduct(id)
	if err != nil {
		return err
	}

	productList[pos] = productList[len(productList) -1]
	productList[len(productList) -1] = nil
	productList = productList[:len(productList)-1]

	return nil
}

var ErrProductNotFound = fmt.Errorf("product not found")

func findProduct(id int) (*Product,int,error) {

	for i,p := range productList{

		if p.ID == id{
			return p,i,nil
		}

	}
	return nil,0,ErrProductNotFound
}

func generateID() int {

	lp := productList[len(productList)-1]
	return  lp.ID + 1
}

var productList = Products {

	&Product {
		ID:          1,
		Name:        "Cafe Latte",
		Description: "Frothy milk coffee",
		Price:       2.45,
		SKU:         "abc232",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},

	&Product {
		ID:          2,
		Name:        "Espresso",
		Description: "Coffee with no milk",
		Price:       3.45,
		SKU:         "why232",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
