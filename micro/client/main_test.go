package main

import (
	"fmt"
	"micro/client/client"
	"testing"
)

func TestOurClient(t *testing.T) {

	c := client.Default
	fmt.Println(c)
}
