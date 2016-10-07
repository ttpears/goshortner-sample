#goshortner

A URL shortner written in Golang as a sample project

##Requirements

* [Go](https://golang.org/)
* [Redis](http://redis.io/)

##Fetch and install:

```bash
go get github.com/ttpears/goshortner-sample
```

##Usage:

```bash
$GOPATH/bin/goshortner-sample
```

###Create a new short url:

```bash
curl localhost:8080/add -d "url=http://example.com"
```

###View stats:

To view stats for a given short URL, add `/stats` to it:

```bash
curl localhost:8080/DmN93nj/stats
```
