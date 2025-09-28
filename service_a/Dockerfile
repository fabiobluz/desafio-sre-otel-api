FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar o serviço
COPY service_a/ ./

# Configurar GOPROXY
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# Download das dependências e build
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o main .

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar o binário
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
