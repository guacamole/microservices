package data

import "testing"

func TestValidateProduct(t *testing.T){

	p := &Product{
		Name: "lassi",
		Price: 3.2,
		SKU: "gfhy-hughy-hfj",
	}

	err := p.Validate()
	if err != nil{
		t.Fatal(err)
	}

}
