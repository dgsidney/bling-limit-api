package utils

import (
	"encoding/json"
	"log"
)

func ExtractCodigoAndDeposito(data string) (string, string) {
	var result struct {
		Retorno struct {
			Estoques []struct {
				Estoque struct {
					Codigo    string `json:"codigo"`
					Depositos []struct {
						Deposito struct {
							ID string `json:"id"`
						} `json:"deposito"`
					} `json:"depositos"`
				} `json:"estoque"`
			} `json:"estoques"`
		} `json:"retorno"`
	}

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Printf("Erro ao extrair código e depósito: %v", err)
		return "", ""
	}

	if len(result.Retorno.Estoques) > 0 && len(result.Retorno.Estoques[0].Estoque.Depositos) > 0 {
		return result.Retorno.Estoques[0].Estoque.Codigo, result.Retorno.Estoques[0].Estoque.Depositos[0].Deposito.ID
	}

	return "", ""
}
