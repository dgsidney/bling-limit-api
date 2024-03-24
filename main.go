package main

import (
	"bling_limit/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	logFileName := "server.log"

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir/criar o arquivo de log: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	http.HandleFunc("/api/filter", handlers.PayloadHandler)

	fmt.Println("Servidor iniciado na porta 5000")

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
