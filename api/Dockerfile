FROM golang:1.10
RUN go get github.com/mongodb/mongo-go-driver/bson
RUN go get github.com/mongodb/mongo-go-driver/mongo
ADD . /go/src/github.com/shreshthkhilani/royal-squid
RUN go install github.com/shreshthkhilani/royal-squid
ENTRYPOINT /go/bin/royal-squid
EXPOSE 8080