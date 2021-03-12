# Running Queries

RestQL provides two forms to run queries: ad-hoc and saved.

## Ad-hoc Queries

Running an ad-hoc query is simply sending the query as a string to the restQL runtime server, for example:

```bash
curl -d "from people" -H "Content-Type: text/plain" http://localhost:9000/run-query?tenant=MYTENANT
```

Although it provides flexibility of building the query in the client, giving it the ability to manipulate the query in ways that restQL does not support or to debug new queries, it is not the recommended way to run queries in a production environment, because of the overhead added by the parsing step.

## Saved Queries

Saved queries are the alternative which deliveries better performance, while also improving debugging. A saved query is just a query that is storage with at least one of the two strategy supported by restQL:

1. Configuration file:

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

2. Database:

You can add support to store queries to a database trough a Database Plugin. You can learn more about it in the [Plugins documentation](/restql/plugins.md).

In a production environment we recommend the use of the [restQL Manager](/restql/manager.md) to manage the queries in a database rather than manually. The restQL Manager automatically enforces the queries' immutability, creating a new revision every time an existing query is updated.

---

Every saved query is defined by three identifiers:

- Namespace: allow grouping logically related queries, like for teams or applications, like `hero-catalog`.
- Name: a descriptive name for the query, like `fetch-dc-heros`.
- Revision: an integer version for the query, restQL favors an immutable approach to updating a query and is recommend that every change to the query results in a new revision.

Based on these identifiers, you can run a saved query by calling the restQL runtime server:

```bash
curl http://localhost:9000/run-query/hero-catalog/fetch-dc-heros/1?tenant=MYTENANT
```

This request will execute the version `1` of the `fetch-dc-heros` query in the `hero-catalog` namespace.

> An important aspect that is present in both forms of execution is the `tenant` query parameter. Tenants are the way restQL organizes mappings, for example `staging` vs `production` or `marvel` vs `dc`. If the restQL instance has a `RESTQL_TENANT` environment variable, this parameter is not used. However, if it is not set, then the client must always provide it.

## RestQL Traits

### Global Status Code

Since restQL access multiple APIs on the same query and it is designed to leverage the REST pattern, the status code returned for the HTTP request to run the query is calculated based on the APIs responses.

For all statements without an `ignore-errors` clause, restQL takes the greater status code and use it as the global status code.

With this, if one of the APIs return `500 Internal Server Error`, this will be the status code returned by restQL.

### Forward Headers

By default, the headers send to restQL on the run query request are forward to all APIs on the query. This simply use cases like tracing headers and authorization and avoids query cluttering, since you do not need to specify every header you wish to send.

### Response Headers

In some cases the client needs to extract information from the headers returned by one the APIs called on the query, for example wehn creating a resource and the API returning its unique id as the `Location` header.

For this use cases restQL automatically return the response headers of each API prefixed with the resource name. For example:

```restql
to hero
  with
    name = "Super Shock"
```

This would have an the following response

```
< HTTP/1.1 422 Unprocessable Entity
< Content-Type: application/json; charset=utf-8
< Content-Length: 281
< Server: restql
< Date: Fri, 12 Mar 2021 17:45:32 GMT
< Connection: keep-alive
< hero-Location: ad1mshcbah1i01hnf
<
{"hero: {"details": {"success": true, "status": 201, "metadata": {}}, "result": {}}}
```

### Cache Control

One of restQL cornerstones is to keep HTTP semantics whenever that's possible. HTTP's headers play a key role in current HTTP tools and servers, worth mentioning the _Cache-Control_ header.

Therefore, restQL support two cache directives: `max-age` and `s-max-age`, both can be applied at query level with `use` or at the statement level.

> All the behaviour described below works the same way with the `s-max-age` directive.

Cache-control is a header returned in an HTTP call that tells the intermediate proxies and the end client on how to handle the cache of the returned content. For example: a _Cache-Control_ header with _max-age=60_ tells that the client can safely cache the request for 1 minute.

Consider a query fetching two resources:

```
from hero

from sidekick
```

What happens if **hero** returns _max-age=60_, and **sidekick** _max-age=30_?

To avoid stale data on the client restQL will return the least common configuration, i.e. the lowest value among them.

If one resource returns _no-cache_, that should have precedence over the max-age headers and _no-cache_ should be the return of the query.

Besides the values returned by the upstream resources, you can also define cache control values in the query. They will be taken into account together with the return values, and the lowest among them all will be returned.

For example, consider the following query:

```
from hero
    max-age 40
```

If **hero** returns _max-age=60_, restQL will assume _max-age=40_ for the **hero** resource (the value defined in the query is lower) and since this is the only resource in the query, this will be used in the return header.

But if you have a query with more resources like below, with **hero** returning _max-age=60_, and **sidekick** returning _max-age=30_:

```
from hero
    max-age 40

from sidekick
    max-age 50
```

Then restQL will assume _max-age=40_ for the **hero** resource (the value defined in the query is lower) and will assume _max-age=30_ for the **sidekick** resource (the returned value is lower).

Hence, restQL will choose the lowest value among the assumed values for each resource, in this case **hero** has _max-age=40_ and **sidekick** has _max-age=30_, so the _Cache-Control_ returned will be _max-age=30_.

You can also define a query global _Cache Control_ directive, for example:

```
use max-age 10

from hero
    max-age 40

from sidekick
    max-age 50
```

Given the same result by the resources (**hero** returning _max-age=60_ and **sidekick** returning _max-age=30_), the _Cache-Control_ returned would be _max-age=30_, but once the global _Cache Control_ is determined restQL will compare it with the query global cache directives and return the lowest.

Hence, the _Cache-Control_ returned will be _max-age=10_.
