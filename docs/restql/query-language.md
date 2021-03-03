# REST Query Language

## Query Syntax

The clause order matters when making restQL queries. The following is a full reference to the query syntax, available clauses and order.

```restql
[ [ use modifier value ] ]

METHOD resource-name [as some-alias] [in some-resource]
  [ headers HEADERS ]
  [ timeout INTEGER_VALUE ]
  [ with WITH_CLAUSES ]
  [ [only FILTERS] OR [hidden] ]
  [ [ignore-errors] ]
```

## Starting a query

The only required clauses in a query are the method and the resource, which will be target. RestQL supports the most popular HTTP methods:

- **from**: HTTP GET
- **to**: HTTP POST
- **into**: HTTP PUT
- **update**: HTTP PATCH
- **delete**: HTTP DELETE

For example, the query `to hero` maps to the following call:

```shell
POST http://some.api/hero/
```

Usually, beyond method and resource, a statement has an alias. It is an optional way to define a custom name reference for the result of that statement. For example, `hero` is the resource being queried and `batman` is the alias which can be used to reference the statement result. If no alias is used, the resource name is then used as a reference.

```restql
from hero as batman
```

## Sending parameters

You can define parameters to be send in the request using the `with` clause, which simply contains a list of key/value pairs.

```restql
from hero as batman
    with
        name = "Batman"
```

You can put more than one item in the same line separating them by a comma or in different lines, so the two examples below are equivalent:

```restql
from hero as batman
    with
        name = "Batman", universe = "DC"
```

```restql
from hero as batman
    with
        name = "Batman"
        universe = "DC"
```

In REST, we can send parameters using path parameters, query parameters, and a request body. RestQL support them all through the `with` clause.

### Path parameters

A path parameter is a variable part of the target URL, hence it is necessary to define in the resource mapping where the path parameters are present. For example, if we define a resource `hero` with URL `http://some.api/hero/:id`, we can send the `id` path parameter using the following query:

```restql
from hero as batman
    with
        id = 1
```

This will map to the following request:

```shell
GET http://some.api/hero/1
```

You can define as many path parameters as you like. Every path element prefixed with a colon (`:`) in the URL will be parsed as a path parameter.

If you don't define a path parameter in the `with` clause, it will be skipped when building the final URL, for example, the following query:

```restql
from hero as batman
```

Will map to the following request:

```shell
GET http://some.api/hero/
```

This is especially useful as the same resource can target an entity by its ID or list all the entities, which is a common design in REST APIs.

### Query parameters

When using the `from` method every parameter in the `with` clause will be mapped to a query parameter, for example:

```restql
from hero as batman
    with
        name = "Batman"
        universe = "DC"
```

Will map to the following request:

```shell
GET http://some.api/hero?name=Batman&universe=DC
```

You can also specify query parameters in resource mappings like the path parameters above. For example, if we define a resource `hero` with URL `http://some.api/hero/:id?:universe`, we can send the `id` path parameter and the `universe` query parameter using the following query:

```restql
to hero as batman
    with
        id = 1
        universe = "DC"
```

Will map to the following request:

```shell
POST http://some.api/hero/1?universe=DC
```

This is specially useful for sending obligatory query parameters in statements with `to`, `into`, `patch` and `delete` methods. 

## Parameters types

Inside the `with` clause, you must specify a list of key/value pairs, the values can be one of the following types:

- A string, which must be enclosed in **double quotes**
- A number, that can have a floating-point. Scientific notation is not supported.
- A boolean: either `true` or `false` in **lower case**
- The null value: `null`
- The variable value. It is a reference value which will be resolved based on data send to restQL by you. To learn more about it read the "Using Variables" section.
- A list of values. Lists are enclosed in **square brackets** and separated by either **commas** or **newlines**
- A key/value structure. Structures must be enclosed in **curly braces**, with each pair separated from each other with a **comma** or a **newline**. The key and the value must be separated with a **colon**, similar to a `json` object.
- A chained value. A chain is a reference value which will be resolved using the results of another statement. You start the chain with name of the bound variable that reference a statement followed by a **dot** specifying the field. You can keep adding dots to go arbitrarily deep within the statement result. The specified field (or path) is resolved using the body returned by the statement and, if not present, is resolved using the headers returned by the statement.

Here is an example containing all the types mentioned above:

```restql
from hero as protagonist
    with
        name = "Super Duper"          // String type
        level = 15                    // Number Type
        honored = true                // Boolean Type
        lastDefeat = null             // Null Value
        weapons = $weapons
        using = ["sword", "shield"]   // List Type
        stats = {health: 100,
                 magic: 100}          // KeyValue Type

from hero as sidekick
    with
        id = protagonist.sidekick.id  // Chaining Type
```

### Body

When using the methods `to`, `into` or `update` every parameter in the `with` clause will be mapped to the request body, for example:

```restql
to hero as batman
    with
        name = "Batman"
        universe = "DC"
```

Will map to the following request:

```shell
POST http://some.api/hero
BODY {"name": "Batman", "universe": "DC"}
```

You can also map a specific `with` parameter as body, for example:

```restql
to hero
    with
        heroes = [
            {name: "Batman"},
            {name: "Superman"},
        ] -> as-body
```

Will map to the following request:

```shell
POST http://some.api/hero
BODY [{"name": "Batman"}, {"name": "Superman"}]
```

This is specially useful in case you want to build an array body using values from previous statements, with a chained value. 

Besides that the `with` clause also support a dynamic body, that allows the client to send a JSON that can contain anything and will be used as the request body, for example:

```restql
update hero as batman
    with
        $newHero
        id = 1
```

If you call restQL with `GET http://your-ip/run-query?newHero={"name":"Batman","universe":"DC"}`, the `newHero` variable will resolve to the JSON `{"name": "Batman", "universe": "DC"}`. To learn more about variable read the section "Using Variables". Then it will map to the following request:

```shell
PATCH http://some.api/hero/1
BODY {"name": "Batman", "universe": "DC"}
```

Any key/value items declared in the `with` clause when using the dynamic body will only be used to supply path parameters.

## Specifying Headers

Before the `with` clause you can add a `headers` clause to define the headers you want to send within that statement. The headers are a list of key/value pairs, like the `with` clause items, but the values must be strings or variables (see below).

```restql
// example using headers

from hero
headers
    Authorization = "Basic user:pass"
    Accept        = "application/json"
with
    id = 1
```

It is important to state that headers present in the query will substitute any request headers with the same name, therefore in the above example, even if the request already has an `Authorization` header, it will be replaced by `"Basic user:pass"`.

## Timeout Control

A specific statement has the default timeout defined in the configurations, which is usually a high value to cover most cases. To change the timeout value for any statement, use the `timeout` clause, which accepts an integer value or a variable (see below), representing the **milliseconds** to wait before the request times out.

The `timeout` clause appears **before** the `with` clause.

```restql
from hero
headers
    Authorization = "Basic user:pass"
    Accept        = "application/json"
timeout 200
with
    id = 1
```

## Using Variables

Alongside directly typing a value or using a chained value, it is possible to define variable that will have their values resolved based on data send to restQL.

Variables can be used inside a statement in the `headers`, `timeout`, `max-age`, `s-max-age` or `with` clauses.

For example, the query below will have its variables resolved using one of the following strategies:

1. Body resolution: if you executed a `POST /run-query`, then the fields in body sent will be used to resolve the variables. If some variable is not found in the body, it will use subsequent strategy.
2. Query Parameter resolution: in either case of a `POST /run-query` or a `GET /run-query`, the query parameters sent will be used to resolve the variables. If some variable is not found in the query parameters, it will use subsequent strategy.
3. Headers resolution: in either case of a `POST /run-query` or a `GET /run-query`, the headers sent will be used to resolve the variables. If some variable is not found in the headers, the query will fail or skip the parameter, depending on where the variable was used.

```restql
from hero
    headers
        Auth = $authentication
    with
        name = $heroName
        level = $heroLevel
```

## Multiplexing

Whenever restQL finds a List value in a `with` parameter, it will perform an **expansion**, which means it will make one request for each item in the list. Suppose we want to fetch the `superheroes` with ids 1, 2 and 3:

```restql
// this query will make THREE requests to the superheroes resource

from superheroes as party
    with
        id = [1, 2, 3]
```

In this case, restQL will perform the following HTTP calls:

`GET http://some.api/superhero?id=1`

`GET http://some.api/superhero?id=2`

`GET http://some.api/superhero?id=3`

If this behaviour is not what you want, but rather pass all values in a single request, you can disable the multiplexing of any parameter by using the **apply** operator `->` with the `no-multiplex` function.

```restql
// now only ONE request will be performed
from superheroes as fused
    with
        id = [1, 2, 3] -> no-multiplex
```

Using `no-multiplex`, restQL will perform just **one** HTTP call, as follows:

`GET http://some.api/superhero?id=1&id=2&id=3`

## Object explosion

Whenever restQL finds an Object value with a list field in a `with` parameter, it will perform an **explosion**, which means it will turn the object into a list of objects for each list value.

```restql
to superheroes
    with
        profiles = {
          cities: ["Gotham", "Metropoles"]
        }
```

In this case, restQL will build the following body:

```json
{
  "profiles": [
    {
      "cities": "Gotham"
    },
    {
      "cities": "Metropoles"
    }
  ]
}
```

If this behaviour is not what you want, but rather send the object as is, you can disable the explosion of any parameter by using the **apply** operator `->` with the `no-explode` function.

```restql
to superheroes
    with
        profiles = {
          cities: ["Gotham", "Metropoles"]
        } -> no-explode
```

Using `no-explode`, restQL will send the following body:

```json
{
  "profiles": {
    "cities": ["Gotham", "Metropoles"]
  }
}
```

## Selecting the returned fields

When the response of a given statement is bloated you may want to filter the fields in order to reduce query payload. You can do this by adding an `only` clause to the end of a statement, simply listing the fields you want:

```restql
from hero
    with
        id = 1
    only
        name
        items
        skills.id
        skills.name
```

For the sub-elements, like `skills.id` and `skills.name` above, the fields `id` and `name` will be nested in a `skills` top-level field.
There is also a special filter `*` which will simply return all the fields. Normally it is redundant but there are special cases where it is useful and you can see in the Functions section (see below).

You also have to option to suppress a statement in the query response. It is usually useful for statements that are only used as an intermediate step to build a parameter to another statement.

```restql
from hero
    with
        name = "Restman"
    hidden

from sidekick
    with
        hero = hero.id
```

## Functions

Sometimes you may need to perform computations a value before sending or returning it. To address this need restQL provides functions, that can be used by specifying its name after a `->` operator. RestQL ships with three built-in functions:

- **base64**: stringify and them hashes the value using a base 64 algorithms.
- **json**: stringify the value using the JSON syntax. For any key/value structure in a `from` statement it is used by default.
- **flatten**: take a list value, usually nested, and return a plain list.
- **matches**: conditionally filter the result of a statement by a regex. If the field contains a string, it only returns the field if it matches the regex. If the field contains a list, it applies the matching to each element, returning a filtered list with the successful matches.

```restql
from hero
    with
        stats = {health: 100,
                 magic: 100} -> base64
    only
        nicknames -> matches("^Super")
        *
```

In this case we use two functions. First, we encode the key/value structure as a base64 hash before sending it to the API. Then, we combine the `matches` function with the all filter selector `*`, this has the effect of returning all fields in the statement response, filtering only the `nickname` field by the specified regex.

## Aggregating result in another statement

RestQL provides an aggregation clause that allows you to easily append a statement result into another. To achieve this use the `in` clause, for example:

```restql
from hero
    with
        name = "Restman"

from sidekick in hero.sidekick
    with
        id = hero.sidekickId
```

The query above will be aggregated as bellow:

```json
{
    "hero": {
        "details": {...},
        "result": {
            "id": 1,
            "name": "Restman",
            "sidekickId": 10,
            "sidekick": {
                "id": 10,
                "name": "Super"
            }
        }
    },
    "sidekick": {
        "details": {...}
    }
}
```

## Ignoring error of a statement

By default, restQL returns the highest HTTP status code returned by the statements. If you'd like restQL to ignore a given statement when calculating the return status code you can use ignore-error modifier on that statement.

This is useful when querying critical and non-critical resources in the same query. As an example, we might want to get all available products and their ratings, but if we get an error from rating we want to ignore it and show the products anyway. For cases like these, we can use the `ignore-errors` expression, as follows:

```restql
from products as product

from ratings
  with
    productId = product.id
  ignore-errors
```

The query above will return a success HTTP status code even when the ratings resources returns an error.

### Explicit dependency

There are two types of statement dependency on restQL: implicit and explicit.

Implicit dependency is inferred by restQL through chained parameters. For example, if we have a query like
```restql
from hero

from sidekick
    with
        id = hero.sidekick.id
```
This declares that the `sidekick` statement has an dependency on the `hero` statement, since to build the request to the `sidekick` API restQL needs data presented on `hero` response.

However, there are cases where we would like statements to have a dependency link, but we don't need to pipe any data from one to another. For example, when we create a new resource and would like to refresh the list of all resource.

We can achieve this by setting an explicit dependency through the `depends-on` keyword, like this:
```restql
to hero
    with
        name = "Super restQL"
        
from hero as allHeroes
    depends-on hero
```

### Cache Control

By default, restQL returns the lowest cache-control value among all statements. You can add a maximum age for the cache control returned by a statement, for example:

```restql

from hero
    max-age 100

from sidekick
    with
        hero = hero.id
```

RestQL will use the lowest value between `max-age 100` and the cache-control returned by the `hero` resource as the cache-control for that statement. Then, restQL will use lowest cache-control value among all statements to calculate the final header.

You can also set a maximum age for the cache-control at the query level, for example:

```restql
use max-age 600

from hero

from sidekick
    with
        hero = hero.id
```

If `max-age 600` is lower than the cache-control for each statement, then it will be used as the final header. But if one of the statements has a cache-control lower than the query level one, this statement cache-control will be used.
