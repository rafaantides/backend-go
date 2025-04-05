package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"backend-go/internal/api/v1/repository/models"
	"backend-go/internal/api/v1/repository/postgresql"
	"backend-go/pkg/ent/category"
	"backend-go/pkg/ent/paymentstatus"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Popula o banco com dados iniciais",
	Run: func(cmd *cobra.Command, args []string) {
		runSeed()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func runSeed() {
	_ = godotenv.Load()

	db := connectDatabase()
	defer db.Close()

	ctx := context.Background()

	fmt.Println("🌱 Iniciando seed do banco...")

	if err := seedPaymentStatuses(ctx, db); err != nil {
		log.Fatalf("erro ao criar payment statuses: %v", err)
	}

	if err := seedCategories(ctx, db); err != nil {
		log.Fatalf("erro ao criar categories: %v", err)
	}

	if err := seedDebts(ctx, db, "./static/json/debts.json"); err != nil {
		log.Fatalf("erro ao criar debts: %v", err)
	}

	fmt.Println("✅ Seed concluído com sucesso!")
}

func seedPaymentStatuses(ctx context.Context, db *postgresql.PostgreSQL) error {
	statuses := []struct {
		Name        string
		Description string
	}{
		{"pending", "Pagamento pendente"},
		{"paid", "Pagamento realizado"},
		{"failed", "Pagamento falhou"},
	}

	for _, s := range statuses {
		exists, err := db.Client.PaymentStatus.Query().Where(paymentstatus.NameEQ(s.Name)).Exist(ctx)

		if err != nil {
			return err
		}
		if exists {
			continue
		}

		_, err = db.Client.PaymentStatus.
			Create().
			SetName(s.Name).
			SetDescription(s.Description).
			Save(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Criado: %s\n", s.Name)
	}
	return nil
}

func seedCategories(ctx context.Context, db *postgresql.PostgreSQL) error {
	categories := []string{
		"Transporte",
		"Bebidas e Conveniência",
		"Restaurantes e Alimentação",
		"Mercado e Compras",
		"Assinaturas e Serviços Digitais",
		"Entretenimento e Eventos",
		"Farmácias e Saúde",
		"Vestuário e Cosméticos",
		"Barbearia e Beleza",
		"Eletrônicos e Tecnologia",
		"Ótica e Acessórios",
	}

	for _, name := range categories {
		exists, err := db.Client.Category.Query().Where(category.NameEQ(name)).Exist(ctx)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		_, err = db.Client.Category.
			Create().
			SetName(name).
			Save(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Criado: %s\n", name)
	}
	return nil
}

func seedDebts(ctx context.Context, db *postgresql.PostgreSQL, jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo JSON: %w", err)
	}

	var debts []models.Debt
	if err := json.Unmarshal(data, &debts); err != nil {
		return fmt.Errorf("erro ao parsear JSON: %w", err)
	}

	for _, d := range debts {
		_, err := db.Client.Debt.
			Create().
			SetTitle(d.Title).
			SetAmount(d.Amount).
			SetPurchaseDate(d.PurchaseDate).
			SetDueDate(d.DueDate).
			// SetInvoiceID e SetCategoryID são opcionais
			Save(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar dívida '%s': %w", d.Title, err)
		}
		fmt.Printf("✅ Dívida criada: %s\n", d.Title)
	}
	return nil
}
