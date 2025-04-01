package handlers

// import (
// 	"encoding/csv"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/xuri/excelize/v2"
// )

// func ProcessFile(c *gin.Context) {
// 	file, _, err := c.Request.FormFile("file")
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "Falha ao obter o arquivo"})
// 		return
// 	}
// 	defer file.Close()

// 	fileType, err := detectFileType(file)
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "Formato de arquivo não suportado"})
// 		return
// 	}

// 	if fileType == "csv" {
// 		err = h.processCSV(file)
// 	} else if fileType == "xlsx" {
// 		err = h.processXLSX(file)
// 	}

// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(200, gin.H{"message": "Arquivo processado com sucesso"})
// }

// func (h *Handler) processCSV(file multipart.File) error {
// 	reader := csv.NewReader(file)
// 	rows, err := reader.ReadAll()
// 	if err != nil {
// 		return err
// 	}

// 	if len(rows) < 2 {
// 		return fmt.Errorf("arquivo CSV inválido")
// 	}

// 	headers := rows[0]
// 	columnIndex := make(map[string]int)
// 	for i, header := range headers {
// 		columnIndex[strings.ToLower(header)] = i
// 	}

// 	for _, line := range rows[1:] {
// 		debt := DebtRequest{
// 			InvoiceID:    line[columnIndex["invoice_id"]],
// 			PurchaseDate: line[columnIndex["purchase_date"]],
// 			DueDate:      line[columnIndex["due_date"]],
// 			Title:        line[columnIndex["title"]],
// 			Amount:       line[columnIndex["amount"]],
// 		}
// 		message, _ := json.Marshal(Message{Data: debt})
// 		h.QueueService.SendMessage(message)
// 	}
// 	return nil
// }

// func (h *Handler) processXLSX(file multipart.File) error {
// 	f, err := excelize.OpenReader(file)
// 	if err != nil {
// 		return err
// 	}

// 	sheetName := f.GetSheetName(0)
// 	rows, err := f.GetRows(sheetName)
// 	if err != nil {
// 		return err
// 	}

// 	if len(rows) < 2 {
// 		return fmt.Errorf("arquivo XLSX inválido")
// 	}

// 	headers := rows[0]
// 	columnIndex := make(map[string]int)
// 	for i, header := range headers {
// 		columnIndex[strings.ToLower(header)] = i
// 	}

// 	for _, row := range rows[1:] {
// 		debt := DebtRequest{
// 			InvoiceID:    row[columnIndex["invoice_id"]],
// 			PurchaseDate: row[columnIndex["purchase_date"]],
// 			DueDate:      row[columnIndex["due_date"]],
// 			Title:        row[columnIndex["title"]],
// 			Amount:       row[columnIndex["amount"]],
// 		}
// 		message, _ := json.Marshal(Message{Data: debt})
// 		h.QueueService.SendMessage(message)
// 	}
// 	return nil
// }

// func detectFileType(file multipart.File) (string, error) {
// 	buffer := make([]byte, 4)
// 	_, err := file.Read(buffer)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Reset file reader position
// 	file.Seek(0, io.SeekStart)

// 	if buffer[0] == 0x50 && buffer[1] == 0x4B {
// 		return "xlsx", nil
// 	}

// 	return "csv", nil
// }
