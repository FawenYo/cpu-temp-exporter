FROM golang:1.23.1-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum cmd/ ./

RUN go mod tidy && go build -o cpu_temp_exporter .

# Final image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cpu_temp_exporter .

EXPOSE 8080

CMD ["./cpu_temp_exporter"]
