# Plugins

restQL have a plugin system that works by compiling the plugins together with restQL application into a custom binary.

This is done by leveraging Go Modules with the help of the [restQL-cli](https://github.com/b2wdigital/restQL-cli).

## Plugin types

Currently, restQL supports following types of plugins:

### Lifecycle 

Defined by the interface `restql.LifecyclePlugin`, it allows you to execute code at various points of the query execution, like before and after an HTTP request is made.

This plugin type is specially useful for monitoring purposes, since it allows you to derive countless metrics from the given data. 

### Database

Defined by the interface `restql.DatabasePlugin`, it allows you to use any an external database to store mappings and queries.

The methods on this interface that support the archiving feature on the Administrative API have some assumptions upon its implementation:
- `UpdateQueryArchiving`: when a query is archived through this method, all its revisions must be also marked as archived. Also, when a query is unarchived its revisions must remain archived.
- `UpdateRevisionArchiving`: when a revision is unarchived its query must also be marked as unarchived.

## Developing plugins

> It is strongly recommended having the [restQL-cli](https://github.com/b2wdigital/restQL-cli) installed locally.

A restQL plugin should be a Go Module project with an `init` function that calls `restql.RegisterPlugin` to set the plugin at restQL initialization.

For example, if you are developing a Lifecycle plugin, the recommended entrypoint file should look like the following:

```go
package main

import "github.com/b2wdigital/restQL-golang/v6/pkg/restql"

func init() {
    restql.RegisterPlugin(restql.PluginInfo{
        Name: "myplugin",
        Type: restql.LifecyclePluginType,
        New: func(logger restql.Logger) (restql.Plugin, error) {
            return NewMyPlugin(logger)
        },
    })
}

func NewMyPlugin(log restql.Logger) (restql.LifecyclePlugin, error) {
    return MyPlugin{}, nil
}
``` 

The `restql.RegisterPlugin` expects a `restql.PluginInfo` with three fields:
- Name: a string used to identify your plugin
- Type: a constant which defines the plugin type, restQL provides this values for each possibility.
- New: a constructor that return a fresh value of your plugin.

If you are using the [restQL-cli](https://github.com/b2wdigital/restQL-cli) you can use it to run and build the plugin locally with restQL to verify the integration. 

### Best Practices

#### Compilation safety

Use a separate function to build your plugin with a signature enforcing the plugin interface that you are implementing, for example the `NewMyPlugin` function specify that the value returned by it should be a `restql.LifecyclePlugin`. 

This will provide compile time verification that your internal type implements all the necessary methods, since the `restql.Plugin` interface is a common interface between all plugin types that only expect a `Name() string` method.

#### Logging

The logger instance given in the `New` constructor has no context since it is the one used during restQL initialization and should be used to log information about the plugin initialization.

In order to log information about the execution of the plugin we suggest the logger to be extracted from the `context.Context` passed to each method using the `restql.GetLogger` helper function. The logger returned by this helper will have all the context of the current query being processed and will improve the debugging when the time comes.