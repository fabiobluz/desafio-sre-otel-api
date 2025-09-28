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
├── test_endpoints.sh   # Script de teste dos endpoints
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

## Troubleshooting

### Serviços não sobem

```bash
# Verificar logs
docker-compose logs

# Rebuild completo
docker-compose down
docker-compose up --build --force-recreate
```

### Erro "unable to get image" ou problemas de download

```bash
# Limpar cache do Docker
docker system prune -a

# Forçar download das imagens
docker-compose pull

# Executar novamente
docker-compose up --build
```

### Erro de versão do Go no Docker

Se receber erro `go.mod requires go >= 1.23.0 (running go 1.21.13)`:

```bash
# Limpar cache do Docker
docker system prune -a

# Rebuild forçado
docker-compose down
docker-compose up --build --force-recreate
```

**Nota**: Os Dockerfiles foram atualizados para usar Go 1.23, mas os go.mod usam Go 1.21 para compatibilidade.

### Erro "no such file or directory" do pacote pkg

Se receber erro `reading /pkg/go.mod: open /pkg/go.mod: no such file or directory`:

```bash
# Limpar cache do Docker
docker system prune -a

# Rebuild forçado
docker-compose down
docker-compose up --build --force-recreate
```

**Nota**: Os Dockerfiles foram corrigidos para copiar o pacote `pkg` corretamente durante o build.

### Erro "traces export: exporter export timeout"

Este erro ocorre quando os serviços não conseguem se conectar ao OpenTelemetry Collector.

**Soluções:**
1. **Verificar rede Docker:**
   ```bash
   docker network ls
   docker network inspect desafio-sre-otel-api_app-network
   ```

2. **Verificar se o collector está rodando:**
   ```bash
   docker-compose logs otel-collector
   ```

3. **Reiniciar os serviços:**
   ```bash
   docker-compose down
   docker-compose up --build
   ```

4. **Verificar conectividade:**
   ```bash
   # Testar se o collector está acessível
   docker exec -it service-a ping otel-collector
   ```

### Erro "go.sum not found" no Docker

Se receber erro `"/go.sum": not found`:

```bash
# 1. Gerar go.sum se não existir
cd service_a && go mod tidy
cd ../service_b && go mod tidy

# 2. Limpar cache do Docker
docker system prune -a

# 3. Rebuild forçado
docker-compose down
docker-compose up --build --force-recreate
```

**Nota**: Os Dockerfiles foram atualizados para usar o contexto correto e copiar arquivos do diretório raiz.

### Erro "go mod download" falha no Docker

Se receber erro `process "/bin/sh -c go mod download" did not complete successfully: exit code: 1`:

```bash
# 1. Verificar se as dependências estão corretas localmente
cd service_a && go mod tidy && go mod download
cd ../service_b && go mod tidy && go mod download

# 2. Limpar cache do Docker completamente
docker system prune -a --volumes

# 3. Rebuild forçado
docker-compose down
docker-compose up --build --force-recreate
```

**Nota**: Os Dockerfiles foram simplificados para copiar todos os arquivos de uma vez e configurar GOPROXY corretamente.

### Erro "pkg/go.mod not found" no Docker

Se receber erro `reading /pkg/go.mod: open /pkg/go.mod: no such file or directory`:

```bash
# 1. Verificar se o pacote pkg está correto
cd pkg && go mod tidy

# 2. Limpar cache do Docker completamente
docker system prune -a --volumes

# 3. Rebuild forçado
docker-compose down
docker-compose up --build --force-recreate
```

**Nota**: Os Dockerfiles foram ajustados para copiar o pacote `pkg` antes dos serviços e padronizar as versões do Go.

### Problemas específicos do Windows

```powershell
# 1. Verificar se o Docker Desktop está rodando
# Inicie o Docker Desktop primeiro

# 2. Verificar se o Docker está funcionando
docker --version
docker-compose --version

# 3. Se Docker Desktop não iniciar, reinicie o Windows ou:
# - Abra o Gerenciador de Tarefas
# - Finalize todos os processos do Docker
# - Reinicie o Docker Desktop

# 4. Limpar tudo e tentar novamente
docker-compose down -v
docker system prune -a
docker-compose up --build
```

### Docker Desktop não inicia

1. **Reiniciar Docker Desktop:**
   - Clique com botão direito no ícone do Docker na bandeja do sistema
   - Selecione "Restart Docker Desktop"

2. **Reiniciar serviços do Windows:**
   ```powershell
   # Como administrador
   net stop com.docker.service
   net start com.docker.service
   ```


### Traces não aparecem no Zipkin

1. Verifique se o OpenTelemetry Collector está rodando
2. Confirme se os serviços estão enviando traces
3. Aguarde alguns segundos para o batch export

### Erro de API Key

Se receber erro 401 da WeatherAPI, configure a variável `WEATHER_API_KEY` ou use a chave demo.

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

**Usando o script de teste:**
```bash
# Executar script de teste
./test_endpoints.sh
```

**Ou manualmente com curl:**
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
