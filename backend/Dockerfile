FROM golang:1.23.4-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /bin/adhp ./cmd/adhp.go

FROM scratch

WORKDIR /app

COPY --from=builder /bin/adhp .

EXPOSE 8301

ENTRYPOINT ["./adhp"]