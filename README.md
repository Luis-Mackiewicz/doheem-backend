# Doheem — Backend API

API de gestão para repúblicas estudantis: finanças compartilhadas, tarefas domésticas e convivência organizada.

## Tecnologias

| Tecnologia | Função |
|---|---|
| **Go 1.26** | Linguagem principal da API |
| **PostgreSQL 18** | Banco de dados relacional |
| **Redis 7** | Rate limiting (sliding window) |
| **SQLC** | Geração de código type-safe a partir de queries SQL |
| **pgx v5** | Driver PostgreSQL |
| **go-redis v9** | Cliente Redis |
| **golang-jwt v5** | Autenticação JWT + Refresh Tokens |
| **go-playground/validator** | Validação de requests |
| **testcontainers-go** | Testes de integração com banco real |
| **golang-migrate** | Migrações de banco de dados |
| **Docker / Docker Compose** | Containerização da infraestrutura |

## Estrutura do Projeto

```
cmd/
  main.go                        # Ponto de entrada da aplicação

internal/
  audit_log/                     # Auditoria de ações (entity, repository)
  config/                        # Configuração via environment variables
  db/
    migrations/                  # Migrations SQL (golang-migrate)
    queries/                     # Queries SQL fonte (sqlc)
    models.go                    # Structs geradas pelo sqlc
    *.sql.go                     # Implementações geradas pelo sqlc
    sqlc.yaml                    # Configuração do sqlc
  dbtest/                        # Helpers para testes com testcontainers
  expense/                       # Despesas (entity, service, repository, tests)
  group/                         # Grupos (entity, service, repository, tests)
  http/                          # Camada HTTP (handlers, middleware, router)
  notification/                  # Notificações (entity, service, repository, tests)
  task/                          # Tarefas (entity, service, repository, tests)
  user/                          # Usuários (entity, service, repository, tests)
```

Cada domínio segue o mesmo padrão:
- `entity.go` — structs de domínio + interfaces de repositório + erros
- `service.go` — lógica de negócio (casos de uso)
- `repository.go` — implementação concreta do repositório (usa `db.Queries` do sqlc)
- `service_test.go` / `repository_test.go` — testes

## Pré-requisitos

- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados
- [Go 1.26+](https://golang.org/dl/) (para execução local sem Docker)
- [golang-migrate](https://github.com/golang-migrate/migrate) (CLI para migrações)

## Como Rodar

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

Variáveis disponíveis (com defaults para desenvolvimento):

| Variável | Default | Descrição |
|---|---|---|
| `DATABASE_URL` | `postgres://doheem_dev_user:simple_pswd@localhost:5432/doheem_dev_db` | Conexão PostgreSQL |
| `REDIS_URL` | `redis://localhost:6379/0` | Conexão Redis |
| `JWT_SECRET` | `doheem-dev-secret-change-in-production` | Chave JWT |
| `JWT_EXPIRES_IN` | `24h` | Duração do token JWT |
| `JWT_REFRESH_EXPIRES_IN` | `168h` (7 dias) | Duração do refresh token |
| `PORT` | `8080` | Porta do servidor HTTP |
| `APP_ENV` | `development` | Ambiente (`development` / `production`) |
| `LOG_FORMAT` | `text` | Formato de log (`text` / `json`) |

### 3. Suba a infraestrutura com Docker Compose

```bash
docker compose up -d
```

Isso iniciará PostgreSQL e Redis automaticamente.

### 4. Execute as migrações

```bash
make migrate-up
```

### 5. Inicie a aplicação

```bash
go run cmd/main.go
```

A API estará disponível em `http://localhost:8080`.

### Docker (produção)

```bash
docker build -t doheem-server .
docker run -p 8080:8080 --env-file .env doheem-server
```

## Endpoints

### Health Check

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/` | Health check (verifica PostgreSQL e Redis) |

### Autenticação — `/api/auth`

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/api/auth/register` | Cadastro de novo morador |
| `POST` | `/api/auth/login` | Login e geração de token |
| `POST` | `/api/auth/refresh` | Renovar access token via refresh token |
| `POST` | `/api/auth/logout` | Invalidação de sessão |

Rotas de auth têm rate limit de **10 requisições por minuto**.

### Usuário — `/api/users`

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/users/me` | Dados do usuário autenticado |
| `PUT` | `/api/users/me` | Atualizar perfil |
| `PUT` | `/api/users/me/password` | Alterar senha |

### Grupos — `/api/groups`

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/api/groups` | Criar república |
| `GET` | `/api/groups` | Listar repúblicas do usuário |
| `GET` | `/api/groups/{id}` | Detalhes da república |
| `PUT` | `/api/groups/{id}` | Atualizar república |
| `GET` | `/api/groups/{id}/members` | Listar membros |
| `POST` | `/api/groups/{id}/members` | Adicionar membro |
| `PUT` | `/api/groups/{id}/members/{userId}` | Alterar cargo do membro |
| `DELETE` | `/api/groups/{id}/members/{userId}` | Remover membro |
| `POST` | `/api/groups/{id}/join` | Entrar na república via invite |
| `POST` | `/api/groups/{id}/regenerate-invite` | Regenerar token de convite |
| `GET` | `/api/groups/{id}/invite-token` | Obter token de convite |

### Despesas — `/api/expenses`

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/api/groups/{groupId}/expenses` | Criar despesa no grupo |
| `GET` | `/api/groups/{groupId}/expenses` | Listar despesas do grupo |
| `GET` | `/api/expenses/{id}` | Detalhes da despesa |
| `PUT` | `/api/expenses/{id}` | Atualizar despesa |
| `DELETE` | `/api/expenses/{id}` | Remover despesa |
| `GET` | `/api/expenses/{id}/splits` | Listar rateios da despesa |
| `PATCH` | `/api/expenses/splits/{id}/pay` | Marcar parcela como paga |
| `GET` | `/api/expenses/{id}/installments` | Listar parcelas de despesa |

### Categorias de Despesa — `/api/categories`

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/api/categories` | Criar categoria |
| `GET` | `/api/categories` | Listar categorias |
| `PUT` | `/api/categories/{id}` | Atualizar categoria |
| `DELETE` | `/api/categories/{id}` | Remover categoria |

### Tarefas — `/api/tasks`

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/api/groups/{groupId}/tasks` | Criar tarefa no grupo |
| `GET` | `/api/groups/{groupId}/tasks` | Listar tarefas do grupo |
| `GET` | `/api/tasks/{id}` | Detalhes da tarefa |
| `PUT` | `/api/tasks/{id}` | Atualizar tarefa |
| `DELETE` | `/api/tasks/{id}` | Remover tarefa |
| `GET` | `/api/tasks/{id}/occurrences` | Listar ocorrências da tarefa |
| `POST` | `/api/tasks/{taskId}/occurrences` | Criar ocorrência |
| `PATCH` | `/api/tasks/occurrences/{id}/complete` | Completar ocorrência |
| `PATCH` | `/api/tasks/occurrences/{id}/discard` | Descartar ocorrência |

### Notificações — `/api/notifications`

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/api/notifications` | Criar notificação |
| `GET` | `/api/notifications` | Listar notificações |
| `GET` | `/api/notifications/unread` | Listar não lidas |
| `PATCH` | `/api/notifications/{id}/read` | Marcar como lida |
| `PATCH` | `/api/notifications/read-all` | Marcar todas como lidas |
| `DELETE` | `/api/notifications/{id}` | Remover notificação |

## Middleware

A cadeia de middleware é aplicada na seguinte ordem (de fora para dentro):

1. **Recovery** — Captura panics e retorna 500
2. **CORS** — Permite qualquer origem (`*`)
3. **Logging** — Gera `request_id` (UUID) e loga método/path/status/duration
4. **Rate Limit** — 100 req/min global; 10 req/min para rotas de auth
5. **Auth (JWT)** — Valida token Bearer e injeta `user_id` no contexto

## Testes

```bash
# Rodar todos os testes
go test ./...

# Testes com cobertura
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Testes de integração usam **testcontainers** para subir PostgreSQL real em container.

## Migrações

```bash
make migrate-up                  # Aplicar migrações
make migrate-down N=1            # Reverter N migrações
make migrate-create NAME=desc    # Criar nova migração
make migrate-version             # Versão atual
```

## CI/CD

O pipeline no GitHub Actions executa:

1. **lint** — `golangci-lint`
2. **test** — `go vet` + `go test -count=1 -timeout=300s`
3. **build** (apenas na `main`) — compila binário, constrói imagem Docker e publica no GitHub Container Registry (`ghcr.io`)

## Licença

Distribuído sob a licença MIT. Veja `LICENSE` para mais informações.
