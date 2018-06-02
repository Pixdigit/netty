package netty

import (
	"fmt"
	"testing"
)

func TestEverything(t *testing.T) {
	nc1 := ConnectionConf{
		TCP,
		"localhost",
		49125,
		"localhost",
		49126,
	}
	c1, err := NewConnection(nc1)
	if err != nil {
		fmt.Println(err)
	}
	nc2 := ConnectionConf{
		TCP,
		"localhost",
		49126,
		"localhost",
		49125,
	}
	c2, err := NewConnection(nc2)
	if err != nil {
		fmt.Println(err)
	}
	err = c1.Connect()
	if err != nil {
		fmt.Println(err)
	}
	err = c2.Connect()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c1)
	fmt.Println(c2)
}
