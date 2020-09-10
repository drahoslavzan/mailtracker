# BUILD
FROM golang:1.14-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main

# RUN
FROM alpine

RUN adduser -S -D -H -h /app app
USER app
COPY --from=builder /build/main /build/database/.env* /app/

WORKDIR /app
EXPOSE 8000
CMD [ "./main" ]