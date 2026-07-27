# Clima por CEP — Google Cloud Run

Serviço HTTP em Go que recebe um CEP, identifica a cidade correspondente e
retorna a temperatura atual em Celsius, Fahrenheit e Kelvin.

## 🔗 URL em produção (Google Cloud Run)

```
https://climacep-903785298812.us-central1.run.app
```

Exemplo de uso, já publicado:

```bash
curl "https://climacep-903785298812.us-central1.run.app/?cep=01001000"
```

## Contrato da API

`GET /?cep=<8 dígitos>`

### Sucesso — `200 OK`

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

### Falhas

| Cenário                | Condição                                         | Status | Corpo                  |
|-------------------------|---------------------------------------------------|--------|--------------------------|
| Formato inválido        | CEP sem 8 dígitos ou com caracteres não numéricos  | `422`  | `invalid zipcode`       |
| CEP não encontrado      | Formato correto, mas CEP inexistente                | `404`  | `can not find zipcode`  |

> Nota sobre a fórmula de Kelvin: o enunciado do desafio define a fórmula
> como `K = C + 273` (sem o `.15`), que é a fórmula implementada aqui. O
> exemplo de payload do próprio enunciado (`301.65` para `28.5°C`) é
> inconsistente com essa fórmula — usar `+273` dá `301.5`. Optou-se por
> seguir a fórmula declarada literalmente.

## Arquitetura

```
.
├── cmd/server/main.go               # entrypoint HTTP, lê WEATHER_API_KEY e PORT
├── internal/
│   ├── cep/                         # cliente ViaCEP (CEP -> cidade)
│   ├── weather/                     # cliente WeatherAPI (cidade -> temp_C)
│   ├── temperature/                 # conversões C -> F/K (funções puras)
│   └── handler/                     # handler HTTP: validação + orquestração + contrato de erros
├── Dockerfile
└── .env.example
```

`cep.Client` e `weather.Client` são interfaces — o `WeatherHandler` depende
apenas delas, o que permite testar toda a lógica HTTP com implementações
falsas, sem chamadas de rede reais.

## Configuração

| Variável          | Descrição                                                             | Obrigatória |
|--------------------|--------------------------------------------------------------------------|-------------|
| `WEATHER_API_KEY` | Chave de API da [WeatherAPI](https://www.weatherapi.com/) (free tier)    | sim         |
| `PORT`            | Porta HTTP do servidor (o Cloud Run injeta essa variável automaticamente) | não (padrão `8080`) |

Copie `.env.example` para `.env` e preencha sua chave da WeatherAPI.

## Como rodar localmente com Docker

```bash
docker build -t climacep .
docker run --rm -p 8080:8080 --env-file .env climacep
```

```bash
curl "http://localhost:8080/?cep=01001000"
```

## Como rodar os testes

Sem dependências externas (só a biblioteca padrão do Go):

```bash
go test ./...
```

Os testes cobrem:

- Conversões de temperatura (`internal/temperature`);
- O cliente ViaCEP, incluindo o caso de CEP não encontrado (`internal/cep`,
  usando `httptest.Server`);
- O cliente WeatherAPI, incluindo o caso de cidade não encontrada
  (`internal/weather`, usando `httptest.Server`);
- O handler HTTP completo — formato inválido (422), CEP não encontrado
  (404), cidade não encontrada no provedor de clima (404), sucesso (200) e
  erro inesperado (500) — usando implementações falsas de `cep.Client` e
  `weather.Client` (`internal/handler`).

## Deploy no Google Cloud Run

Pré-requisitos: [gcloud CLI](https://cloud.google.com/sdk/docs/install)
autenticado (`gcloud auth login`) e um projeto GCP com a API do Cloud Run
habilitada.

```bash
gcloud run deploy climacep \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=<sua-chave-da-weatherapi>
```

O comando acima builda a imagem a partir do `Dockerfile` (via Cloud Build)
e publica o serviço. Ao final, o próprio `gcloud` imprime a URL pública —
copie-a para a seção **URL em produção** deste README.
