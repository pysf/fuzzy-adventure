package main

import (
	"context"
	"sync"
)

type pipeline struct {
	wg     sync.WaitGroup
	output chan searchEngineResult
	input  chan delivery
	ctx    context.Context
}

func (p *pipeline) sendData(d interface{}) {

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		select {
		case <-p.ctx.Done():
		case p.input <- d.(delivery):
		}
	}()
}

func (p *pipeline) outputChannel() <-chan searchEngineResult {
	return p.output
}

func (p *pipeline) closeInputChannel() {
	go func() {
		defer close(p.input)
		p.wg.Wait()
	}()
}
