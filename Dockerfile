FROM golang:1.23.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /app/main .
COPY --from=build /app/.env .

EXPOSE 8080

CMD ["./main"]