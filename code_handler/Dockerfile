# Etapa de build do Go
FROM golang:1.23.0-alpine as builder

WORKDIR /app

# Copiar os arquivos de dependência e baixar os módulos Go
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copiar o código-fonte e construir o binário
COPY . .
RUN go build -o /bin/code_handler ./code_handler.go

# Etapa final
FROM alpine:latest

WORKDIR /app

# Copiar o binário Go do build anterior
COPY --from=builder /bin/code_handler .

# Instalar Python e pytest usando apk
RUN apk add --no-cache python3 py3-pytest

EXPOSE 8083

# Definir o ponto de entrada para o binário Go
ENTRYPOINT ["./code_handler"]
