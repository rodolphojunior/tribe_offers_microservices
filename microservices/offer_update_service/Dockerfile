# Dockerfile para o Offer Update Service
FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY ./microservices/offer_update_service .

RUN go build -o /go/bin/offer_update_service .

FROM alpine:3.18
WORKDIR /root/
COPY --from=build /go/bin/offer_update_service .
CMD ["./offer_update_service"]