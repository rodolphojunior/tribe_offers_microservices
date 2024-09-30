# Dockerfile para o Offer Ranking Service
FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./ 

RUN go build -o /go/bin/offer_ranking_service .

FROM alpine:3.18

WORKDIR /root/

COPY --from=build /go/bin/offer_ranking_service .

EXPOSE 8088

CMD ["./offer_ranking_service"]