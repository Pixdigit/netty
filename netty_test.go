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
	c2, err := NewConnection(nc2);	if err != nil {tools.WrapErr(err, "could not create second test connection", t)}

	c1.Start()
	c2.Start()

	c2.Accept(errChan)
	err = c1.Connect();	if err != nil {tools.WrapErr(err, "could not connect to second connection instance", t)}
	tools.TestAgainstStrings(
		func(str string) (error) {return c1.Send(str, str)},
		func() (string, error){
			key, value, err := c2.Recv();	if err != nil {return "", err}
			if key != value.(string) {
				return "", nil
			}
			return key, nil
		},
		"error while transmitting data",
		t,
	)
	err = c1.Send("test", "test2");	if err != nil {tools.WrapErr(err, "could not send data", t)}


	end := time.After(1 * time.Second)

	run := true
	for run {
		select {
		case err := <-errChan:
			if err != nil {
				tools.WrapErr(err, "error", t)
			} else {
				t.Log("some function pushed error with value nil")
			}
		case <-end:
			run = false
		default:
		}
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
