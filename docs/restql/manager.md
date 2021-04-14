# restQL Manager

It is a Web UI over the restQL Admin API in order to facilitate daily operations like query creation, resource mappings, debugging and experimentation.

restQL Manager requires an [restQL](https://github.com/b2wdigital/restQL-golang) instance with the Admin API enabled.

## Features

With the restQL Manager you can play with queries using the **Run** command and save once you complete your solution.

You can also search for other queries by namespace and/or name and search for mappings by tenants and/or name.

Finally, you can also archive queries or revisions and see the list of archived items on the side menu.

When a query is archived it is only shown if the _Archived_ filter is enabled on the side menu. If a query is not archived but all its revisions are, then it will be shown on the normal queries list, but disabled.

When a non archived query is selected, just non archived revisions are listed on the select menu. If an archived query is selected, all revisions are listed on the select menu.

## Deploying

First build the project, which will produce a directory with static files.

The restQL Manager depends on one environment variable, `REACT_APP_RESTQL_URL`. Hence, the command to built it must be like this:

```bash
REACT_APP_RESTQL_URL=http://my-restql.corp/ npm run build
```

The static files generate can then be served from a object storage like AWS S3 or, more traditionally, with an Nginx or Expresse server.

<div style="display:flex; width: 100%; align-items: center; justifiy-content: stretch">
  <img src="./assets/images/manager.png"/>
</div>
