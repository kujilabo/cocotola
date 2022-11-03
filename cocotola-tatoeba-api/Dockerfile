FROM golang:1.18-alpine as builder

RUN apk add --no-cache build-base

WORKDIR /app
COPY . .

RUN go build -o cocotola ./src/main.go

# Application image.
FROM alpine:latest

RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=builder /app/cocotola .

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser /app

USER appuser

EXPOSE 8080 8081

CMD ["./cocotola"]
