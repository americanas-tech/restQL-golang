# Resource Mappings

RestQL provides three ways to map resources

1. Environment variables
2. Configuration File
3. Database

### Environment variables

RestQL will detect that an environment variable is a mapping if it starts with `RESTQL_MAPPING_` followed by the tenant (that will be used as is) name, another underline and then resource name (that will be lowercased) and the value should be the target url, for example, `RESTQL_MAPPING_MYTENANT_HERO=http://hero.api/` will create a mapping with name `hero` and target `http://hero.api/`.

These mappings overwrite any mapping with the same name present in database or configuration file.

### Configuration file

You can define mappings in the configuration file like:

```yaml
tenants:
  my-tenant:
    hero: http://hero.api/
```

These mappings can be overwritten by a mapping with the same name present in the database or the environment.

### Database

You can add support to store mappings to a database trough a Database Plugin. You can learn more about it in the [Plugins documentation](/restql/plugins.md). 

In a production environment we recommend the use of the [restQL Manager](/restql/manager.md) to manage the mappings in a database rather than manually.
