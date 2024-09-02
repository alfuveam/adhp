FROM golang:1.23.0-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /bin/tcc ./cmd/tcc.go

FROM scratch

WORKDIR /app

COPY --from=builder /bin/tcc .

EXPOSE 8080

ENTRYPOINT ["./tcc"]