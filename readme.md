# Aggregation Plugin

Aggregation Plugin is a middleware for [Traefik](https://github.com/traefik/traefik) which performs simple request aggregation that receives POST with JSON containing URLS on configured route and returns aggregated JSON.

## Configuration

## Static

```toml
[pilot]
    token="xxx"

[experimental.plugins.aggregation]
    modulename = "github.com/traefik/aggregation-plugin"
    version = "v0.1.0"
```

## Dynamic

To configure the `Aggregation` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example illustates
the usage of `aggregation` middleware plugin on specific path (`/aggregate`). 

```toml
[http.routers]
  [http.routers.my-router]
    rule = "Path(`/aggregate`)"
    middlewares = ["aggregate"]
    service = "my-service"

# Specify server for aggregation
[http.middlewares]
  [http.middlewares.aggregate.plugin]
    server = "http://127.0.0.1"

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```
## Request example

```
POST http://localhost/aggregate

{
    "foo": "foo/1",
    "bar": "bar/2?foo=3"
}
```

## Response example

```
{
   "foo":{
      "foo-key":"foo-value"
   },
   "bar":{
      "bar-key":"bar-value"
   }
}
```