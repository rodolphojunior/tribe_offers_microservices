# Dockerfile para o Offer Update Service
FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./ 

RUN go build -o /go/bin/coupon_redemption_service .

FROM alpine:3.18
WORKDIR /root/
COPY --from=build /go/bin/coupon_redemption_service .
# Expor a porta que o microsserviço vai usar
EXPOSE 8089
CMD ["./coupon_redemption_service"]
