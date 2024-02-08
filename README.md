# RateLimiter

In general, a rate limiter caps how many requests a sender can issue in a specific time window.

It then blocks requests once the cap is reached.

There are different techniques for measuring and limiting rates, each with their own uses and implications.

`TokenBucket` is a Go implementation of a rate limiter using the `Token Bucket Algorithm`.

It allows you to limit the rate of events, such as requests, in the specific time window.

The Token Bucket Algorithm works as follows: 
> A token bucket is a container that has pre-defined capacity. Tokens are put in the bucket at preset rates periodically. Once the bucket is full, no more tokens are added.

## Usage

Import the `ratelimiter` package into your code:

```go
import "github.com/deep-quality-dev/ratelimiter"
```

Create a TokenBucket with the desired parameters:
```go
tb := ratelimiter.NewTokenBucket(capacity, refillAmount, refillInterval)
```

Use the ConsumeTokens method to check and consume tokens:
```go
if tb.ConsumeTokens(1) {
    // Perform the action for which tokens were consumed
    fmt.Println("Action allowed")
} else {
    // Action not allowed due to rate limiting
    fmt.Println("Rate limit exceeded")
}
```

## Run Test

Execute the following command to run unit test:

```shell
go test ./...
```

### Run Tests with coverage

```shell
go test -cover ./...
```

If you want to export coverage report, execute the following command:

```shell
go tool cover -html=coverage.out
```