# Saved Queries

RestQL provides two forms to run queries: ad-hoc and saved.

Running an ad-hoc query is simply sending the query as a string to the restQL runtime server, for example:

```bash
curl -d "from people" -H "Content-Type: text/plain" http://localhost:9000/run-query
```

Although it provides flexibility of building the query in the client, giving it the ability to manipulate the query in ways that restQL does not support or to debug new queries, it is not the recommended way to run queries in a production environment, because of the overhead added by the parsing step.

Saved queries are the alternative which deliveries better performance, while also improving debugging. A saved query is just a query that is storage with at least one of the two strategy supported by restQL, the database or the configuration file. Every saved query is defined by three identifiers:

- Namespace: allow grouping logically related queries, like for teams or applications, like `hero-catalog`.
- Name: a descriptive name for the query, like `fetch-dc-heros`.
- Revision: an integer version for the query, restQL favors an immutable approach to updating a query and is recommend that every change to the query results in a new revision.

Based on these identifiers, you can run a saved query by calling the restQL runtime server:

```bash
curl http://localhost:9000/run-query/hero-catalog/fetch-dc-heros/1
```

This request will execute the version `1` of the `fetch-dc-heros` query in the `hero-catalog` namespace.

## Configuration file

You can store queries in the configuration file, for example:

```yaml
queries:
  myNamespace:
    myQuery:
      - |
        from hero
    otherQuery:
      - |
        from sidekick
```

You can define multiples namespaces and multiples queries in each namespace. Each query has a list of revisions, hence if you want call the endpoint `/run-query/myNamespace/myQuery/1` you would execute the query `from hero`, which is the first element in the array associated with the query in the YAML.

Note that this query will only be used if one with the same namespace, name and revision is not present in the database.

## Database

You can add support to store queries to a database trough a Database Plugin. You can learn more about it in the [Plugins documentation](/restql/plugins.md).

In a production environment we recommend the use of the [restQL Manager](/restql/manager.md) to manage the queries in a database rather than manually. The restQL Manager automatically enforces the queries' immutability, creating a new revision every time an existing query is updated.