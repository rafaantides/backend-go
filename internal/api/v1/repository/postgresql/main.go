package postgresql

import (
	"backend-go/pkg/ent"
	"backend-go/pkg/hooks"
	"context"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	Client *ent.Client
}

// func NewPostgreSQL(dsn string) (interfaces.Database, error) {
func NewPostgreSQL(dsn string) (*PostgreSQL, error) {
	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	client := ent.NewClient(ent.Driver(drv))

	// TODO: ver como isso funciona na pratica depois
	// Colocar em outro lugar??
	client.Debt.Use(hooks.SetDefaultStatusHook(client))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Schema.Create(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("erro ao criar schema: %w", err)
	}

	log.Println("Banco de dados conectado com sucesso via Ent")
	return &PostgreSQL{Client: client}, nil
}

func (d *PostgreSQL) Close() {
	if err := d.Client.Close(); err != nil {
		log.Println("Erro ao fechar conexão com o banco:", err)
	} else {
		log.Println("Conexão com o banco fechada.")
	}
}
