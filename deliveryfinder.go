package main

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"
)

func newDeliveryFinder(ctx context.Context, query findDeliveryQuery) *deliveryFinder {

	df := &deliveryFinder{
		&pipeline{
			input:  make(chan delivery, 2),
			output: make(chan searchEngineResult, 2),
			ctx:    ctx,
			wg:     sync.WaitGroup{},
		},
	}
	df.start(query)
	return df
}

type deliveryFinder struct {
	*pipeline
}

func (rc *deliveryFinder) sendData(d interface{}) {
	rc.pipeline.sendData(d)
}

func (rc *deliveryFinder) outputChannel() <-chan searchEngineResult {
	return rc.pipeline.output
}

func (rc *deliveryFinder) closeInputChannel() {
	rc.pipeline.closeInputChannel()
}

func (rf *deliveryFinder) start(query findDeliveryQuery) {

	go func() {
		defer close(rf.output)

		sendResult := func(r searchEngineResult) {
			select {
			case <-rf.ctx.Done():
			case rf.output <- r:
			}
		}

		var deliveryCount int
		for {
			select {
			case <-rf.ctx.Done():
				return
			case order, ok := <-rf.input:

				if !ok {
					sendResult(searchEngineResult{
						data: map[string]interface{}{
							"count_per_postcode_and_time": map[string]interface{}{
								"postcode":       query.postcode,
								"from":           query.from,
								"to":             query.to,
								"delivery_count": deliveryCount,
							},
						},
					})
					return
				}

				if order.Postcode == query.postcode {

					inRange, err := isInRange(order, struct {
						from string
						to   string
					}{query.from, query.to})

					if err != nil {
						sendResult(searchEngineResult{
							err: err,
						})
						return
					}

					if *inRange {
						deliveryCount++
					}
				}
			}
		}
	}()

}

func isInRange(order delivery, rangeQuery struct {
	from string
	to   string
}) (*bool, error) {

	layout := "3PM"
	from, err := time.Parse(layout, rangeQuery.from)
	if err != nil {
		return nil, fmt.Errorf("isInRage: invalid query %w", err)
	}

	to, err := time.Parse(layout, rangeQuery.to)
	if err != nil {
		return nil, fmt.Errorf("isInRage: invalid query %w", err)
	}

	rgx := regexp.MustCompile(`\w+ ([\d]+\w{2}) - ([\d]+\w{2})`)
	res := rgx.FindStringSubmatch(order.Delivery)
	if len(res) != 3 {
		return nil, fmt.Errorf("isInRage: failed to extract %v order, got %v", order.Delivery, res)
	}

	deliveryFrom, err := time.Parse(layout, res[1])
	if err != nil {
		return nil, fmt.Errorf("isInRage: failed to create dilivery from %w", err)
	}

	deliveryTo, err := time.Parse(layout, res[2])
	if err != nil {
		return nil, err
	}

	inRange := false
	if (deliveryFrom.After(from) || deliveryFrom.Equal(from)) && (deliveryTo.Before(to) || deliveryTo.Equal(to)) {
		inRange = true
	}

	return &inRange, nil

}

type findDeliveryQuery struct {
	postcode string
	from     string
	to       string
}
