<img name="logo" src="./assets/images/logo_text.svg?sanitize=true">

<br/><br/>

# What is restQL?

**restQL** is a data platform that makes it easy to fetch information from multiple services in the most efficient manner.

It is powered by a query language designed to take advantage of the REST architectural pattern, and a performant runtime built as an HTTP server.

An example restQL query can retrieve information from a `search` resource and use its result to fetch data in other resource, like `hero`. In this case, restQL will only execute the request to `hero` and to `enemy` when the `search` request returns, i.e. sequentially. However, since `hero` and `enemy` don't depend on each other, when their dependency returns, both requests will be executed, i.e. in parallel.

```
from search
    with
        role = "hero"

from hero as heroList
    with
        name = search.name

from enemy as enemiList
    with
        hero = search.name
```

## Next steps

1. Setup with our [getting started](/restql/getting-started.md)
2. Learn restQL [query language](/restql/query-language.md),
3. Get involved :) We're looking for contributors, if you're interested open a Pull Request at our [GitHub Project](https://github.com/b2wdigital/restQL-golang).

## Help and community

If you need help you can reach the community on Telegram:

- https://t.me/restQL

## Useful Links

- [Tackling microservice query complexity](https://medium.com/b2w-engineering/restql-tackling-microservice-query-complexity-27def5d09b40): Project motivation and history
- [@restQL](https://t.me/restQL): restQL Telegram Group
- [restql.b2w.io](http://restql.b2w.io): Project home page,
- [game.b2w.io](http://game.restql.b2w.io/): A game developed to teach the basics of restQL language,
- [restQL-golang](https://github.com/b2wdigital/restQL-golang): Main platform implementation, including the language specification and server runtime.
- [restQL-clojure](https://github.com/b2wdigital/restQL-clojure): _Deprecated implementation, allowed to embed the restQL runtime in your Clojure application_
- [restQL-java](https://github.com/b2wdigital/restQL-java): _Deprecated implementation, allowed to embed the restQL runtime in your Java application_
- [restQL-manager](https://github.com/b2wdigital/restQL-manager): To manage saved queries and resources endpoints (requires a MongoDB instance).
- [Wiki](https://github.com/b2wdigital/restQL-golang/wiki/RestQL-Query-Language): Git Hub documentation.

Who're talking about restQL:

- [infoQ: restQL, a Microservices Query Language, Released on GitHub](https://www.infoq.com/news/2018/01/restql-released)
- [infoQ: 微服务查询语言 restQL 已在 GitHub 上发布](http://www.infoq.com/cn/news/2018/01/restql-released)
- [OSDN Mag: マイクロサービスクエリ言語「restQL 2.3」公開](https://mag.osdn.jp/18/01/12/160000)

## Releases

You can find our latest releases at our [release page](https://github.com/b2wdigital/restQL-golang/releases).

## License

Copyright © 2016-2020 B2W Digital

Distributed under the MIT License.
