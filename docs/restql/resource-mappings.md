# Resource Mappings

RestQL provides three ways to map resources

1. Environment variables
2. Configuration File
3. Database

### Environment variables

RestQL will detect that an environment variable is a mapping if it follows the pattern `RESTQL_MAPPING_MYTENANT_RESOURCENAME`, for example `RESTQL_MAPPING_UNIVERSE_HERO=http://hero.api/` will create a mapping with name `hero` and target `http://hero.api/` under the tenant `UNIVERSE`.

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
