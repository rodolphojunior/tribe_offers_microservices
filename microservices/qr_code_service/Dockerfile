# Etapa 1: Construir o binário Go
FROM golang:1.22-alpine AS build

# Definir o diretório de trabalho
WORKDIR /app

# Copiar go.mod e go.sum para baixar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código do microsserviço de QR code
COPY . ./ 

# Construir o binário do microsserviço de QR code
RUN go build -o /go/bin/qr_code_service .

# Etapa 2: Executar o binário
FROM alpine:3.18

WORKDIR /root/

# Copiar o binário construído na etapa anterior
COPY --from=build /go/bin/qr_code_service .

# Expor a porta que o microsserviço vai usar
EXPOSE 8086

# Comando para iniciar o microsserviço
CMD ["./qr_code_service"]
