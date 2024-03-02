# go-cache

```go

// cache instance that will expire objects if not used for 30 seconds
c := cache.New(cache.LastAccess[*Object](30*time.Second))

// put some object to the cache
c.Put("1", &Object{})

// try to get object
object, ok := c.Peek("1")
if ok {
  // object present in cache
}

ctx := context.Background()
db := NewDatabase()

// loading object if not in cache, it will remain in cache for next 30 sec
object, err := c.Get("2", func(id string) (*Object, error) {
  return loadObjectFromDatabase(ctx, db, id)
})

c.Remove("1", "2") // remove objects from cache by id

```
