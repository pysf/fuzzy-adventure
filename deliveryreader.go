package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

func newDeliveryReader(ctx context.Context) *deliveryReader {

	return &deliveryReader{
		wg:  sync.WaitGroup{},
		ctx: ctx,
	}

}

type deliveryReader struct {
	wg  sync.WaitGroup
	ctx context.Context
}

func (d deliveryReader) start(input <-chan downloadResult) <-chan searchEnginInput {

	resultCh := make(chan searchEnginInput, 100)

	go func() {
		defer close(resultCh)

		sendResult := func(r searchEnginInput) {
			select {
			case <-d.ctx.Done():
			case resultCh <- r:
			}
		}

		for {
			select {
			case <-d.ctx.Done():
				return
			case dr, ok := <-input:
				if !ok {
					return
				}

				if dr.err != nil {
					sendResult(searchEnginInput{
						err: dr.err,
					})
					return
				}

				dec := json.NewDecoder(dr.f)
				if _, err := dec.Token(); err != nil {
					sendResult(searchEnginInput{
						err: fmt.Errorf("readFile: Failed to read opening [ , %w", err),
					})
					return
				}

				defer dr.f.Close()

				for dec.More() {

					var delivery delivery
					err := dec.Decode(&delivery)
					if err != nil {
						sendResult(searchEnginInput{
							err: fmt.Errorf("readFile: Failed to decode json, %w", err),
						})
						return
					}

					sendResult(searchEnginInput{
						data: delivery,
					})

				}

				if _, err := dec.Token(); err != nil {
					sendResult(searchEnginInput{
						err: fmt.Errorf("readFile: Failed to read closing ] , %w", err),
					})
				}
			}
		}

	}()

	return resultCh
}

type delivery struct {
	Postcode string
	Recipe   string
	Delivery string
}
