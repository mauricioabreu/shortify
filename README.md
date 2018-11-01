# shortify [![Build Status](https://travis-ci.org/mauricioabreu/shortify.svg?branch=master)](https://travis-ci.org/mauricioabreu/shortify)

URL shortener written in golang.

`redis` is used to store the shortened URLs.

## What is it?

It is an URL shortener. It is a service you can use to shorten your URLs. Similar to services like `bit.ly` and `bl.ink`.

## Why?

It is just an experiment, a toy project to learn and exercise my systems design skills.
I want this project to be:

* Fast
* Testable
* Deployable
* Debuggable

## Using

There is support for docker-compose:

```
docker-compose up
```

Now you can run the server and use the service.

Shortening an URL:

```
curl -sv http://localhost:8000/ -XPOST -d "url=https://www.maugzoide.com/"
```

Retrieving the URL:

```
curl -sv http://localhost:8000/d41d8cda2f -XGET
```

## Benchmark

One of the premises of this service is to be *fast*. To show how fast this service is, we need to benchmark the service running. We run the server and use a tool to perform a *load testing*, sending multiple HTTP requests as if it is facing a high number of people using it.

This report uses the [hey](https://github.com/rakyll/hey) tool. It is a replacement for `Apache Benchmark`. 

What are we going to measure:

* **Latency** - how fast a server responds to the requests in the load test. HTTP requests are a bit 
* **Throughput** - how many requests the server can succesfully respond in a specific interval.

### First test (low number of requests, high concurrency):

Requests: 1000
Concurrency: 500

```
hey -n 1000 -c 500 -m GET -disable-redirects http://localhost:8000/5629ae52dd
```

Summary:

Total:        0.1454 secs

Slowest:      0.1197 secs

Fastest:      0.0028 secs

Average:      0.0475 secs

Requests/sec: 6879.2310

### Second test (high number of requests, low concurrency):

Requests: 10000
Concurrency: 50

```
hey -n 10000 -c 50 -m GET -disable-redirects http://localhost:8000/5629ae52dd
```

Summary:

Total:        0.5026 secs

Slowest:      0.0288 secs

Fastest:      0.0002 secs

Average:      0.0024 secs

Requests/sec: 19897.3339

### Third test (high number of requests, high concurrency):

Requests: 100000
Concurrency: 500

```
hey -n 100000 -c 500 -m GET -disable-redirects http://localhost:8000/5629ae52dd
```

Summary:

Total:        4.4532 secs

Slowest:      1.0243 secs

Fastest:      0.0002 secs

Average:      0.0217 secs

Requests/sec: 22455.7735
