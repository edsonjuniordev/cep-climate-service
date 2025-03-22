# Cep Climate Service

Este é um sistema desenvolvido em Go que permite obter informações climáticas a partir de um CEP válido. O sistema utiliza APIs externas para determinar a localização com base no CEP e retorna a temperatura atual nas unidades Celsius, Fahrenheit e Kelvin. O sistema foi implementado com suporte para deployment no Google Cloud Run.

## Funcionalidades

- Recebe um CEP válido de 8 dígitos.
- Pesquisa a localização correspondente ao CEP usando a API ViaCEP.
- Obtém a temperatura atual da localização usando a WeatherAPI.
- Retorna a temperatura em Celsius, Fahrenheit e Kelvin.
- Responde adequadamente em cenários de erro.

## Cenários de Resposta

### Sucesso
- **Código HTTP**: `200`
- **Response Body**:
```json
  {
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.65
  }
```

### Erro: CEP inválido
- Código HTTP: 422

- Mensagem: invalid zipcode

### Erro: CEP não encontrado

- Código HTTP: 404

- Mensagem: can not find zipcode

## Pré-requisitos

- Go instalado

- Conta no Google Cloud Platform com Google Cloud Run configurado

- Chaves de API:

    - WeatherAPI

## Instalação e Configuração

1. Clone o repositório:

```bash
git clone https://github.com/edsonjuniordev/cep-climate-service.git
cd cep-climate-service
```

2. Configure as variáveis de ambiente no arquivo .env:

```bash
WEATHER_API_KEY=your_weather_api_key
```

3. Compile e execute o servidor localmente:

```bash
go run main.go
```

4. O servidor estará disponível em http://localhost:8080.

## Endpoints

### Obter clima por CEP
- URL: /weather

- Método: GET

- Query Param: 
```bash
?cep=
```

## Testes

- Execute os testes automatizados:

```bash
go test ./...
```

## Docker Compose

Para executar localmente com Docker Compose:

1. Certifique-se de que o Docker Compose está instalado.

2. Adicione a variável de ambiente no arquivo docker-compose.yaml.

3. Suba o serviço:

```bash
docker-compose up
```

4. O serviço estará disponível em http://localhost:8080.

## Cloud Run

Para acessar o projeto no Google Cloud Run, acesse:

```bash
https://cep-climate-service-m5ywbqwhhq-uc.a.run.app
```