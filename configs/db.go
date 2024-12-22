package configs

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn
var err error

func InitDatabase() {
	DB, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Erro ao configurar o banco de dados: %v", err)
	}
}

func CloseDatabase() {
	DB.Close(context.Background())
}
