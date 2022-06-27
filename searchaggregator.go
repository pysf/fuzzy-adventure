package main

import (
	"context"
	"sync"
)

func newSearchAggregator(ctx context.Context, input <-chan searchEnginInput, blocks ...searchEngine) *searchAggregator {
	searchEngine := &searchAggregator{
		input:         input,
		searchEngines: blocks,
		ctx:           ctx,
	}
	searchEngine.tee()
	return searchEngine
}

type searchAggregator struct {
	searchEngines []searchEngine
	ctx           context.Context
	input         <-chan searchEnginInput
}

func (agg *searchAggregator) tee() {

	closeAllInputs := func(sbs []searchEngine) {
		for _, sb := range sbs {
			sb.closeInputChannel()
		}
	}

	go func() {

		defer closeAllInputs(agg.searchEngines)

		for {
			select {
			case <-agg.ctx.Done():
				return
			case d, ok := <-agg.input:
				if !ok {
					return
				}

				if d.err != nil {
					panic(d.err)
				}

				for _, blk := range agg.searchEngines {
					blk.sendData(d.data)
				}
			}

		}

	}()

}

func (eng *searchAggregator) faninResult() chan searchEngineOutput {

	resultCh := make(chan searchEngineOutput, 2)
	var wg sync.WaitGroup

	for _, blk := range eng.searchEngines {
		outputCh := blk.outputChannel()

		wg.Add(1)
		go func() {
			defer wg.Done()

			sendResult := func(r searchEngineOutput) {
				select {
				case <-eng.ctx.Done():
				case resultCh <- r:
				}
			}

			for {
				select {
				case <-eng.ctx.Done():
					return
				case r, ok := <-outputCh:
					if !ok {
						return
					}

					if r.err != nil {
						sendResult(searchEngineOutput{
							err: r.err,
						})
						break
					}

					data := make(map[string]interface{})
					for k, v := range r.data {
						data[k] = v
					}

					sendResult(searchEngineOutput{
						data: data,
					})

				}
			}

		}()

	}

	go func() {
		defer close(resultCh)
		wg.Wait()
	}()

	return resultCh
}

type searchEngineOutput struct {
	err  error
	data map[string]interface{}
}

type searchEnginInput struct {
	err  error
	data interface{}
}
