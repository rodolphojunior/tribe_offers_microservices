# Etapa 1: Construir o binário Go
FROM golang:1.22-alpine AS build

WORKDIR /app

# Copiar go.mod e go.sum para baixar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código do projeto
COPY . .

# Construir o binário do servidor backend
RUN go build -o /go/bin/tribe_offer_microservices_backend cmd/main.go

# Etapa 2: Executar o binário
FROM alpine:3.18

WORKDIR /root/

# Copiar o binário da etapa anterior
COPY --from=build /go/bin/tribe_offer_microservices_backend .

# Expor a porta que o servidor vai usar
EXPOSE 8080

# Comando para iniciar o servidor
CMD ["./tribe_offer_microservices_backend"]


