# build stage
FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/bc-server ./cmd/server

# run stage
FROM alpine:3.20
WORKDIR /app

COPY --from=build /out/bc-server /app/bc-server

ENV PORT=8080
EXPOSE 8080

RUN adduser -D -H appuser && chown -R appuser:appuser /app
USER appuser

CMD ["/app/bc-server"]
