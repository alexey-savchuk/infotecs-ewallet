FROM golang:1.21.7-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o bin/ewallet cmd/ewallet/ewallet.go


FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/ewallet .

EXPOSE 8080

CMD ["./ewallet"]
