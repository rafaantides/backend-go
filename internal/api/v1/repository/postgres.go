package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewPostgreSQL(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao pingar o banco: %w", err)
	}

	// TODO: rever a necessidade disso daqui
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Banco de dados conectado")
	return &Database{DB: db}, nil
}

func (d *Database) Close() {
	if err := d.DB.Close(); err != nil {
		log.Println("Erro ao fechar conexão com o banco:", err)
	}
}
