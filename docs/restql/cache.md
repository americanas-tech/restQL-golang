# Cache Control

One of restQL cornerstones is to keep HTTP semantics whenever that's possible. HTTP's headers play a key role in current HTTP tools and servers, worth mentioning the _Cache-Control_ header.

Therefore, restQL support two cache directives: `max-age` and `s-max-age`, both can be applied at query level with `use` or at the statement level.

> All the behaviour described below works the same way with the `s-max-age` directive.

Cache-control is a header returned in an HTTP call that tells the intermediate proxies and the end client on how to handle the cache of the returned content. For example: a _Cache-Control_ header with _max-age=60_  tells that the client can safely cache the request for 1 minute.

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
