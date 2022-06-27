package main

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

type DeliveryFinderResult struct {
	deliverie []delivery
	expected  searchEngineResult
	query     findDeliveryQuery
}

var DeliveryFinderResults = []DeliveryFinderResult{
	{
		query: findDeliveryQuery{
			postcode: "10177",
			from:     "7AM",
			to:       "1PM",
		},
		deliverie: []delivery{

			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 7AM - 1PM",
			},
			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 8AM - 12PM",
			},
			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 6AM - 1PM",
			},
			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 7AM - 2PM",
			},
			{
				Postcode: "10557",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 7AM - 1PM",
			},
		},
		expected: searchEngineResult{
			data: map[string]interface{}{
				"count_per_postcode_and_time": map[string]interface{}{
					"postcode":       "10177",
					"from":           "7AM",
					"to":             "1PM",
					"delivery_count": 2,
				},
			},
		},
	},
	{
		query: findDeliveryQuery{
			postcode: "10557",
			from:     "7AM",
			to:       "1PM",
		},
		deliverie: []delivery{

			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 7AM - 1PM",
			},
			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 8AM - 12PM",
			},
			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 6AM - 1PM",
			},
			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 7AM - 2PM",
			},
			{
				Postcode: "10667",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 7AM - 1PM",
			},
		},
		expected: searchEngineResult{
			data: map[string]interface{}{
				"count_per_postcode_and_time": map[string]interface{}{
					"postcode":       "10557",
					"from":           "7AM",
					"to":             "1PM",
					"delivery_count": 0,
				},
			},
		},
	},
}

func TestDeliveryFinderStart(t *testing.T) {

	for i, c := range DeliveryFinderResults {

		resultCh := make(chan searchEngineResult)
		inputCh := make(chan delivery)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		postcodeFinder := deliveryFinder{
			&pipeline{
				wg:     sync.WaitGroup{},
				output: resultCh,
				input:  inputCh,
				ctx:    ctx,
			},
		}
		postcodeFinder.start(c.query)

		go func() {
			defer close(inputCh)
			for _, d := range c.deliverie {
				inputCh <- d
			}
		}()

		select {
		case res, ok := <-resultCh:
			if !ok {
				break
			}

			if res.err != c.expected.err {
				t.Fatalf("case (%v), expected %v got %v", i, c.expected.data, res.data)
			}

			if !reflect.DeepEqual(c.expected.data, res.data) {
				t.Fatalf("case (%v), \n expected %v \n got      %v", i, c.expected.data, res.data)
			}

		case <-time.After(1 * time.Second):
			t.Fatalf("%d ,timedout!", i)
		}

	}
}
