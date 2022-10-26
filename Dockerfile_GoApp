FROM golang:1.19-alpine3.16 as builder
WORKDIR /app
COPY . .
RUN go build main.go

FROM alpine:3.16 as final
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD [ "/app/main" ]