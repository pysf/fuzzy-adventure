# Recipe Stats Calculator

## Descriptions

It was a very interesting task to do, I learned and enjoyed it.
It is worth mentioning that I learned Golang recently (nearly two months ago) It's my second app with Go and all comments are welcome. (The First app is [GoPipelines](https://github.com/pysf/go-pipelines) if you are interested.)  
One last thing, You may see commits from two different usernames (parmenides and pysf), they are both me! sorry for that it just happened by mistake.

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

## Instructions

In the given assignment we suggest you to process an automatically generated JSON file with recipe data and calculated some stats.

1. Clone this repository.
2. Create a new branch called `dev`.
3. Create a pull request from your `dev` branch to the master branch.
4. Reply to the thread you're having with your recruiter telling them we can start reviewing your code

## Given

Json fixtures file with recipe data. Download [Link](https://test-golang-recipes.s3-eu-west-1.amazonaws.com/recipe-calculation-test-fixtures/hf_test_calculation_fixtures.tar.gz)

_Important notes_

1. Property value `"delivery"` always has the following format: "{weekday} {h}AM - {h}PM", i.e. "Monday 9AM - 5PM"
2. The number of distinct postcodes is lower than `1M`, one postcode is not longer than `10` chars.
3. The number of distinct recipe names is lower than `2K`, one recipe name is not longer than `100` chars.

## Functional Requirements

1. Count the number of unique recipe names.
2. Count the number of occurences for each unique recipe name (alphabetically ordered by recipe name).
3. Find the postcode with most delivered recipes.
4. Count the number of deliveries to postcode `10120` that lie within the delivery time between `10AM` and `3PM`, examples _(`12AM` denotes midnight)_:
   - `NO` - `9AM - 2PM`
   - `YES` - `10AM - 2PM`
5. List the recipe names (alphabetically ordered) that contain in their name one of the following words:
   - Potato
   - Veggie
   - Mushroom

## Non-functional Requirements

1. The application is packaged with [Docker](https://www.docker.com/).
2. Setup scripts are provided.
3. The submission is provided as a `CLI` application.
4. The expected output is rendered to `stdout`. Make sure to render only the final `json`. If you need to print additional info or debug, pipe it to `stderr`.
5. It should be possible to (implementation is up to you):  
   a. provide a custom fixtures file as input  
   b. provide custom recipe names to search by (functional reqs. 5)  
   c. provide custom postcode and time window for search (functional reqs. 4)

## Expected output

Generate a JSON file of the following format:

```json5
{
  unique_recipe_count: 15,
  count_per_recipe: [
    {
      recipe: "Mediterranean Baked Veggies",
      count: 1,
    },
    {
      recipe: "Speedy Steak Fajitas",
      count: 1,
    },
    {
      recipe: "Tex-Mex Tilapia",
      count: 3,
    },
  ],
  busiest_postcode: {
    postcode: "10120",
    delivery_count: 1000,
  },
  count_per_postcode_and_time: {
    postcode: "10120",
    from: "11AM",
    to: "3PM",
    delivery_count: 500,
  },
  match_by_name: [
    "Mediterranean Baked Veggies",
    "Speedy Steak Fajitas",
    "Tex-Mex Tilapia",
  ],
}
```

## Review Criteria

We expect that the assignment will not take more than 3 - 4 hours of work. In our judgement we rely on common sense
and do not expect production ready code. We are rather instrested in your problem solving skills and command of the programming language that you chose.

It worth mentioning that we will be testing your submission against different input data sets.

**General criteria from most important to less important**:

1. Functional and non-functional requirements are met.
2. Prefer application efficiency over code organisation complexity.
3. Code is readable and comprehensible. Setup instructions and run instructions are provided.
4. Tests are showcased (_no need to cover everything_).
5. Supporting notes on taken decisions and further clarifications are welcome.
