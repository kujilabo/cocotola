FROM golang:1.18-alpine as builder

RUN apk add --no-cache build-base

WORKDIR /go/app
COPY . .

RUN go build -o cocotola ./main.go

# Application image.
FROM alpine:latest

RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=builder /go/app/cocotola .
COPY --from=builder /go/app/configs ./configs
COPY --from=builder /go/app/sqls ./sqls

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN chown -R appuser /app

USER appuser

EXPOSE 8080

CMD ["./cocotola"]
