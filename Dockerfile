# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -o /out/bookshelf ./cmd/server


# Stage 2: Runtime
FROM alpine:3.24

RUN apk add --no-cache ca-certificates

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/bookshelf ./bookshelf

USER app

EXPOSE 8080

CMD ["./bookshelf"]