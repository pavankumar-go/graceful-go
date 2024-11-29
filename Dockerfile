FROM golang:1.18 as builder

WORKDIR ~/go/
COPY . .

RUN GOOS=linux go build  -o /graceful  main.go

ENTRYPOINT ["/graceful"]
