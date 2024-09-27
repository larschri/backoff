# backoff

A small library for implementing exponential backoff in a simple for loop.

```go
for range backoff.Exponential(context.Background()) {
    t, err := Thing()
    if err == nil {
        return t
    }
}
```

Some properties can be configured. The code above is equivalent to:

```go
for range backoff.Exponential(context.Background(),
    backoff.WithMultiplier(2),    // wait time is multiplied by this value after each iteration
    backoff.WithMin(time.Second), // wait time after the first iteration
    backoff.WithMax(time.Hour),   // the maximum wait time
    backoff.WithTerminate(false), // run forever, do not stop when maximum wait time is reached
) {
    t, err := Thing()
    if err == nil {
        return t
    }
}
```
