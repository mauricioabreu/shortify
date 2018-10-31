# shortify

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
curl -sv http://localhost:8000/d41d8cd -XGET
```
