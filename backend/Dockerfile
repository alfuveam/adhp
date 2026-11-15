FROM golang:1.23.4-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /bin/adhp ./cmd/adhp.go

FROM scratch

WORKDIR /app

COPY --from=builder /bin/adhp .
# scratch não tem certificados CA: necessários para o TLS do SMTP e do Turnstile
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8301

ENTRYPOINT ["./adhp"]