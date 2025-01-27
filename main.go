package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Endereco unifica os campos principais que queremos exibir:
type Endereco struct {
	CEP        string
	Logradouro string
	Bairro     string
	Cidade     string
	UF         string
	Fonte      string // Para indicar qual API retornou mais rápido
}

// Estrutura para decodificar a resposta da BrasilAPI
// (Exemplo de retorno: https://brasilapi.com.br/api/cep/v1/01153000)
type brasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

// Estrutura para decodificar a resposta da ViaCEP
// (Exemplo de retorno: https://viacep.com.br/ws/01153000/json/)
type viaCEPResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

// buscaBrasilAPI faz a requisição para a BrasilAPI (https://brasilapi.com.br/api/cep/v1/01153000)
func buscaBrasilAPI(ctx context.Context, cep string) (Endereco, error) {
	endereco := Endereco{}
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	// Cria a requisição com o context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return endereco, err
	}

	// Executa a chamada
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return endereco, err
	}
	defer resp.Body.Close()

	// Se a API não retornar 200, consideramos erro
	if resp.StatusCode != http.StatusOK {
		return endereco, fmt.Errorf("BrasilAPI retornou status %d", resp.StatusCode)
	}

	// Decodifica o JSON específico da BrasilAPI
	var data brasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return endereco, err
	}

	// Mapeia para nossa struct unificada Endereco
	endereco = Endereco{
		CEP:        data.Cep,
		Logradouro: data.Street,
		Bairro:     data.Neighborhood,
		Cidade:     data.City,
		UF:         data.State,
		Fonte:      "BrasilAPI",
	}

	return endereco, nil
}

// buscaViaCEP faz a requisição para a ViaCEP (http://viacep.com.br/ws/01153000/json/)
func buscaViaCEP(ctx context.Context, cep string) (Endereco, error) {
	endereco := Endereco{}
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	// Cria a requisição com o context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return endereco, err
	}

	// Executa a chamada
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return endereco, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return endereco, fmt.Errorf("ViaCEP retornou status %d", resp.StatusCode)
	}

	// Decodifica o JSON específico da ViaCEP
	var data viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return endereco, err
	}

	endereco = Endereco{
		CEP:        data.Cep,
		Logradouro: data.Logradouro,
		Bairro:     data.Bairro,
		Cidade:     data.Localidade,
		UF:         data.Uf,
		Fonte:      "ViaCEP",
	}

	return endereco, nil
}

func main() {
	// Defina o CEP que deseja consultar
	cep := "06341650" //01153000

	// Criamos um contexto com timeout de 1 segundo
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Canal para receber o resultado de qualquer uma das APIs
	resultChan := make(chan Endereco, 2)
	errChan := make(chan error, 2)

	// Inicia duas goroutines, cada uma chamando uma API
	go func() {
		end, err := buscaBrasilAPI(ctx, cep)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- end
	}()

	go func() {
		end, err := buscaViaCEP(ctx, cep)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- end
	}()

	// Usamos select para pegar o resultado da primeira que chegar
	select {
	case <-ctx.Done():
		// Se o contexto expirar (timeout de 1s), exibimos erro
		log.Println("Timeout: Nenhuma resposta recebida em 1 segundo.")
	case end := <-resultChan:
		// Recebemos o primeiro resultado, então cancelamos para descartar o outro
		cancel()
		log.Printf("Resultado recebido da %s:\n", end.Fonte)
		log.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nCidade: %s\nUF: %s\n",
			end.CEP,
			end.Logradouro,
			end.Bairro,
			end.Cidade,
			end.UF,
		)
	case err := <-errChan:
		// Caso a primeira coisa que chegar seja um erro, mostramos
		// e, em seguida, poderíamos até esperar a próxima.
		cancel()
		log.Printf("Erro ao consultar: %v\n", err)
	}
}
