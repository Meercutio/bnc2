# build stage
FROM golang:1.24-alpine AS build
WORKDIR /src

# важно для go install в alpine
RUN apk add --no-cache git ca-certificates

# (по желанию) чтобы всегда работало через модульный прокси
ENV GOPROXY=https://proxy.golang.org,direct

RUN go env && go install -x github.com/pressly/goose/v3/cmd/goose@v3.21.1
ENV PATH="/go/bin:${PATH}"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/bc-server ./cmd/server

# run stage
FROM alpine:3.20
WORKDIR /app

COPY --from=build /out/bc-server /app/bc-server
COPY --from=build /go/bin/goose /usr/local/bin/goose

COPY db/migrations/ /app/migrations/

ENV MIGRATIONS_DIR=/app/migrations

ENV PORT=8080
EXPOSE 8080

RUN adduser -D -H appuser && chown -R appuser:appuser /app
USER appuser

CMD ["/app/bc-server"]
