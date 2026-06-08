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
| **Docker / Docker Compose** | Containerização e orquestração do ambiente |

---

## 🏛️ Arquitetura

O projeto adota **Arquitetura Hexagonal** (Ports & Adapters), garantindo desacoplamento entre regras de negócio e dependências externas.

### 🤔 Por que Arquitetura Hexagonal?

O Doheem possui uma característica que torna essa escolha natural: **múltiplas integrações externas simultâneas**. O sistema conversa com PostgreSQL (persistência), Redis (cache) e uma API HTTP — ou seja, três "portas" de entrada e saída diferentes.

Em uma arquitetura tradicional em camadas, a lógica de negócio ficaria acoplada a essas dependências. Isso significa que, por exemplo, a regra de *"dividir uma despesa igualmente entre os moradores"* estaria misturada com código de banco de dados ou HTTP — tornando o sistema frágil e difícil de testar.

Com a **Arquitetura Hexagonal**, o domínio de negócio fica no centro, completamente isolado:

| Problema sem a arquitetura | Como a Hexagonal resolve |
|---|---|
| Trocar o banco de dados exigiria alterar regras de negócio | O repositório é uma interface (Port); só o Adapter muda |
| Testar a lógica de divisão de despesas exige subir o banco | O domínio é testável puro, sem infraestrutura |
| Difícil adicionar novo canal (ex: WebSocket) | Basta criar um novo Adapter HTTP sem tocar no domínio |

Em termos práticos para o TCC: essa arquitetura permite que a **regra de negócio seja o coração do sistema**, independente de qual tecnologia está ao redor. Se amanhã o projeto migrar de PostgreSQL para outro banco, as regras de divisão de despesas e gestão de tarefas continuam intactas.

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
└── adapters/                    # Implementações externas (HTTP, DB, Redis)
    ├── http/                    # Handlers e rotas
    ├── repository/              # Queries SQLC / PostgreSQL
    └── cache/                   # Integração Redis
```

### Fluxo de dados

```
HTTP Request → Handler (Adapter) → Service (Port) → Domain → Repository (Adapter)
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

Isso iniciará PostgreSQL e Redis automaticamente.

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
