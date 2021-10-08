# Configuration

RestQL can be configured either via `Environment Variables` or through a `Config File`. It will search the file first at the path defined in `RESTQL_CONFIG` environment variable, then at `./restql.yml` and finally at `${HOME}/restql.yml`.

Configuration options follows the precedence `Environment > Config File > Default`.

For configuration details of plugins, like database plugin or others, please refer to their documentation.

> Any times this documentation uses "duration string" it refers to the [Go time duration syntax](https://golang.org/pkg/time/#ParseDuration).

## Tenants

If you are working with a database for storing mappings, then you will be able to scope then in tenants. The tenant that a query should use to resolve its statements can be defined in two ways:

1. Through a `tenant` query parameter, which allow for a same restQL deployment to run each query with a possibility different tenant.
2. Through a `RESTQL_TENANT` environment variable, which will lock the tenant allowed to be used by that deployment, ignoring the `tenant` query parameter if present.

When using tenants the second approach is recommended if you aim to provide isolation between the tenants and guarantee that a tenant cannot produce load in the APIs of other tenants.

## HTTP layer

**Forward prefix**: you can customize restQL to proxy query parameters with the given prefix to the APIs, it is useful to send context query parameters. To set it, use the `http.forwardPrefix` field or the `RESTQL_FORWARD_PREFIX` environment variable, both accept a string.

**Query timeout**: you can define the default maximum time for the query to be executed, that is, the maximum time spent calling the APIs (not including database and parsing latency), if a timeout is defined in the query with `use timeout = <timeout>`, this timeout will be ignored. To set it, use the `RESTQL_QUERY_GLOBAL_TIMEOUT` environment variable, both accept duration string, with a default of 30 seconds.

**Resource timeout**: you can define the default maximum time spent waiting for an API to response, if a timeout is defined for in the query statement for that API, this timeout will be ignored. To set it, use the `RESTQL_QUERY_RESOURCE_TIMEOUT` environment variable, both accept duration string, with a default of 5 seconds.

### Profiling

You can use the `pprof` tool to investigate restQL performance. To enable it set `RESTQL_ENABLE_PPROF` environment variable to `true`, which will expose the basic endpoints for profiling (cpu, heap, threadcreate and goroutine). Setting the variable `RESTQL_ENABLE_FULL_PPROF` will also enable the profiling endpoints for block and mutexes. _Note that enabling all the profiling endpoints can result in serious performance degradation_.

### HTTP Server

**HTTP Ports**: You can customize the ports where the restQL API, health and profiling will run.

- API port: set through `RESTQL_PORT` environment variable.
- Health port: set through `RESTQL_HEALTH_PORT` environment variable.
- Profiler port: set through `RESTQL_PPROF_PORT` environment variable.

**Enable Administrative API**: restQL exposes a set of endpoints to configure queries and mappings stored on the database. One can enable it through the `http.server.admin.enable` field or the `RESTQL_ADMIN_ENABLE` environment variable. To find more about it go to [Administrative API](/restql/admin.md).

**Graceful shutdown**: when restQL receives a `SIGTERM` signal it starts the shutdown, avoiding accepting new requests and waiting for the ongoing ones to finish before exiting. You can define a timeout for this process using `http.server.gracefulShutdownTimeout` field in the YAML configuration, after which restQL will break all running requests and exit.

**Read timeout**: you can specify the maximum time taken to read the client request to the restQL API through the `http.server.readTimeout` field.

**Middlewares**: currently restQL support 3 built-in middlewares, setting any of the fields automatically enable the given middleware.

> From version 6.0.0 all middlewares have an `enable` field that must be set to `true` in order for them to be activated.

- Request ID: this middleware generates a unique id for each request restQL API receives. The `http.server.middlewares.requestId.header` field define the header name use to return the generated id. The `http.server.middlewares.requestId.strategy` defines how the id will be generated and can be either `base64` or `uuid`.
- Timeout: this middleware limits the maximum time any request can take. The `http.server.middlewares.timeout.duration` field accept a time duration value.
- Request Cancellation: this middleware stops query execution when the client drops the connection. This improves fault response as it avoids unnecessary computation and reduces traffic on downstream APIs. You can also manage the connection watching interval with the field `http.server.middlewares.requestCancellation.watchingInterval`, which accepts a duration string.
- CORS: Cross-Origin Resource Sharing is a specification that enables truly open access across domain-boundaries.
  You can configure your own CORS headers either via the configuration file:
  ```yaml
  http:
    server:
      middlewares:
        cors:
          allowOrigin: "example.com, hero.api"
          allowMethods: "GET, POST"
          allowHeaders: "X-TID, X-Custom"
          allowCredentials: false
          exposeHeaders: "X-TID"
          maxAge: 10 # seconds, as per specification
  ```
  Or via environment variables:
  ```shell script
  RESTQL_CORS_ALLOW_ORIGIN=${allowed_custom_origin}
  RESTQL_CORS_ALLOW_METHODS=${allowed_custom_methods}
  RESTQL_CORS_ALLOW_HEADERS=${allowed_custom_headers}
  RESTQL_CORS_EXPOSE_HEADERS=${allowed_custom_expose_headers}
  RESTQL_CORS_ALLOW_CREDENTIALS=${allowed_credentials}
  RESTQL_CORS_MAX_AGE=${allowed_max_age}
  ```

### Http Client

RestQL primary feature is performing optimized HTTP calls, but since each environment has different characteristics like workload and latency, it is important that you tune the parameters for the internal HTTP client in order to achieve the best performance. You can set these parameters throught the configuration file.

- `http.client.connectionTimeout`: limits the time taken to establish a TCP connection with a host.
- `http.client.maxIdleConnectionDuration`: set the time a connection will be kept open in idle state, after it the connection will be closed. It accepts a duration string.
- `http.client.maxConnectionsPerHost`: limits the size of the connection pool for each host.
- `http.client.dnsRefreshInterval`: defines the time a DNS query result will be cached.

#### Concurrency

RestQL provides configuration parameters to limit the workload that it will accept.

With these guardrails in place, once these thresholds are reach restQL will deny service with a _507 Insufficient Storage_ status code. This was chosen because in practice these limits exists to avoid unbound growth of goroutines and memory usage that would lead to application restarts.

It is best practice to always run production restQL deployments with these parameters set. The absence of those or defining it to `0` will disable the mechanisms that avoid linear resource consumption when workload increase.

**Maximum concurrent queries**: this is the first limiter and will accept or reject a query entirely. It can be defined with the configuration field `http.client.maxConcurrentQueries` or through the environment variable `RESTQL_MAX_CONCURRENT_QUERIES`. This parameter should be more strict because it has more impact on reducing resource consumption and fails the execution faster.

**Maximum concurrent goroutines**: this is the second limiter and will accept or reject a goroutine call when running a query. It can be defined with the configuration field `http.client.maxConcurrentGoroutines` or through the environment variable `RESTQL_MAX_CONCURRENT_GOROUTINES`. This parameter should be more loose since the numbers of goroutines can vary drastically depending on runtime data that define the number of multiplexed calls that should be made. Also, if during a query execution one goroutine fails to be accepted, the entire query will be discarted, and a _507 Insufficient Storage_ status code will be returned.

> P.S.: The goroutine limiter only applies to goroutines used to process and dispatch HTTP requests to upstream APIs. If measuring the total number of goroutines in your deployment, it will be greater than the maximum concurrent goroutine, since it does not impact the usage of goroutines to accept new connections and other tasks.

_Deprecated on v4.2.0:_

- `http.client.maxRequestTimeout`: although every the timeout for calling a resource can be defined by the client in the query you can set a upper limit to request time, for example, if you set it to `2s` even though a query specifies a timeout of `10s` restQL will drop the request when it reachs its maximum timeout. It accepts a duration string.
- `http.client.maxIdleConnections`: limits the size of the global idle connection pool.
- `http.client.maxIdleConnectionsPerHost`: limits the size of the idle connection pool for each host.

## Caching

RestQL uses cache to avoid excessive database calls and grammar parsing. The cache used for the parser and for the fetching queries from databases uses a simple LRU strategy.

For fetching mappings from the database restQL uses a stale-cache strategy, which runs an update task in background when the TTL for an entry expire and only replace the cached value if the fetching is successful. This allows restQL to stay updated but not break if the database goes offline.

You can customize each cache maximum size and, for the mappings cache, other parameters. You can also disable all caching using `cache.disable: true` or `RESTQL_CACHE_DISABLE=true`.

**Queries and Parser**

Both caches have only one parameter, maximum cache size.

To set it for the query cache use the field `cache.query.maxSize` or the `RESTQL_CACHE_QUERY_MAX_SIZE` environment variable, they accept a integer value greater than zero.

And, to set it for the parser cache use the field `cache.parser.maxSize` or the `RESTQL_CACHE_PARSER_MAX_SIZE` environment variable, they accept a integer value greater than zero.

**Mappings**:

This cache has a maximum size, an expiration used for all entries and parameters for the background routine responsible for the update expired entries.

To set the size use the field `cache.mappings.maxSize` or the `RESTQL_CACHE_MAPPINGS_MAX_SIZE` environment variable, they accept an integer value greater than zero.
In order to customize the expiration duration use the field `cache.mappings.expiration` or the `RESTQL_CACHE_MAPPINGS_EXPIRATION` environment variable, they accept a duration string.

The update background routine has two parameters

- Refresh interval: for example if it is set to `30s` then the routine will run every thirty seconds. To set it, use the `cache.mappings.refreshInterval` field or the `RESTQL_CACHE_MAPPINGS_REFRESH_INTERVAL` environment variable, both accept a duration string.
- Refresh Queue Length: when an entry is hit and expired, a task in added to the background update routine queue. Every time the routine run, all tasks in this queue are executed. You can limit the size of this queue, which effectively limits the batch size which the background routine will receive every time it runs and, therefore, limits the time which will be spent in the background routine every time. To set it, use the `cache.mappings.refreshQueueLength` field or the `RESTQL_CACHE_MAPPINGS_REFRESH_QUEUE_LENGTH` environment variable, both accept an integer value.

## Logging

Due to the traffic restQL is designed to handle it takes a conservative approach to logging, placing the most of it in the `DEBUG` level. You can customize this log level and others parameters through the configuration file:

- `logging.enable`: boolean value that can disable all logging.
- `logging.timestamp`: boolean value that indicate with a timestamp field should be added to the log entry.
- `logging.level`: the minimum log level required for a log entry to be output. You can see the list of available levels on the [zerolog documentation](https://github.com/rs/zerolog#leveled-logging).

## Debugging

RestQL supports a debug mode where its response verbosity is increased to include details about the HTTP request made to the upstream API resource. Currently, restQL support debug activation through a query parameter or header. By default, debugging with a query parameter is enabled. You can customize this with the following parameters in the configuration file or environment variables:

- `debugging.queryParam` or `RESTQL_DEBUGGING_QUERY_PARAM`: enables debugging using the `_debug` query param, `true` by default.
- `debugging.header` or `RESTQL_DEBUGGING_HEADER`: enables debugging using the `X-Restql-Debug` header, `false` by default.

## Alternative storage for mappings and queries

To understand others stores besides a database for mappings and queries please refer to [Resource Mappings](/restql/resource-mappings.md) and [Running Queries](/restql/running-queries.md) pages.
