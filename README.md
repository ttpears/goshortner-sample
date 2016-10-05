#goshortner
A URL shortner written in Golang as a sample project

##Prerequisites
####Get redigo:
`go get github.com/garyburd/redigo/redis`

####Get mux:
`go get github.com/gorilla/mux`

###Get code:
`git clone https://github.com/ttpears/goshortner-sample.git`

##Build:
`go build`

##Usage:
`./goshortner`

###Create a new short url: 
`curl localhost:8080/add -d "url=http://example.com"`

###View stats:
`curl localhost:8080/shorturl/stats`
