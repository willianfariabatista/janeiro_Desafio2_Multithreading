package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Estrutura que representa a o Endereço.
type Endereco struct {
	CEP        string // Código de Endereçamento Postal.
	Logradouro string // Nome da rua ou avenida.
	Bairro     string // Nome do bairro.
	Cidade     string // Nome da cidade.
	UF         string // Unidade Federativa (Estado).
	Fonte      string // Para indicar qual API retornou mais rápido.
}

// Estrutura para decodificar a resposta da API: BrasilAPI.
type brasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

// Estrutura para decodificar a resposta da API: ViaCEP.
type viaCEPResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

// buscaBrasilAPI faz a requisição para a brasilapi.com.br.
func buscaBrasilAPI(ctx context.Context, cep string) (Endereco, error) {

	endereco := Endereco{}
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	// Cria uma requisição HTTP GET com um contexto para buscar o endereço na API brasilapi.com.br.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	// Se houver erro na criação da requisição, retorna a estrutura de endereço vazia e o erro.
	if err != nil {
		return endereco, err
	}

	// Envia a requisição HTTP para obter os dados do endereço.
	resp, err := http.DefaultClient.Do(req)

	// Se houver erro na requisição, retorna a estrutura de endereço vazia e o erro.
	if err != nil {
		return endereco, err
	}

	// Garante que o corpo da resposta HTTP seja fechado.
	defer resp.Body.Close()

	// Verifica se a resposta HTTP tem status diferente de 200 (OK). Se for diferente, retorna a estrutura de endereço vazia e uma mensagem de erro.
	if resp.StatusCode != http.StatusOK {
		return endereco, fmt.Errorf("BrasilAPI retornou status %d", resp.StatusCode)
	}

	var data brasilAPIResponse

	// Decodifica o JSON da resposta HTTP para a estrutura brasilAPIResponse. Se ocorrer um erro na decodificação, retorna a estrutura de endereço vazia e o erro.
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return endereco, err
	}

	// Preenche a estrutura Endereco com os dados retornados pela API brasilapi.com.br.
	// Atribui os valores correspondentes do JSON recebido à estrutura Endereco.
	endereco = Endereco{
		CEP:        data.Cep,
		Logradouro: data.Street,
		Bairro:     data.Neighborhood,
		Cidade:     data.City,
		UF:         data.State,
		Fonte:      "BrasilAPI",
	}

	// Retorna a estrutura Endereco preenchida.
	return endereco, nil
}

// buscaViaCEP faz a requisição para a API: viacep.com.br.
func buscaViaCEP(ctx context.Context, cep string) (Endereco, error) {

	endereco := Endereco{}
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	// Cria uma requisição HTTP GET com um contexto para buscar o endereço na API viacep.com.br.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	// Se houver erro na criação da requisição, retorna a estrutura de endereço vazia e o erro.
	if err != nil {
		return endereco, err
	}

	// Envia a requisição HTTP para obter os dados do endereço.
	resp, err := http.DefaultClient.Do(req)

	// Se houver erro na requisição, retorna a estrutura de endereço vazia e o erro.
	if err != nil {
		return endereco, err
	}

	// Garante que o corpo da resposta HTTP seja fechado.
	defer resp.Body.Close()

	// Verifica se a resposta HTTP tem status diferente de 200 (OK). Se for diferente, retorna a estrutura de endereço vazia e uma mensagem de erro.
	if resp.StatusCode != http.StatusOK {
		return endereco, fmt.Errorf("ViaCEP retornou status %d", resp.StatusCode)
	}

	var data viaCEPResponse

	// Decodifica o JSON da resposta HTTP para a estrutura viaCEPResponse. Se ocorrer um erro na decodificação, retorna a estrutura de endereço vazia e o erro.
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return endereco, err
	}

	// Preenche a estrutura Endereco com os dados retornados pela API viacep.com.br.
	// Atribui os valores correspondentes do JSON recebido à estrutura Endereco.
	endereco = Endereco{
		CEP:        data.Cep,
		Logradouro: data.Logradouro,
		Bairro:     data.Bairro,
		Cidade:     data.Localidade,
		UF:         data.Uf,
		Fonte:      "ViaCEP",
	}

	// Retorna a estrutura Endereco preenchida.
	return endereco, nil
}

func main() {

	// Defina o CEP que deseja consultar.
	cep := "06341650"

	// Cria um contexto com timeout de 1 segundo para limitar o tempo da requisição HTTP.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	// Garante que o contexto sejam liberados após a execução.
	defer cancel()

	// Canal para receber o resultado de qualquer uma das APIs.
	resultChan := make(chan Endereco, 2)

	// Armazena erros que podem ocorrer nas requisições.
	errChan := make(chan error, 2)

	// Inicia duas goroutines para buscar o endereço em APIs diferentes de forma concorrente.
	go func() {

		// Goroutine para buscar o endereço na BrasilAPI.
		end, err := buscaBrasilAPI(ctx, cep)
		if err != nil {

			// Envia o erro para o canal de erros se a requisição falhar.
			errChan <- err
			return
		}

		// Envia o resultado para o canal de resultados.
		resultChan <- end
	}()

	go func() {
		// Goroutine para buscar o endereço na ViaCEP.
		end, err := buscaViaCEP(ctx, cep)
		if err != nil {

			// Envia o erro para o canal de erros se a requisição falhar.
			errChan <- err
			return
		}

		// Envia o resultado para o canal de resultados.
		resultChan <- end
	}()

	// Select para pegar o resultado da primeira que chegar.
	select {

	// Caso o contexto expire (timeout atingido), interrompe a espera por respostas.
	case <-ctx.Done():

		// Se o contexto expirar (timeout de 1s), exibimos erro
		log.Println("Timeout: Nenhuma resposta recebida em 1 segundo.")

	// Recebe o primeiro resultado disponível no canal de respostas.
	case end := <-resultChan:

		// Se um resultado for recebido antes do timeout, cancela o contexto para interromper a outra requisição.
		cancel()

		// Registra no log os detalhes do endereço recebido da API que respondeu mais rápido.
		log.Printf("Resultado recebido da %s:\n", end.Fonte)
		log.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nCidade: %s\nUF: %s\n",
			end.CEP,
			end.Logradouro,
			end.Bairro,
			end.Cidade,
			end.UF,
		)

	// Recebe um erro do canal de erros caso a requisição falhe.
	case err := <-errChan:

		// Se um erro for recebido, cancela o contexto e exibe a mensagem de erro no log.
		cancel()
		log.Printf("Erro ao consultar: %v\n", err)
	}
}
