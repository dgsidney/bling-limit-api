package handlers

import (
	"bling_limit/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ViolationReport struct {
	Codigo     string
	Deposito   string
	Tentativas int64
}

type ViolationNotification struct {
	Codigo     string `json:"codigo"`
	Deposito   string `json:"deposito"`
	Tentativas int64  `json:"tentativas"`
	Message    string `json:"message"`
}

func sendViolationNotification(violation *ViolationReport) {

	message := fmt.Sprintf("Violação: Código %s excedeu o limite com + de %d tentativas. (Depósito ID: %s)", violation.Codigo, violation.Tentativas, violation.Deposito)

	callbackURL := os.Getenv("CALLBACK_ENDPOINT")
	if callbackURL == "" {
		log.Println("CALLBACK_ENDPOINT não está definido.")
		return
	}

	notification := ViolationNotification{
		Codigo:     violation.Codigo,
		Deposito:   violation.Deposito,
		Tentativas: violation.Tentativas,
		Message:    message,
	}
	payloadBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Erro ao serializar o payload: %v\n", err)
		return
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", callbackURL, body)
	if err != nil {
		log.Printf("Erro ao criar a requisição: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar a notificação: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("A notificação não foi bem-sucedida: Código de status %d\n", resp.StatusCode)
	}
}

var violations = make(map[string]*ViolationReport)

func PayloadHandler(w http.ResponseWriter, r *http.Request) {
	var payloads []map[string]interface{}

	_ = json.NewDecoder(r.Body).Decode(&payloads)

	client := utils.NewRedisClient()

	for _, payload := range payloads {
		body, ok := payload["body"].(map[string]interface{})
		if !ok {
			continue
		}

		data, ok := body["data"].(string)
		if !ok {
			continue
		}

		codigo, deposito := utils.ExtractCodigoAndDeposito(data)
		if codigo == "" {
			continue
		}

		blockedKey := "blocked:" + codigo
		blocked, _ := client.Exists(utils.Ctx, blockedKey).Result()
		if blocked > 0 {
			continue
		}

		key := "codigo:" + codigo
		tentativas, _ := client.Incr(utils.Ctx, key).Result()

		if tentativas == 1 {
			client.Expire(utils.Ctx, key, 5*time.Second)
		}

		if tentativas > 3 {
			client.Set(utils.Ctx, blockedKey, "1", 1*time.Minute)

			if _, exists := violations[codigo]; !exists {
				violations[codigo] = &ViolationReport{Codigo: codigo, Deposito: deposito, Tentativas: tentativas}
			} else {
				violations[codigo].Tentativas = tentativas
			}

			sendViolationNotification(violations[codigo])

			continue
		}

		_ = json.NewEncoder(w).Encode(payload)
	}

	for _, report := range violations {
		log.Printf("Violação: Código %s excedeu o limite com + de %d tentativas. (Depósito ID: %s)\n", report.Codigo, report.Tentativas, report.Deposito)
	}
	violations = make(map[string]*ViolationReport)

}
