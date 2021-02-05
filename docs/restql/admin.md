# Administrative API

RestQL have an optional Database plugin that allows it to store resource mappings and queries on an external storage.

In order to support the administration of these configurations and better support the [RestQL Manager](/restql/manager.md).

The reading endpoints expose queries and mappings stored on the database, config file and environment. However, the writing endpoints only allow operations on entities stored on the database.

Any operation on those endpoints is protected by an Authorization code, configured via the `http.server.admin.authorizationCode` field or the `RESTQL_ADMIN_AUTHORIZATION_CODE` environment variable, which should be sent on the `Authorization` header as a bearer token like `Bearer <auth-code>`.

### REST endpoints

All the endpoints described would be placed under `/admin/` endpoint, i.e. `GET /tenant` means `GET /admin/tenant`.

### `GET /tenant`
List all tenants available

**Return**:
```json
{
  "tenants": [
    "dc",
    "marvel",
    "vertigo"
  ]
}
```

### `GET /tenant/:name/mapping`
List mappings of name to URL under the given tenant `:tenant`.

Optionally the client can send an `source` query parameter to filter mappings by its storage.

**Return**:
```json
{
  "tenant": "marvel",
  "mappings": {
    "hero": {
      "url": "http://marvel.api/hero/:id",
      "source": "database"
    },
    "weapons": {
      "url": "http://marvel.api/weapons",
      "source": "database"
    }
  }
}
```

### `POST  /tenant/:name/mapping/:name`
Update the URL associated with the mapping `:name` under the tenant `:tenant`

**Body**: 
```json
{ 
  "url": "http://some.api/resource/:id",
}
```

### `GET /namespace`
List all query namespaces available

**Return**:
```json
{
  "namespaces": [
    "cardgame",
    "moba"
  ]
}
```

### `GET /namespace/:namespace/query`
List all queries under the given `:namespace`

Optionally the client can send an `source` query parameter to filter query revisions by its storage.

**Return**:
```json
{
  "namespace": "cardgame",
  "queries": [
    {
      "name": "my-query",
      "revisions": [
          { "text": "from hero" },
          { "text": "from hero as h" }
      ]
    }
  ]
}
```

### `GET /namespace/:namespace/query/:name`
Fetch all query revisions of the `:query` under the namespace `:namespace`

Optionally the client can send an `source` query parameter to filter query revisions by its storage.

**Return**:
```json
{
  "namespace": "my-namespace",
  "name": "my-query",
  "revisions": [
      { "text": "from hero" },
      { "text": "from hero as h" }
  ]
}
```

### `GET /namespace/:namespace/query/:name/revision/:revision`
Fetch query revision `:revision` of the `:query` under the namespace `:namespace`

**Return**:
```json
{
  "namespace": "my-namespace",
  "name": "my-query",
  "revision": { 
    "text": "from hero" 
  }
}
```

### `POST /namespace/:namespace/query/:name`
Create a new revision of query `:query` under namespace `:namespace`. If the query does not exist, create it.

**Body**:
```json
{
  "text": "from hero as h" 
}
```
