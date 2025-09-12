# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags "-s -w" -o app .
RUN apk update && apk add upx 

RUN upx --best --lzma app

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/internal/pkg/database/migrations/postgresql/ ./internal/pkg/database/migrations/postgresql/

EXPOSE 8080/tcp

CMD ["./app"]