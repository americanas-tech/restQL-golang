<p align="center">
  <a href="http://restql.b2w.io">
    <img width="537px" height="180px" alt="restQL" src="./docs/assets/images/logo_text.png?sanitize=true">
  </a>
</p>

<p align="center">
  restQL-golang allows you to use the <strong>restQL</strong> language to run queries that define how you want to access your REST services, making it easy to fetch information from multiple points in the most efficient manner.
</p>

[![GoDoc](https://godoc.org/github.com/b2wdigital/restQL-golang?status.svg)](https://pkg.go.dev/github.com/b2wdigital/restQL-golang)
[![Go Report Card](https://goreportcard.com/badge/github.com/b2wdigital/restQL-golang)](https://goreportcard.com/report/github.com/b2wdigital/restQL-golang)

# Getting Started

For a throughout explanation of how to run restQL follow [this tutorial](http://docs.restql.b2w.io/#/restql/tutorial/intro).

## Query language
The clause order matters when making restQL queries. The following is a full reference to the query syntax, available clauses and order.

```
[ [ use modifier = value ] ]

METHOD resource-name [as some-alias] [in some-resource]
  [ headers HEADERS ]
  [ timeout INTEGER_VALUE ]
  [ with WITH_CLAUSES ]
  [ [only FILTERS] OR [hidden] ]
  [ [ignore-errors] ]
```
e.g:
```restQL
from search
    with
        role = "hero"

from hero as heroList
    with
        name = search.results.name
```
Learn more about [**restQL** query language](http://docs.restql.b2w.io/#/restql/query-language)

# Links
* [Docs](http://docs.restql.b2w.io)
* [restQL-manager](https://github.com/B2W-BIT/restQL-manager): To manage saved queries and resources endpoints. restQL-manager requires a MongoDB instance.
* [Tackling microservice query complexity](https://medium.com/b2w-engineering/restql-tackling-microservice-query-complexity-27def5d09b40): Project motivation and history

## Reach the community
* [@restQL](https://t.me/restQL): restQL Telegram Group

## Who's talking about restQL

* [infoQ: restQL, a Microservices Query Language, Released on GitHub](https://www.infoq.com/news/2018/01/restql-released)
* [infoQ: 微服务查询语言restQL已在GitHub上发布](http://www.infoq.com/cn/news/2018/01/restql-released)
* [OSDN Mag: マイクロサービスクエリ言語「restQL 2.3」公開](https://mag.osdn.jp/18/01/12/160000)
* [Build API's w/ GraphQL, RestQL or RESTful?](https://www.youtube.com/watch?v=OeUGswoYrvA)

## License

Copyright © 2016-2020 B2W Digital

Distributed under the MIT License.

