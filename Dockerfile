FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o grpc-orders ./cmd/grpc-orders/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/grpc-orders ./
COPY --from=builder /app/.env ./

EXPOSE 50055

ENTRYPOINT ["./grpc-orders"]