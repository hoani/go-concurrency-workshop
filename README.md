## Go Concurrency Testing Workshop

This code is part of a workshop on testing concurrent code in golang.

The `main` package provides an `http` server on port `8080`:

```sh
go run .
```

In another terminal
```sh
curl localhost:8080/line1
```

Responds with:
```
Knock knock
```