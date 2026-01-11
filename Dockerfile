# build stage
FROM golang:1.24-alpine AS build
WORKDIR /src

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/bc-server ./cmd/server

# run stage
FROM alpine:3.20
WORKDIR /app

# бинарник и статика
COPY --from=build /out/bc-server /app/bc-server
COPY web /app/web

# Fly будет проксировать наружу, нам достаточно слушать 8080
ENV ADDR=:8080
EXPOSE 8080

CMD ["/app/bc-server"]
