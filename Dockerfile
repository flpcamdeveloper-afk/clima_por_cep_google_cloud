FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 go build -o /climacep ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /climacep .

# O Cloud Run injeta a variável PORT em tempo de execução; 8080 é o padrão
# usado tanto localmente quanto pelo Cloud Run quando PORT não é definida.
EXPOSE 8080

ENTRYPOINT ["./climacep"]
