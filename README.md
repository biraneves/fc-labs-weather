# 🌦️ FC Labs Weather

Aplicação em Go que recebe um CEP brasileiro de 8 dígitos, resolve a localidade via ViaCEP
e retorna as temperaturas atuais daquela cidade em graus Celsius, graus Fahrenheit e Kelvin
consumindo a WeatherAPI. O serviço está pronto para execução local com **Docker Compose** e
para deploy no **Google Cloud Run**.

## 📌 Visão Geral da API

- **Endpoint**: `GET /weather?cep={cep}`
- **Sucesso (200)**:

    Resposta: `{"temp_C":27.1,"temp_F":80.8,"temp_K":300.3}`

- **CEP com formato válido porém inexistente (404)**:

    Resposta: `{"error":"cannot find zipcode"}`

- **CEP inválido (422)**:

    Resposta: `{"error":"invalid zipcode"}`

## 🧱 Arquitetura

- `cmd/server`: ponto de entrada que carrega a configuração, instancia o cliente HTTP com timeouts,
  cria os clients do ViaCEP/WeatherAPI, monta o caso de uso `GetWeatherByCEP` e registra os
  middlewares e rotas no `http.Server`.
- `internal/application/dto`: define os DTOs usados na borda da aplicação (`RequestIn/Out`,
  `ViaCEPRequest/Response`, `WeatherAPIRequest/Response`) e garante a serialização esperada do JSON.
- `internal/application/ports`: interfaces que separam o caso de uso das infraestruturas;
  `inbound.GetWeatherByCEPUseCase` descreve a porta de entrada e
  `outbound.ZipcodeLookupPort / WeatherProviderPort` abstraem os clientes externos.
- `internal/application/usecase/get_weather_by_cep.go`: regra de negócio que valida o CEP, busca 
  localização na ViaCEP, consome a WeatherAPI e converte as temperaturas para graus Celsius, graus 
  Fahrenheit e Kelvin.
- `internal/domain/entity`: value objects do domínio (`Cep`, `TemperatureCelsius`, 
  `TemperatureFahrenheit`, `TemperatureKelvin`) com validações e conversões encapsuladas.
- `internal/infrastructure/config`: `config.go` carrega variáveis do `.env` ou ambiente aplicando
  defaults; arquivos `.env` e `.env.example` documentam os parâmetros.
- `internal/infrastructure/http/server`: camada HTTP com handler principal, rota de healthcheck,
  middlewares de logging e recovery e utilitários de resposta.
- `internal/infrastructure/http/viacep`: client REST responsável por consultar a ViaCEP e mapear 
  respostas e erros.
- `internal/infrastructure/http/weather_api`: client REST para a WeatherAPI, incluindo montagem de
  query, tratamento de status e log de diagnósticos.

## ⚙️ Variáveis de Ambiente

Copie `internal/infrastructure/config/.env.example` para `internal/infrastructure/config/.env` e
ajuste conforme necessário:

```
HTTP_TIMEOUT=5s
VIACEP_URL=https://viacep.com.br/ws/
VIACEP_TIMEOUT=5s
WEATHER_URL=https://api.weatherapi.com/v1
WEATHER_API_KEY=<sua_chave_weatherapi>
WEATHER_TIMEOUT=10s
```

> `WEATHER_API_KEY` é obrigatório em produção; demais variáveis possuem defaults seguros.

## 🛠️ Makefile

Optou-se pela criação de um `Makefile` para centralizar os comandos mais utilizados e evitar
inconsistências entre ambientes. Alvos como `make run`, `make test`, `make docker-up` e 
`make docker-down` encapsulam as instruções corretas (inclusive flags como `-race`) e facilitam
a automação local e em pipelines, sem necessidade de memorizar sequências longas de comandos.

## 🚀 Execução Local

### 1. Go direto

Após criar o arquivo `.env` com os valores adequados, copie-o para a raiz do projeto e execute:

```bash
make run
```

### 2. Docker Compose (recomendado)

```bash
export IMAGE_TAG=local
docker compose up --build
```

Ou utilize o atalho:

```bash
make docker-up
```

Para interromper a execução:

```bash
make docker-down
```

## ✅ Testes Automatizados

Executar testes com verificação de corrida:

```bash
make test
```

## 🔄 Pipeline CI/CD

- Workflow em `.github/workflows/ci.yml`.
- Gatilhos: `push` / `pull_request` em `main` e tags `v*.*.*`.
- Etapas principais:
  1. Checkout do código.
  2. Configuração go Go `1.24`.
  3. Execução de `go test -v -race ./...`.
  4. Preparação de QEMU + Buildx.
  5. Login no Docker Hub (exceto em PRs).
  6. Build multi-arch com `Dockerfile.prod`.
  7. Geração de tags (`latest`, `sha`, semver) e push para `biraneves/weather-app`.
- Resultado: toda entrega em `main` publica automaticamente uma nova imagem no Docker Hub, pronta para
  o deploy no Cloud Run.

## ☁️ Deploy no Google Cloud Run

### 1. Autenticar e definir projeto

```bash
gcloud auth login
gcloud config set project <PROJECT_ID>
```

### 2. Deploy utilizando a imagem do Docker Hub gerata pela pipeline

```bash
gcloud run deploy fc-labs-weather \
  --image=docker.io/biraneves/weather-app:latest \
  --region=us-central1 \
  --allow-unauthenticated \
  --port=8080 \
  --set-env-vars=VIACEP_URL=https://viacep.com.br/ws/,WEATHER_URL=https://api.weatherapi.com/v1,WEATHER_API_KEY=<sua_chave>
```

### 3. URL pública

Exposta pelo Cloud Run no formato `https://fc-labs-weather-<hash>-<region>.run.app/weather?cep={cep}`.

### 4. Atualizações contínuas

```bash
gcloud run deploy fc-labs-weather \
  --image=docker.io/biraneves/weather-app:<nova_tag>
```

## 🚧 Dificuldades Encontradas

Inicialmente foi utilizado o Viper para gerenciar as configurações da aplicação, mas foram encontrados
problemas de consistência com o carregamento das variáveis (especialmente no ambiente de produção).

Para simplificar o fluxo e ter controle explícito dos defaults, migrou-se para configuração manual
baseada em `godotenv` + `os.Getenv`, garantindo a previsibilidade no comportamento dos ambientes.

## 📦 Entrega

- Código-fonte completo neste repositório.
- Testes automatizados cobrindo fluxos de sucesso e falha (invalid `zipcode`, `cannot find zipcode`).
- Conteinerização com `Dockerfile.prod` e `docker-compose.yaml`.
- Pipeline GitHub Actions compilando, testando e publicando a imagem em `biraneves/weather-app`.
- Deploy no Google Cloud Run consumindo diretamente a imagem do Docker Hub, com
  [endpoint público ativo](https://weather-app-74472728841.europe-west1.run.app/weather?cep=07190050).