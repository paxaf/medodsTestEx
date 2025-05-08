FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -o my_app ./cmd

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache curl ca-certificates

RUN apk update

RUN curl -fsSL \
https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
sh

RUN mkdir migrations && mkdir config

COPY --from=builder app/config/*.yaml config/
COPY --from=builder /app/my_app my_app
COPY migrations/ ./migrations/
COPY entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh
RUN sed -i 's/\r$//' entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]