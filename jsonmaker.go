package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

func newJSONMaker(ctx context.Context, input chan searchEngineOutput) *jsonMaker {
	jsonFileMaker := &jsonMaker{
		ctx:   ctx,
		input: input,
	}
	return jsonFileMaker
}

type jsonMaker struct {
	ctx   context.Context
	input chan searchEngineOutput
}

func (j *jsonMaker) create() chan jsonMakerResult {

	resultCh := make(chan jsonMakerResult)

	go func() {
		defer close(resultCh)
		data := make(map[string]interface{})

		for {
			select {
			case <-j.ctx.Done():
				return
			case d, ok := <-j.input:

				if !ok {
					d, err := json.MarshalIndent(data, "", "	")
					if err != nil {
						select {
						case <-j.ctx.Done():
						case resultCh <- jsonMakerResult{
							err: err,
						}:
						}
						return
					}

					fmt.Fprintf(os.Stdout, "%v\n", string(d))
					select {
					case <-j.ctx.Done():
					case resultCh <- jsonMakerResult{}:
					}
					return
				}

				if d.err != nil {
					panic(d.err)
				}

				for k, v := range d.data {
					data[k] = v
				}

			}
		}
	}()

	return resultCh
}

type jsonMakerResult struct {
	err error
}
