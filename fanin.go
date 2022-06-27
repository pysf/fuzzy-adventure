package main

import (
	"context"
	"sync"
)

func fanin(ctx context.Context, chans ...<-chan interface{}) <-chan interface{} {

	var wg sync.WaitGroup
	multiplexedCh := make(chan interface{})

	multiplexer := func(ch <-chan interface{}) {
		defer wg.Done()
		for d := range ch {
			select {
			case <-ctx.Done():
				return
			case multiplexedCh <- d:
			}
		}
	}

	wg.Add(len(chans))

	for _, ch := range chans {
		go multiplexer(ch)
	}

	go func() {
		defer close(multiplexedCh)
		wg.Wait()
	}()

	return multiplexedCh
}
