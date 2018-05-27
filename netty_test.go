package netty

import (
	"fmt"
	"testing"
)

func TestEverything(t *testing.T) {
	nc := NetConf{
		TCP,
		"localhost",
		49125,
		"localhost",
		49125,
	}
	conn, err := NewConnector(nc)
	fmt.Println(err)
	fmt.Println(conn)
}
