# Sistema de Temperatura por CEP com OpenTelemetry e Zipkin

Este projeto implementa um sistema distribuído em Go que recebe um CEP, identifica a cidade e retorna o clima atual com temperaturas em Celsius, Fahrenheit e Kelvin, utilizando OpenTelemetry para observabilidade e Zipkin para visualização de traces.

## Arquitetura

O sistema é composto por:

- **Service A**: Valida CEPs e encaminha para Service B
- **Service B**: Busca cidade pelo CEP e consulta temperatura
- **OpenTelemetry Collector**: Coleta traces dos serviços
- **Zipkin**: Visualiza traces distribuídos

## Requisitos

- Docker
- Docker Compose
- Go 1.21+ (para desenvolvimento local)

## Como Executar

### 1. Configuração da API Key (Opcional)

Para usar a WeatherAPI com chave real, configure a variável de ambiente. Escolha um dos métodos:

#### Método 1: Arquivo .env (Recomendado)
```bash
# Copie o arquivo de exemplo
cp config.env.example .env

# Edite o arquivo .env e adicione sua chave real
# WEATHER_API_KEY=sua_chave_real_aqui
```


#### Método 2: Variável de ambiente no terminal
```bash
export WEATHER_API_KEY="sua_chave_aqui"
```

#### Método 3: Linha de comando direta
```bash
WEATHER_API_KEY="sua_chave_aqui" docker-compose up --build
```


Se não configurada, o sistema usará a chave demo (limitada).

### 2. Executar com Docker Compose

```bash
# Subir todos os serviços
docker-compose up --build

# Ou em background
docker-compose up -d --build
```

### 3. Verificar Serviços

- **Service A**: http://localhost:8080
- **Service B**: http://localhost:8081  
- **Zipkin UI**: http://localhost:9411

## Uso da API

### Endpoint Principal (Service A)

```bash
curl -X POST http://localhost:8080/cep-weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'
```

### Resposta de Sucesso

```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

### Códigos de Erro

- **422**: CEP inválido (não tem 8 dígitos)
- **404**: CEP não encontrado
- **500**: Erro interno do servidor

## Observabilidade

### Visualizar Traces no Zipkin

1. Acesse http://localhost:9411
2. Clique em "Run Query" para ver todos os traces
3. Clique em um trace para ver detalhes

### Traces Implementados

- **CEPValidationHandler**: Validação do CEP no Service A
- **call-service-b**: Chamada HTTP entre serviços
- **fetch-city-via-viacep**: Busca da cidade via ViaCEP
- **fetch-weather-api**: Consulta temperatura via WeatherAPI
- **temperature-conversion**: Conversão de temperaturas

## Desenvolvimento Local

### Estrutura do Projeto

```
├── service_a/          # Serviço de validação de CEP
│   ├── model/         # Structs específicas do Service A
│   ├── main.go        # Código principal
│   ├── telemetry_setup.go  # Configuração de telemetria
│   ├── go.mod         # Dependências
│   └── Dockerfile     # Build do container
├── service_b/          # Serviço de consulta de temperatura
│   ├── model/         # Structs específicas do Service B
│   ├── main.go        # Código principal
│   ├── telemetry_setup.go  # Configuração de telemetria
│   ├── go.mod         # Dependências
│   └── Dockerfile     # Build do container
├── docker-compose.yml  # Orquestração dos serviços
├── otel-collector-config.yml  # Configuração do collector
├── config.env.example  # Exemplo de configuração de ambiente
└── README.md          # Documentação do projeto
```

### Executar Serviços Individualmente

```bash
# Service A
cd service_a
go run main.go

# Service B  
cd service_b
go run main.go
```

## APIs Externas Utilizadas

- **ViaCEP**: https://viacep.com.br/ (busca cidade por CEP)
- **WeatherAPI**: https://www.weatherapi.com/ (consulta temperatura)

## Modelos de Dados

Cada serviço possui seu próprio pacote `model` local com structs específicas:

### Service A (service_a/model/):

- **`CEPRequest`**: Requisição de validação de CEP
- **`CEPValidationError`**: Erros de validação do CEP
- **`ServiceAResponse`**: Resposta padrão do Service A
- **`HTTPStatus`**: Códigos de status HTTP
- **`ErrorResponse`**: Resposta de erro padrão

### Service B (service_b/model/):

- **`CEPRequest`**: Requisição de validação de CEP
- **`ViaCEPResponse`**: Resposta completa da API ViaCEP
- **`WeatherAPIResponse`**: Resposta completa da API WeatherAPI
- **`WeatherResponse`**: Resposta final do sistema
- **`ServiceBError`**: Erros específicos do Service B
- **`TemperatureConversion`**: Funções de conversão de temperatura

### Uso nos Serviços:

```go
// Service A
import "service_a/model"
var req model.CEPRequest

// Service B
import "service_b/model"
var weatherResp model.WeatherResponse
var conversion model.TemperatureConversion
```

## Conversões de Temperatura

- **Fahrenheit**: F = C × 1.8 + 32
- **Kelvin**: K = C + 273

## Limitações

- A chave demo da WeatherAPI tem limitações de uso
- CEPs devem ter exatamente 8 dígitos
- Sistema funciona apenas com CEPs brasileiros válidos

## Testando a Aplicação

### 1. Verificar se os serviços estão rodando

```bash
# Verificar containers
docker-compose ps

# Verificar logs
docker-compose logs service-a
docker-compose logs service-b
```

### 2. Testar endpoints
```bash
# Teste 1: CEP válido
curl -X POST http://localhost:8080/cep-weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Teste 2: CEP inválido
curl -X POST http://localhost:8080/cep-weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'
```

### 3. Verificar traces no Zipkin

Acesse: http://localhost:9411

- Procure por traces dos serviços
- Verifique se os spans estão sendo criados corretamente
- Analise os tempos de resposta