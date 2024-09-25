# backoff

A small library for implementing exponential backoff in a simple for loop.

```go
	for range Exponential(context.Background(), WithMax(4*time.Second)) {
		if _, err := os.Hostname(); err == nil {
			fmt.Println("yay, hostname")
			break
		}
	}

```

