# Recipe Stats Calculator

## Implementatios Details

1. Any type that satisfies the searchEngine interface can be added to the SearchAggregator easily
2. The SearchAggregator accepts searcEngineInput type and produces searchEnginOutput type
3. All searchEngines are autonomous so they can utiles goroutine and fanin/fanout patterns.
4. All main and optional features are Implemented.

![rsc](https://github.com/hellofreshdevtests/pysf-recipe-count-test-2020/blob/dev/high-level-architecture.jpg?raw=true)

## How to run

0. To build the docker container run:

   ```sh
   ./setup.sh
   ```

1. To run the CLI app with default parameters run:

   ```sh
   docker run rsc
   ```

2. To read the app help run:
   ```sh
   docker run rsc -help
   ```
   ```sh
   Usage of /rsc:
   -from string
    	form time e.g., 10AM (default "10AM")
   -postcode string
    	postcode (default "10120")
   -recipe value
    	comma separated recipes to search (case insensitive) e.g., Potato,Veggie,Mushroom (default Potato,Veggie,Mushroom)
   -to string
    	to time e.g., 3PM (default "3PM")
   -url string
    	[required] file URL (must be a .tar.gz file) (default "https://test-golang-recipes.s3-eu-west-1.amazonaws.com/recipe-calculation-test-fixtures/hf_test_calculation_fixtures.tar.gz")
   ```
3. You can run the app with different parameters e.g.
   ```sh
   docker run rsc -recipe Pizzas,Veggie
   ```
