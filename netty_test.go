package netty

import (
	"testing"
	"time"

	tools "github.com/Pixdigit/goTestTools"
)

func TestEverything(t *testing.T) {
	errChan := make(chan error)

	nc1 := ConnectionConf{
		TCP,
		"localhost",
		49125,
		"localhost",
		49126,
	}
	c1, err := NewConnection(nc1); if err != nil {tools.WrapErr(err, "could not create first test connection", t)}

	nc2 := ConnectionConf{
		TCP,
		"localhost",
		49126,
		"localhost",
		49125,
	}
	c2, err := NewConnection(nc2); if err != nil {tools.WrapErr(err, "could not create second test connection", t)}

	c1.Start()
	c2.Start()

	c2.Accept(errChan)
	err = c1.Connect(); if err != nil {tools.WrapErr(err, "could not connect to second connection instance", t)}
	/*tools.TestAgainstStrings(
		func(str) {return c1.Send(str, str)},
		func() {
			end := false
			time.AfterFunc(2 * time.Second, func() {end = true})
			for result := nil; result == nil; result, err = c1.Poll(){
				if err != nil {
					return "", err
				}
				if end {
					return "", errors.New("timeout")
				}
			}
		},
		"error while transmitting data",
		t,
	)*/
	c1.Send("test", "test2")


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

}
