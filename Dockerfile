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

# Defina o comando padrão
CMD ["air"]
