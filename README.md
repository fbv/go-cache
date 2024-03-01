# go-chache

```
func loadObjectFromDatabase(ctx context.Context, db *Database, id string) (*Object, error){
  ...
  return object, nil
}


ctx := context.Background()
db := NewDatabase()

c := cache.New(cache.LastAccess[*Object](30*time.Second))
object, err := c.Get("id", func(id string) (*Object, error) {
  return loadObjectFromDatabase(ctx, db, id)
})
```
