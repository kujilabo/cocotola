FROM golang:1.18-alpine as builder

RUN apk add --no-cache build-base=0.5-r3

WORKDIR /app
COPY . .

RUN go build -o cocotola ./src/main.go

# Application image.
FROM alpine:3.16.2

RUN apk --no-cache add tzdata=2022f-r0

WORKDIR /app

COPY --from=builder /app/cocotola .

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser /app

USER appuser

EXPOSE 8080 8081

CMD ["./cocotola"]
