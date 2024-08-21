# Etapa de construção
FROM golang:1.23 AS build

# Defina o diretório de trabalho
WORKDIR /app

# Copie go.mod e go.sum e faça o download das dependências
COPY go.mod go.sum ./
RUN go mod download

# Copie o código fonte
COPY . .

# Instale o Air para recarregamento automático
RUN go install github.com/air-verse/air@latest

# Compile o aplicativo
RUN go build -o main .

# Etapa de execução
FROM debian:bullseye-slim

# Defina o diretório de trabalho
WORKDIR /app

# Copie o binário da etapa de construção
COPY --from=build /app/main .
COPY --from=build /go/bin/air /usr/local/bin/air

# Mapeie o volume de trabalho
VOLUME ["/app"]

# Defina o comando padrão
CMD ["air"]
