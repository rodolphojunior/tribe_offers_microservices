# Etapa 1: Construir o binário Go
FROM golang:1.22-alpine AS build

# Definir o diretório de trabalho
WORKDIR /app

# Copiar go.mod e go.sum para baixar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código do projeto
COPY . .

# Construir o binário do servidor backend
RUN go build -o /go/bin/tribo_ofertas_backend cmd/main.go

# Construir o binário para o atualizador de ofertas
RUN go build -o /go/bin/offer_updater cmd/offer_updater.go

# Etapa 2: Executar o binário
FROM alpine:3.18

WORKDIR /root/

# Copiar os binários construídos na etapa anterior
COPY --from=build /go/bin/tribo_ofertas_backend .
COPY --from=build /go/bin/offer_updater .

# Expor a porta que o servidor vai usar
EXPOSE 8080

# Comando para iniciar o servidor
CMD ["./tribo_ofertas_backend"]

