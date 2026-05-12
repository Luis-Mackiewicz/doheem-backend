# 🏠 Doheem — Backend API

> API de gestão para repúblicas estudantis: finanças compartilhadas, tarefas domésticas e convivência organizada.

---

## 📌 Sobre o Projeto

O **Doheem Backend** é a camada de negócios e dados da plataforma Doheem. Construído em **Go** com **Arquitetura Hexagonal**, oferece uma API REST de alta performance para gerenciar despesas compartilhadas, tarefas domésticas e autenticação de moradores de repúblicas estudantis.

---

## 🛠️ Tecnologias

| Tecnologia | Função |
|---|---|
| **Go** | Linguagem principal da API |
| **PostgreSQL** | Banco de dados relacional |
| **SQLC** | Geração de código type-safe a partir de queries SQL |
| **Redis** | Cache de sessões e dados temporários |
| **Kafka** | Mensageria assíncrona para eventos de tarefas e notificações |
| **Docker / Docker Compose** | Containerização e orquestração do ambiente |

---

## 🏛️ Arquitetura

O projeto adota **Arquitetura Hexagonal** (Ports & Adapters), garantindo desacoplamento entre regras de negócio e dependências externas.

```
cmd/
└── doheem/
    └── main.go                  # Ponto de entrada da aplicação

internal/
├── domain/                      # Entidades e regras de negócio puras
│   ├── expense/
│   ├── task/
│   └── user/
├── services/                    # Casos de uso e orquestração
│   ├── expense_service.go
│   ├── task_service.go
│   └── auth_service.go
└── adapters/                    # Implementações externas (HTTP, DB, Kafka, Redis)
    ├── http/                    # Handlers e rotas
    ├── repository/              # Queries SQLC / PostgreSQL
    ├── cache/                   # Integração Redis
    └── messaging/               # Producers e Consumers Kafka
```

### Fluxo de dados

```
HTTP Request → Handler (Adapter) → Service (Port) → Domain → Repository (Adapter)
                                                   ↓
                                           Kafka Producer → Event Bus → Consumer
```

---

## ✅ Pré-requisitos

- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados
- [Go 1.22+](https://golang.org/dl/) (para execução local sem Docker)

---

## 🚀 Como Rodar

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/doheem-backend.git
cd doheem-backend
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
# Edite o .env com suas configurações locais
```

### 3. Suba a infraestrutura com Docker Compose

```bash
docker-compose up -d
```

Isso iniciará PostgreSQL, Redis e Kafka automaticamente.

### 4. Execute a aplicação

```bash
go run cmd/doheem/main.go
```

A API estará disponível em `http://localhost:8080`.

---

## 📡 Endpoints Principais

### 🔐 Autenticação — `/auth`

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/auth/register` | Cadastro de novo morador |
| `POST` | `/auth/login` | Login e geração de token |
| `POST` | `/auth/logout` | Invalidação de sessão |
| `GET` | `/auth/me` | Dados do usuário autenticado |

### 💸 Despesas — `/expenses`

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/expenses` | Lista despesas da república |
| `POST` | `/expenses` | Cria nova despesa compartilhada |
| `PATCH` | `/expenses/:id/pay` | Marca parcela como paga |
| `GET` | `/expenses/summary` | Resumo financeiro por morador |

### ✅ Tarefas — `/tasks`

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/tasks` | Lista tarefas da república |
| `POST` | `/tasks` | Cria nova tarefa |
| `PATCH` | `/tasks/:id/status` | Atualiza status da tarefa |
| `DELETE` | `/tasks/:id` | Remove uma tarefa |

---

## ⚡ Destaque Técnico — Processamento Orientado a Eventos

O módulo de tarefas utiliza **Apache Kafka** para processamento assíncrono de eventos, garantindo alta escalabilidade e desacoplamento entre os serviços.

```
Morador conclui tarefa
       ↓
  HTTP Handler
       ↓
  Task Service → Publica evento `task.completed` no Kafka
                          ↓
               Notification Consumer
                          ↓
          Envia notificação para os outros moradores
```

Isso permite que notificações, atualizações de pontuação e logs de auditoria sejam processados de forma independente, sem impactar a latência da requisição principal.

---

## 🧪 Testes

```bash
# Rodar todos os testes
go test ./...

# Testes com cobertura
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 📄 Licença

Distribuído sob a licença MIT. Veja `LICENSE` para mais informações.

---

> Feito com ☕ e Go para facilitar a vida nas repúblicas universitárias.
