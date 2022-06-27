package main

import (
	"context"
	"flag"
	"os"
)

const MAX_POSSIBLE_RECIPE int = 2000
const MAX_POSSIBLE_POSTCODE int = 1000000

func main() {

	// start := time.Now()

	defaultFileURL := "https://test-golang-recipes.s3-eu-west-1.amazonaws.com/recipe-calculation-test-fixtures/hf_test_calculation_fixtures.tar.gz"
	fileURL := flag.String("url", defaultFileURL, "[required] file URL (must be a .tar.gz file)")
	postcode := flag.String("postcode", "10120", "postcode")
	from := flag.String("from", "10AM", "form time e.g., 10AM")
	to := flag.String("to", "3PM", "to time e.g., 3PM")

	recipes := recipeFlag{"Potato", "Veggie", "Mushroom"}
	flag.CommandLine.Var(&recipes, "recipe", "comma separated recipes to search (case insensitive) e.g., Potato,Veggie,Mushroom")
	flag.Parse()
	flag.CommandLine.SetOutput(os.Stderr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	downloader, err := newDownloader(ctx)
	if err != nil {
		panic(err)
	}
	defer downloader.cleanup()
	downladerRes := downloader.download(*fileURL)

	deliveryReader := newDeliveryReader(ctx)
	delivryChannel := deliveryReader.start(downladerRes)

	recipeCounter := newRecipeCounter(ctx)
	postcodeFinder := newPostcodeFinder(ctx)
	recipeFinder := newRecipeFinder(ctx, recipes...)
	deliveryFinder := newDeliveryFinder(ctx, findDeliveryQuery{
		postcode: *postcode,
		from:     *from,
		to:       *to,
	})

	searchEngine := newSearchAggregator(ctx, delivryChannel, recipeCounter, deliveryFinder, recipeFinder, postcodeFinder)
	searchResult := searchEngine.faninResult()

	jsonFileMaker := newJSONMaker(ctx, searchResult)
	res := jsonFileMaker.create()

	if result := <-res; result.err != nil {
		panic(result.err)
	}

	// elapsed := time.Since(start)
	// fmt.Fprintf(os.Stderr, "took %s", elapsed)

}
