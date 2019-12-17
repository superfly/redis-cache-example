# Redis cache example

This is a simple app written in go that connects to the fly redis cache defined in the FLY_REDIS_CACHE_URL env variable. It shows how to get, set, and delete cached data using the request path as the key. 


**Setting data**

```bash
curl -XPOST -d "Michael" https://go-redis-cache-example.fly.dev/name
```

**Getting data**

```bash
curl https://go-redis-cache-example.fly.dev/name
```


**Deleting data**

```bash
curl -XDELETE https://go-redis-cache-example.fly.dev/name
```
