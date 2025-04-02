package postgresql

import (
	"backend-go/internal/api/v1/interfaces"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	DB *sql.DB
}

func NewPostgreSQL(dsn string) (interfaces.Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao pingar o banco: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Banco de dados conectado")
	return &PostgreSQL{DB: db}, nil
}

func (d *PostgreSQL) Close() {
	if err := d.DB.Close(); err != nil {
		log.Println("Erro ao fechar conexão com o banco:", err)
	}
}
