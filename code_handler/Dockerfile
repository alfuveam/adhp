FROM golang:1.23.4-alpine as builder

WORKDIR /app

COPY ./handler/go.mod ./handler/go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /bin/code_handler ./handler/code_handler.go

# para usar o apk
FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/local/go /usr/local/go
COPY --from=builder /bin/code_handler .

RUN apk add --no-cache python3 py3-pytest
ENV PATH=$PATH:/usr/local/go/bin

EXPOSE 8082

ENTRYPOINT ["./code_handler"]
