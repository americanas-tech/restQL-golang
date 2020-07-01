# Getting Started

## Running restQL

If you have Go installed then you can run `go get -u github.com/b2wdigital/restQL-golang` and have restQL in your `$GOBIN` directory. However, if you prefer you can download the binary for you OS in the [release page](https://github.com/b2wdigital/restQL-golang/releases).

Then you can start the server runtime with `RESTQL_PORT=9000 RESTQL_MAPPING_PEOPLE=https://swapi.co/api/people restql`

The provided environment variable `RESTQL_MAPPING_PEOPLE` is what is called a mapping, it associates a name with a REST resource URL, in this case the resource will be mapped to the name `people`.

With this we can run an ad-hoc query, i.e. a query that is send in the request to the restQL rather than stored in configuration or database (for more information, please read the [configuration page](/restql/config.md)). To run it you should send a POST request to restQL with the query string as our body content, for example:

```bash
curl -d "from people" -H "Content-Type: text/plain" http://localhost:9000/run-query
```

This query will execute a GET request on the URL fetching all people in the API and return it to you.

```json
{
  "people": {
    "details": {
      "success": true,
      "status": 200,
      "metadata": {}
    },
    "result": {
      "count": 82,
      "next": "http://swapi.dev/api/people/?page=2",
      "previous": null,
      "results": [
        {
          "name": "Luke Skywalker",
          "height": "172",
          "mass": "77",
          "hair_color": "blond",
          "skin_color": "fair",
          "eye_color": "blue",
          "birth_year": "19BBY",
          "gender": "male"
        },
        {
          "name": "C-3PO",
          "height": "167",
          "mass": "75",
          "hair_color": "n/a",
          "skin_color": "gold",
          "eye_color": "yellow",
          "birth_year": "112BBY",
          "gender": "n/a"
        }
      ]
    }
  }
}
```

This is the simplest query that can be made with restQL, however the language that supports this platform is much more powerful. Follow the next steps to learn more.

## Next steps

1. Learn restQL [query language](/restql/query-language.md),
2. Learn about the [manager and saved queries](/restql/manager),
3. Get involved and [contribute ツ](/restql/how-to-contribute)
