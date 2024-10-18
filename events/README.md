# How to use this

```go
func run() error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

    queries := events.New(conn)
    queries.GetEvent(ctx, id)
```
