# Setting Up

## Configuring a REST API resource

To add a resource to RestQL we simply need to edit `restql.yml`, the file you created a few moments ago. Let's add Elon Musk's space business' API:

```yaml
mappings:
  launches: "https://api.spacexdata.com/v3/launches"
```

Any value assigned to that a key in `mappings` should point to a valid API.

## Running it

To get a restQL instance up you just need to run the following command from the project root:
```shell script
$ RESTQL_CONFIG=/path/to/restql.yml make dev
```

Nice! Now `restQL-http` is being served at `localhost:9000`!