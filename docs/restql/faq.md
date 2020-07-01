# Frequently Asked Questions

### **How does RestQL orchestrate multiple APIs?**

RestQL uses communicating sequential processes (CSP) in order to make concurrent calls. Anything that can be concurrent (e.g. API calls that aren't dependent on each other) will be concurrent.

### **Which response time should be expected in chained calls?**

When there are chained calls, the response time is the sum of the dependent services. Because of parallelism, if there are no dependent services in a query, the response time will be the maximum of the response times of all the requests.

### **Is it better to save queries at the database or at the YAML configuration file?**

For learning purposes, saving your queries at `restql.yml` should be easier, since it needs no database configuration, but it has the downside that at each change in the configuration file will require a restart of the application in order for the changes take place.

Saving your queries at the database requires a bit of configuration by the plugin, but you won't need to restart the application each time you want to add a new query or resource.

### **How does RestQL ordenate responses?**

In most cases, ordenation responsability is delegated to the API which is being consulted.

### **And what about pagination?**

This is the same case as in **ordenation**, pagination responsability is delegated to the API which is being consulted.

### **Is there any cache in restQL?**

RestQL caches only the queries text, mappings and parsed query, not the result of its execution.
