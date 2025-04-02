package handlers

import (
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func ProcessDebts(messageBody []byte) error {
	var msgData dto.DebtRequest
	err := json.Unmarshal(messageBody, &msgData)
	if err != nil {
		log.Printf("Erro ao decodificar JSON: %v", err)
		return err
	}

	processRequest(msgData)
	return nil
}

func processRequest(msgData dto.DebtRequest) error {

	url := os.Getenv("API_URL")

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Erro ao fazer requisição: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler resposta: %v", err)
		return err
	}

	var response models.Debt
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Erro ao decodificar o response: %v", err)
		return err
	}

	return nil
}
