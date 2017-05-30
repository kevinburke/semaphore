
semaphore
=========

Semaphore implementation in golang

[![GoDoc](https://godoc.org/github.com/kevinburke/semaphore?status.svg)](https://godoc.org/github.com/kevinburke/semaphore)
[![Go Report Card](https://goreportcard.com/badge/github.com/kevinburke/semaphore)](https://goreportcard.com/report/github.com/kevinburke/semaphore)

### Usage
Initiate
```go
import "github.com/kevinburke/semaphore"
...
sem := semaphore.New(5) // new semaphore with 5 permits
```
Acquire
```go
sem.Acquire() // one
sem.AcquireMany(n) // multiple
sem.AcquireWithin(n, time.Second * 5) // timeout after 5 sec
```
Release
```go
sem.Release() // one
sem.ReleaseMany(n) // multiple
```

### documentation

See here: https://godoc.org/github.com/kevinburke/semaphore
