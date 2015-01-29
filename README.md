A URL shortnet writting in Golang, as a sample project

To use:

Get code:
git clone https://github.com/ttpears/goshortner-sample.git

Get redigo:
go get github.com/garyburd/redigo/redis

Get mux:
go get github.com/gorilla/mux

Build:
go build

Run:
./goshortner

Usage:

Create a new short url: 
curl <host>:8080/add -d "url=<longurl>"

View stats:
curl <host>:8080/<shorturl>/stats
