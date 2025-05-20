FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -buildvcs=false -o main

FROM alpine:3

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main /main
COPY .env_example .env 

EXPOSE 8000

CMD ["/main"]
