package main

import (
	"fmt"
	"github.com/guacamole/microservices/micro/client/client"
	"testing"
)

func TestOurClient(t *testing.T) {

	c := client.Default
	fmt.Println(c)
}
