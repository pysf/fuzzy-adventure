package main

import (
	"context"
	"sort"
	"sync"
)

func newRecipeCounter(ctx context.Context) *recipeCounter {

	rc := &recipeCounter{
		&pipeline{
			input:  make(chan delivery, 10),
			output: make(chan searchEngineResult, 2),
			ctx:    ctx,
			wg:     sync.WaitGroup{},
		},
	}

	rc.start()
	return rc
}

type recipeCounter struct {
	*pipeline
}

func (rc *recipeCounter) sendData(d interface{}) {
	rc.pipeline.sendData(d)
}

func (rc *recipeCounter) outputChannel() <-chan searchEngineResult {
	return rc.pipeline.output
}

func (rc *recipeCounter) closeInputChannel() {
	rc.pipeline.closeInputChannel()
}

func (rc *recipeCounter) start() {

	go func() {
		defer close(rc.pipeline.output)
		uniqueRecipe := make(map[string]int, MAX_POSSIBLE_RECIPE)

		for {
			select {
			case <-rc.pipeline.ctx.Done():
				return
			case order, ok := <-rc.pipeline.input:
				if !ok {
					select {
					case <-rc.pipeline.ctx.Done():
						return
					case rc.pipeline.output <- searchEngineResult{
						data: map[string]interface{}{
							"unique_recipe_count": len(uniqueRecipe),
							"count_per_recipe":    recipeToArray(uniqueRecipe),
						},
					}:
					}
					return
				}

				uniqueRecipe[order.Recipe] += 1
			}

		}

	}()

}

func recipeToArray(recipes map[string]int) []map[string]interface{} {

	keys := make([]string, 0, len(recipes))

	for k := range recipes {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	result := make([]map[string]interface{}, 0, len(recipes))
	for _, name := range keys {
		result = append(result, map[string]interface{}{
			"count":  recipes[name],
			"recipe": name,
		})
	}
	return result
}
