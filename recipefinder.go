package main

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const MAX_CONCERRENT_PROCESS int = 3

func newRecipeFinder(ctx context.Context, recipe ...string) *recipeFinder {
	recipeFinder := &recipeFinder{
		&pipeline{
			wg:     sync.WaitGroup{},
			input:  make(chan delivery, 10),
			output: make(chan searchEngineResult, 2),
			ctx:    ctx,
		},
	}

	recipeFinder.start(recipe...)
	return recipeFinder
}

func (rf *recipeFinder) start(receipe ...string) {

	go func() {
		defer close(rf.output)

		channels := make([]<-chan interface{}, MAX_CONCERRENT_PROCESS)

		rgxs := createRecipeRgx(receipe...)
		for i := 0; i < MAX_CONCERRENT_PROCESS; i++ {
			channels[i] = rf.findRecipe(rgxs...)
		}

		uniqueRecipes := make(map[string]bool, MAX_POSSIBLE_RECIPE)
		for r := range fanin(rf.ctx, channels...) {
			uniqueRecipes[r.(string)] = true
		}

		recepieNames := mapValues(&uniqueRecipes)

		select {
		case <-rf.ctx.Done():
		case rf.output <- searchEngineResult{
			data: map[string]interface{}{
				"match_by_name": recepieNames,
			},
		}:
		}

	}()

}

func (rf *recipeFinder) findRecipe(rgxs ...*regexp.Regexp) <-chan interface{} {

	resultCh := make(chan interface{})

	go func() {
		defer close(resultCh)

		sendResult := func(r string) {
			select {
			case <-rf.pipeline.ctx.Done():
			case resultCh <- r:
			}
		}

		for {
			select {
			case <-rf.pipeline.ctx.Done():
				return
			case order, ok := <-rf.pipeline.input:
				if !ok {
					return
				}

				for _, rgx := range rgxs {
					if rgx.MatchString(order.Recipe) {
						sendResult(order.Recipe)
						break
					}
				}
			}
		}
	}()

	return resultCh
}

type recipeFinder struct {
	*pipeline
}

func (rc *recipeFinder) sendData(d interface{}) {
	rc.pipeline.sendData(d)
}

func (rc *recipeFinder) outputChannel() <-chan searchEngineResult {
	return rc.pipeline.output
}

func (rc *recipeFinder) closeInputChannel() {
	rc.pipeline.closeInputChannel()
}

func createRecipeRgx(recipe ...string) []*regexp.Regexp {
	rgxs := make([]*regexp.Regexp, 0, len(recipe))
	for _, rcp := range recipe {
		rgxs = append(rgxs, regexp.MustCompile(fmt.Sprintf("(?i)%v", rcp)))
	}
	return rgxs
}

func mapValues(recipes *map[string]bool) []string {

	result := make([]string, 0, len(*recipes))

	for k := range *recipes {
		result = append(result, k)
	}
	sort.Strings(result)

	return result
}

type recipeFlag []string

func (i *recipeFlag) String() string {
	return strings.Join(*i, ",")
}

func (i *recipeFlag) Set(value string) error {
	*i = strings.Split(value, ",")
	return nil
}
