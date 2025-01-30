# janeiro_Desafio2_Multithreading
Desafio de Janeiro 2025 Faculdade Full Cycle - Go Expert - Multithreading

## Busca CEP Rápida
### Descrição
Busca CEP Rápida é uma aplicação desenvolvida em Go que utiliza multithreading e APIs para consultar informações de endereço a partir de um CEP (Código de Endereçamento Postal). O projeto realiza requisições simultâneas para duas APIs diferentes: BrasilAPI e ViaCEP. O objetivo é obter o resultado da API que responder mais rapidamente, descartando a resposta mais lenta. Além disso, a aplicação impõe um limite de tempo de 1 segundo para as respostas, exibindo um erro de timeout caso nenhuma das APIs responda dentro desse período.

## Funcionalidades
Consulta Simultânea: Realiza requisições paralelas para duas APIs distintas.
Seleção da Resposta Mais Rápida: Utiliza goroutines e canais para capturar a primeira resposta recebida.
Limite de Tempo: Implementa um timeout de 1 segundo para evitar esperas longas.
Exibição de Resultados: Mostra os dados do endereço no terminal, indicando qual API forneceu a informação.
Requisitos
Go versão 1.16 ou superior
Instalação

### Clone o repositório:

git clone https://github.com/willianfariabatista/seu-repositorio.git

### Acesse a pasta do projeto:

cd /janeiro_Desafio2_Multithreading

### Compile o projeto:

go build -o buscaCEP main.go

### Uso

###Execute o binário gerado para iniciar a aplicação. Por padrão, o CEP utilizado é 06341650, mas você pode alterar o valor diretamente no código ou adaptar para receber como argumento.

./buscaCEP

### Exemplo de saída:

#### Resultado recebido da BrasilAPI:

CEP: 06341-655
Logradouro: Rua João Marcos
Bairro: Vila Flora
Cidade: Barueri
UF: SP

#### Caso ocorra um timeout:

Timeout: Nenhuma resposta recebida em 1 segundo.

#### Estrutura do Projeto

/janeiro_Desafio2_Multithreading/main.go

main.go: Código principal da aplicação que gerencia as requisições concorrentes às APIs e o processamento dos resultados.

## Tecnologias Utilizadas

Go: Linguagem de programação utilizada para desenvolver a aplicação.

#### APIs:

BrasilAPI

ViaCEP

#### Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma issue ou enviar um pull request com melhorias e sugestões.
