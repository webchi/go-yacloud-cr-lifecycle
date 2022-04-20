FROM golang:1.18.1  AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR $GOPATH/src
ADD go.mod .
ADD go.sum .
ADD . .
RUN go mod download
RUN go build -o /go/main

FROM alpine
COPY --from=builder /go/main /go/main
EXPOSE 8080
WORKDIR /go
ENTRYPOINT ["./main"]
