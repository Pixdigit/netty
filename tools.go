package netty

func pushErr(errChan chan error, err error) {
	go func() {
		errChan <- err
	}()
}
