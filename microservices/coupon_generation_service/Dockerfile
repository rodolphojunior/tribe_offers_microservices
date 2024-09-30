# Etapa 1: Construir o binário Go
FROM golang:1.22-alpine AS build

WORKDIR /app
# Copiar o go.mod e go.sum para o container
COPY go.mod go.sum ./
# Baixar as dependências
RUN go mod download
COPY . ./ 
RUN go build -o /go/bin/coupon_generation_service .
FROM alpine:3.18
WORKDIR /root/
COPY --from=build /go/bin/coupon_generation_service .
EXPOSE 8084
CMD ["./coupon_generation_service"]
