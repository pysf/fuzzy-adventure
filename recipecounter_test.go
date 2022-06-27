package main

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

type RecipeCounterResult struct {
	deliverie []delivery
	expected  searchEngineResult
}

var RecipeCounterResults = []RecipeCounterResult{
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
			{
				Postcode: "10151",
				Recipe:   "Cheesy Chicken Enchilada Bake",
				Delivery: "Saturday 10AM - 1PM",
			},
		},
		expected: searchEngineResult{
			data: map[string]interface{}{
				"unique_recipe_count": 2,
				"count_per_recipe": []map[string]interface{}{
					{
						"count":  2,
						"recipe": "Cheesy Chicken Enchilada Bake",
					},
					{
						"count":  1,
						"recipe": "Spanish One-Pan Chicken",
					},
				},
			},
		},
	},
}

func TestRecipeCounterStart(t *testing.T) {

	for i, c := range RecipeCounterResults {

		resultCh := make(chan searchEngineResult)
		inputCh := make(chan delivery)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		recipeCounter := recipeCounter{
			&pipeline{
				wg:     sync.WaitGroup{},
				output: resultCh,
				input:  inputCh,
				ctx:    ctx,
			},
		}
		recipeCounter.start()

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
