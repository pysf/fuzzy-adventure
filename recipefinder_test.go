package main

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

type RecipeFinderResult struct {
	deliverie []delivery
	recipe    []string
	expected  searchEngineResult
}

var RecipeFinderResults = []RecipeFinderResult{
	{
		deliverie: []delivery{
			{
				Postcode: "10118",
				Recipe:   "Spanish One-Pan Chicken",
				Delivery: "Saturday 11AM - 1PM",
			},
			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 6AM - 1PM",
			},
		},
		recipe: []string{"Cheesy"},
		expected: searchEngineResult{
			data: map[string]interface{}{
				"match_by_name": []string{"Cheesy Chicken Enchilada Bake"},
			},
		},
	},
	{
		deliverie: []delivery{
			{
				Postcode: "10118",
				Recipe:   "Spanish One-Pan Chicken",
				Delivery: "Saturday 11AM - 1PM",
			},
			{
				Postcode: "10177",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 6AM - 1PM",
			},
		},
		recipe: []string{"Chicken"},
		expected: searchEngineResult{
			data: map[string]interface{}{
				"match_by_name": []string{"Cheesy Chicken Enchilada Bake", "Spanish One-Pan Chicken"},
			},
		},
	},
}

func TestRecipeFinederStart(t *testing.T) {

	for i, c := range RecipeFinderResults {

		resultCh := make(chan searchEngineResult)
		inputCh := make(chan delivery)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		recipeFinder := recipeFinder{
			&pipeline{
				wg:     sync.WaitGroup{},
				output: resultCh,
				input:  inputCh,
				ctx:    ctx,
			},
		}
		recipeFinder.start(c.recipe...)

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
				t.Fatalf("case (%v), expected %v got %v", i, c.expected.data, res.data)
			}

		case <-time.After(1 * time.Second):
			t.Fatalf("%d ,timedout!", i)
		}

	}

}
