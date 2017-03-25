package async

import (
	"fmt"
	"sync"
)

type TypeFunc func(chan bool) error

func Go(fs ...TypeFunc) error {
	cerr := make(chan error, 1)
	ccancel := make(chan bool, len(fs))

	go func() {
		var wg sync.WaitGroup

		wg.Add(len(fs))
		for _, rawF := range fs {
			f := rawF

			go func() {
				defer func() {
					wg.Done()
				}()
				defer func() {
					if e := recover(); e != nil {
						cerr <- fmt.Errorf("async: panic %v", e)
					}
				}()

				if e := f(ccancel); e != nil {
					cerr <- e
				}
			}()
		}

		wg.Wait()
		close(cerr)
	}()

	err := <-cerr
	defer func() {
		defer close(ccancel)
		for _ = range fs {
			ccancel <- err != nil
		}
	}()
	return err
}
