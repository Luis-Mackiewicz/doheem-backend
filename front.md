# Doheem — Frontend Angular 19+ (PWA Mobile-First)

## Stack Obrigatória

| Ferramenta | Versão | Uso |
|---|---|---|
| Angular | 19+ standalone components | Framework |
| TypeScript | 5.7+ | Linguagem |
| Tailwind CSS | 4+ | Estilização utilitária |
| Angular PWA | `@angular/pwa` schematics | Service worker + manifest |
| Angular CDK | 19+ | Drag & Drop (Kanban) |
| RxJS | 7+ | Reactive streams |
| Node.js | 22+ | Runtime |

---

## 1. Setup Inicial

```bash
ng new doheem-front --routing --style=css --ssr=false
cd doheem-front
ng add @angular/pwa
npm install -D tailwindcss @tailwindcss/postcss
npm install @angular/cdk
```

### Configurar Tailwind
`postcss.config.js`:
```js
export default { plugins: { '@tailwindcss/postcss': {} } }
```

`src/styles.css`:
```css
@import "tailwindcss";
```

### tsconfig.json (strict)
```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "forceConsistentCasingInFileNames": true,
    "baseUrl": "src",
    "paths": { "@/*": ["app/*"] }
  }
}
```

---

## 2. Arquitetura de Pastas

```
src/
├── main.ts
├── index.html
├── manifest.webmanifest
├── ngsw-config.json
├── styles.css
├── app/
│   ├── app.component.ts
│   ├── app.config.ts
│   ├── app.routes.ts
│   ├── app.components.ts             # layout: bottom-nav + router-outlet
│   │
│   ├── core/
│   │   ├── guards/
│   │   │   └── auth.guard.ts         # canActivate: redireciona para /login se sem token
│   │   ├── interceptors/
│   │   │   ├── auth.interceptor.ts   # anexa Authorization: Bearer <token>
│   │   │   └── refresh.interceptor.ts# em 401, tenta POST /api/auth/refresh + retry
│   │   ├── services/
│   │   │   ├── auth.service.ts       # login, register, refresh, logout, token storage
│   │   │   ├── api.service.ts        # HttpClient wrapper base
│   │   │   ├── group.service.ts
│   │   │   ├── expense.service.ts
│   │   │   ├── payment.service.ts
│   │   │   ├── task.service.ts
│   │   │   ├── invite.service.ts
│   │   │   ├── notification.service.ts
│   │   │   ├── split-tag.service.ts
│   │   │   └── theme.service.ts      # dark/light mode toggle
│   │   └── models/
│   │       ├── user.model.ts
│   │       ├── group.model.ts
│   │       ├── expense.model.ts
│   │       ├── payment.model.ts
│   │       ├── task.model.ts
│   │       ├── invite.model.ts
│   │       ├── notification.model.ts
│   │       ├── split-tag.model.ts
│   │       └── paginated.model.ts    # { data: T[], total: number }
│   │
│   ├── pages/
│   │   ├── onboarding/
│   │   │   ├── onboarding.component.ts      # carrossel 3 slides
│   │   │   └── onboarding.component.html
│   │   ├── auth/
│   │   │   ├── login.component.ts
│   │   │   ├── login.component.html
│   │   │   ├── register.component.ts
│   │   │   └── register.component.html
│   │   ├── dashboard/
│   │   │   ├── dashboard.component.ts       # resumo financeiro + últimas ops
│   │   │   └── dashboard.component.html
│   │   ├── groups/
│   │   │   ├── group-list.component.ts
│   │   │   ├── group-list.component.html
│   │   │   ├── group-settings.component.ts  # admin: gerenciar moradores
│   │   │   ├── group-settings.component.html
│   │   │   ├── group-create.component.ts
│   │   │   └── group-create.component.html
│   │   ├── expenses/
│   │   │   ├── expense-list.component.ts    # lista cronológica
│   │   │   ├── expense-list.component.html
│   │   │   ├── expense-form.component.ts    # formulário rápido (<30s)
│   │   │   ├── expense-form.component.html
│   │   │   ├── expense-detail.component.ts  # detalhes + splits
│   │   │   └── expense-detail.component.html
│   │   ├── tasks/
│   │   │   ├── task-kanban.component.ts     # quadro Kanban com CDK drag & drop
│   │   │   ├── task-kanban.component.html
│   │   │   ├── task-form.component.ts
│   │   │   └── task-form.component.html
│   │   ├── payments/
│   │   │   ├── payment-balance.component.ts # matriz "quem deve para quem"
│   │   │   ├── payment-balance.component.html
│   │   │   └── payment-form.component.ts
│   │   │   └── payment-form.component.html
│   │   ├── invites/
│   │   │   ├── invite-list.component.ts
│   │   │   └── invite-list.component.html
│   │   ├── notifications/
│   │   │   ├── notification-list.component.ts
│   │   │   └── notification-list.component.html
│   │   └── profile/
│   │       ├── profile.component.ts
│   │       ├── profile.component.html
│   │       ├── profile-edit.component.ts
│   │       └── profile-edit.component.html
│   │
│   └── shared/
│       ├── components/
│       │   ├── fab/
│       │   │   ├── fab.component.ts         # Floating Action Button
│       │   │   └── fab.component.html
│       │   ├── skeleton/
│       │   │   ├── skeleton.component.ts    # placeholder de carregamento
│       │   │   └── skeleton.component.html
│       │   ├── avatar-group/
│       │   │   ├── avatar-group.component.ts# fotos sobrepostas dos moradores
│       │   │   └── avatar-group.component.html
│       │   ├── pagination/
│       │   │   ├── pagination.component.ts  # controles + "X de Y"
│       │   │   └── pagination.component.html
│       │   ├── empty-state/
│       │   │   ├── empty-state.component.ts # ilustração + mensagem
│       │   │   └── empty-state.component.html
│       │   ├── confirm-dialog/
│       │   │   ├── confirm-dialog.component.ts
│       │   │   └── confirm-dialog.component.html
│       │   └── bottom-nav/
│       │       ├── bottom-nav.component.ts  # navegação inferior 5 ícones
│       │       └── bottom-nav.component.html
│       ├── pipes/
│       │   ├── currency-mask.pipe.ts        # R$ 1.234,56
│       │   ├── date-mask.pipe.ts            # 12 jan 2026
│       │   └── time-ago.pipe.ts             # "há 2 horas"
│       └── theme/
│           ├── theme.service.ts             # alterna classe .dark no <html>
│           └── theme-toggle.component.ts    # botão sol/lua
```

---

## 3. Models (Interfaces TypeScript)

### `core/models/paginated.model.ts`
```typescript
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
}
```

### `core/models/user.model.ts`
```typescript
export interface User {
  id: string;
  name: string;
  email: string;
  avatar_url?: string | null;
  created_at: string;   // ISO 8601
  updated_at: string;   // ISO 8601
}

export interface AuthResponse {
  user: User;
  token: string;
  // refresh_token está em cookie httpOnly (não vem no body)
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;   // min 6 chars
  avatar_url?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface UpdateProfileRequest {
  name?: string;
  email?: string;
  avatar_url?: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;  // min 6 chars
}
```

### `core/models/group.model.ts`
```typescript
export interface Group {
  id: string;
  name: string;
  currency: string;     // "BRL"
  is_active: boolean;
  inactive_since?: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateGroupRequest {
  name: string;
  currency: string;     // "BRL", "USD", etc
}

export interface GroupMember {
  id: string;
  group_id: string;
  user_id: string;
  role: string;         // "owner" | "admin" | "member"
  is_active: boolean;
  joined_at: string;
  user_name: string;
  user_email: string;
  avatar_url?: string | null;
}

export interface AddMemberRequest {
  user_id: string;
  role: string;         // "admin" | "member"
}
```

### `core/models/expense.model.ts`
```typescript
export interface Expense {
  id: string;
  group_id: string;
  created_by: string;
  description: string;
  total_amount: number;
  expense_date: string;   // "2006-01-02"
  due_date?: string | null;
  category_id?: string | null;
  split_type: string;     // "equal" | "custom"
  is_installment: boolean;
  installment_count?: number | null;
  created_at: string;
  updated_at: string;
}

export interface CreateExpenseRequest {
  description: string;
  total_amount: number;
  expense_date: string;       // "YYYY-MM-DD"
  due_date?: string;
  category_id?: string;
  split_type: string;         // "equal" | "custom"
  is_installment?: boolean;
  installment_count?: number;
  participants: string[];     // user_ids
  custom_splits?: { user_id: string; amount: number }[];
}

export interface ExpenseSplit {
  id: string;
  expense_id: string;
  user_id: string;
  amount: number;
  is_paid: boolean;
  created_at: string;
  user_name: string;
  user_email: string;
}

export interface ExpenseCategory {
  id: string;
  group_id: string;
  name: string;
  created_at: string;
}
```

### `core/models/payment.model.ts`
```typescript
export interface Payment {
  id: string;
  group_id: string;
  payer_id: string;
  receiver_id: string;
  amount: number;
  payment_date: string;     // "2006-01-02"
  status: string;           // "pending" | "confirmed" | "cancelled"
  notes?: string | null;
  created_at: string;
  confirmed_at?: string | null;
  cancelled_at?: string | null;
  payer_name: string;
  receiver_name: string;
}

export interface CreatePaymentRequest {
  payer_id: string;
  receiver_id: string;
  amount: number;
  payment_date: string;     // "YYYY-MM-DD"
  notes?: string;
}

export interface PaymentAttachment {
  id: string;
  payment_id: string;
  file_path: string;
  file_type: string;
  file_size: number;
  uploaded_at: string;
}
```

### `core/models/task.model.ts`
```typescript
export interface Task {
  id: string;
  group_id: string;
  title: string;
  description?: string | null;
  assigned_to: string;
  category?: string | null;
  is_recurring: boolean;
  recurring_period?: string | null;  // "daily" | "weekly" | "monthly"
  recurring_ended_at?: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTaskRequest {
  title: string;
  description?: string;
  assigned_to: string;
  category?: string;
  is_recurring?: boolean;
  recurring_period?: string;
}

export interface TaskOccurrence {
  id: string;
  task_id: string;
  due_date: string;         // "2006-01-02"
  status: string;           // "pending" | "completed" | "discarded"
  completed_by?: string | null;
  created_at: string;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  assigned_to?: string;
  category?: string;
}
```

### `core/models/invite.model.ts`
```typescript
export interface Invite {
  id: string;
  group_id: string;
  code: string;
  created_by: string;
  expires_at: string;
  used_at?: string | null;
  revoked_at?: string | null;
  created_at: string;
  created_by_name: string;
  group_name?: string;
}

export interface CreateInviteRequest {
  expires_in_hours?: number;  // default 48
}
```

### `core/models/notification.model.ts`
```typescript
export interface Notification {
  id: string;
  user_id: string;
  group_id?: string | null;
  type: string;
  title: string;
  message: string;
  is_read: boolean;
  created_at: string;
}
```

### `core/models/split-tag.model.ts`
```typescript
export interface SplitTag {
  id: string;
  group_id: string;
  name: string;
  created_by: string;
  created_at: string;
}

export interface SplitTagMember {
  id: string;
  split_tag_id: string;
  user_id: string;
  user_name: string;
}
```

---

## 4. API — Todas as Rotas

Base URL: `http://localhost:8080`

### 4.1 Saúde
```
GET  /api                       → { status, checks: { database, redis, kafka } }
```

### 4.2 Autenticação
```
POST /api/auth/register         → RegisterRequest    → AuthResponse   + Set-Cookie RefreshToken (httpOnly)
POST /api/auth/login            → LoginRequest       → AuthResponse   + Set-Cookie RefreshToken (httpOnly)
POST /api/auth/refresh          → (cookie automático) → AuthResponse  + Set-Cookie RefreshToken (httpOnly)
POST /api/auth/logout           → (cookie automático) → 204           + Clear-Cookie RefreshToken (httpOnly)
```

**Fluxo de Auth:**
1. `login()` → salva `token` do body em `localStorage`
2. `AuthInterceptor` → lê `token` do storage, anexa `Authorization: Bearer <token>`
3. Se `401`: `RefreshInterceptor` tenta `POST /api/auth/refresh` (cookie enviado automático)
4. Se refresh OK → atualiza `token` no storage → retry a request original
5. Se refresh falha → limpa storage → redireciona `/login`

### 4.3 Usuário (autenticado)
```
GET    /api/users/me             → User
PUT    /api/users/me             → UpdateProfileRequest → User
PUT    /api/users/me/password    → ChangePasswordRequest → 204
```

### 4.4 Grupos
```
POST   /api/groups                    → CreateGroupRequest → Group
GET    /api/groups?limit=&offset=     → PaginatedResponse<Group>
GET    /api/groups/{id}               → Group
PUT    /api/groups/{id}               → CreateGroupRequest → Group
DELETE /api/groups/{id}               → 204
PATCH  /api/groups/{id}/deactivate    → 204
PATCH  /api/groups/{id}/activate      → 204
GET    /api/groups/{id}/members?limit=&offset= → PaginatedResponse<GroupMember>
POST   /api/groups/{id}/members       → AddMemberRequest → 204
PUT    /api/groups/{id}/members/{userId} → { role: string } → 204
DELETE /api/groups/{id}/members/{userId} → 204
```

### 4.5 Despesas
```
POST   /api/groups/{groupId}/expenses               → CreateExpenseRequest → Expense
GET    /api/groups/{groupId}/expenses?limit=&offset= → PaginatedResponse<Expense>
GET    /api/expenses/{id}                            → Expense
PUT    /api/expenses/{id}                            → CreateExpenseRequest → Expense
DELETE /api/expenses/{id}                            → 204
GET    /api/expenses/{id}/splits?limit=&offset=      → PaginatedResponse<ExpenseSplit>
PATCH  /api/expenses/splits/{id}/pay                 → 204
GET    /api/expenses/{id}/installments?limit=&offset=→ PaginatedResponse<Installment>
PATCH  /api/expenses/installments/{id}/pay           → 204
POST   /api/groups/{groupId}/categories              → { name: string } → ExpenseCategory
GET    /api/groups/{groupId}/categories?limit=&offset=→ PaginatedResponse<ExpenseCategory>
PUT    /api/categories/{id}                          → { group_id, name } → ExpenseCategory
DELETE /api/categories/{id}                          → 204
```

### 4.6 Pagamentos
```
POST   /api/groups/{groupId}/payments               → CreatePaymentRequest → Payment
GET    /api/groups/{groupId}/payments?limit=&offset= → PaginatedResponse<Payment>
GET    /api/payments/{id}                            → Payment
PATCH  /api/payments/{id}/confirm                    → 204 (publica evento Kafka)
PATCH  /api/payments/{id}/cancel                     → 204
DELETE /api/payments/{id}                            → 204
POST   /api/payments/{id}/attachments               → FormData (file) → PaymentAttachment
```

### 4.7 Tarefas
```
POST   /api/groups/{groupId}/tasks                   → CreateTaskRequest → Task
GET    /api/groups/{groupId}/tasks?limit=&offset=    → PaginatedResponse<Task>
GET    /api/tasks/{id}                                → Task
PUT    /api/tasks/{id}                                → UpdateTaskRequest → Task
DELETE /api/tasks/{id}                                → 204
GET    /api/tasks/{id}/occurrences?limit=&offset=    → PaginatedResponse<TaskOccurrence>
POST   /api/tasks/{taskId}/occurrences                → { due_date: "2026-06-01" } → TaskOccurrence
PATCH  /api/tasks/occurrences/{id}/complete           → 204
PATCH  /api/tasks/occurrences/{id}/discard            → 204
```

### 4.8 Invites
```
POST   /api/groups/{groupId}/invites                  → CreateInviteRequest → Invite
GET    /api/groups/{groupId}/invites?limit=&offset=   → PaginatedResponse<Invite>
GET    /api/invites/pending?limit=&offset=            → PaginatedResponse<Invite>
POST   /api/invites/{id}/accept                        → 204
PATCH  /api/invites/{id}/revoke                        → 204
```

### 4.9 Notificações
```
GET    /api/notifications?limit=&offset=              → PaginatedResponse<Notification>
GET    /api/notifications/unread?limit=&offset=       → PaginatedResponse<Notification>
PATCH  /api/notifications/{id}/read                    → 204
PATCH  /api/notifications/read-all                     → 204
DELETE /api/notifications/{id}                         → 204
```

### 4.10 Split Tags (etiquetas de rateio)
```
POST   /api/groups/{groupId}/split-tags               → { name: string } → SplitTag
GET    /api/groups/{groupId}/split-tags?limit=&offset=→ PaginatedResponse<SplitTag>
DELETE /api/split-tags/{id}                            → 204
GET    /api/split-tags/{id}/members?limit=&offset=    → PaginatedResponse<SplitTagMember>
POST   /api/split-tags/{id}/members                    → { user_id: string } → 204
DELETE /api/split-tags/{id}/members/{userId}           → 204
```

---

## 5. Interceptors e Guards

### `auth.interceptor.ts`
```typescript
@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const token = localStorage.getItem('access_token');
    if (token) {
      req = req.clone({ setHeaders: { Authorization: `Bearer ${token}` } });
    }
    return next.handle(req);
  }
}
```

### `refresh.interceptor.ts`
```typescript
@Injectable()
export class RefreshInterceptor implements HttpInterceptor {
  private isRefreshing = false;
  private refreshSubject = new Subject<string | null>();

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(req).pipe(
      catchError(err => {
        if (err.status === 401 && !req.url.includes('/auth/')) {
          return this.handleRefresh(req, next);
        }
        return throwError(() => err);
      })
    );
  }

  private handleRefresh(req: HttpRequest<any>, next: HttpHandler) {
    if (!this.isRefreshing) {
      this.isRefreshing = true;
      this.refreshSubject.next(null);

      return this.authService.refresh().pipe(
        switchMap(res => {
          this.isRefreshing = false;
          localStorage.setItem('access_token', res.token);
          this.refreshSubject.next(res.token);
          return next.handle(req.clone({
            setHeaders: { Authorization: `Bearer ${res.token}` }
          }));
        }),
        catchError(err => {
          this.isRefreshing = false;
          this.authService.logout();
          return throwError(() => err);
        })
      );
    }
    return this.refreshSubject.pipe(
      filter(token => token !== null),
      take(1),
      switchMap(token => next.handle(req.clone({
        setHeaders: { Authorization: `Bearer ${token}` }
      })))
    );
  }
}
```

---

## 6. PWA — Service Worker

### `ngsw-config.json`
```json
{
  "$schema": "./node_modules/@angular/service-worker/config/schema.json",
  "index": "/index.html",
  "assetGroups": [
    {
      "name": "app",
      "installMode": "prefetch",
      "resources": { "files": ["/favicon.ico", "/index.html", "/*.css", "/*.js"] }
    },
    {
      "name": "assets",
      "installMode": "lazy",
      "updateMode": "prefetch",
      "resources": { "files": ["/**/*.svg", "/assets/**"] }
    }
  ],
  "dataGroups": [
    {
      "name": "api",
      "urls": ["/api/**"],
      "cacheConfig": {
        "strategy": "freshness",
        "maxSize": 100,
        "maxAge": "5m"
      }
    }
  ]
}
```

### `manifest.webmanifest`
```json
{
  "name": "Doheem",
  "short_name": "Doheem",
  "theme_color": "#6C3FC5",
  "background_color": "#1E1E2E",
  "display": "standalone",
  "scope": "/",
  "start_url": "/",
  "icons": [
    { "src": "assets/icons/icon-192.png", "sizes": "192x192", "type": "image/png" },
    { "src": "assets/icons/icon-512.png", "sizes": "512x512", "type": "image/png" }
  ]
}
```

### Prompt de Instalação PWA
```typescript
// app.component.ts
@HostListener('window:beforeinstallprompt', ['$event'])
onBeforeInstallPrompt(event: Event) {
  event.preventDefault();
  this.deferredPrompt = event;
  this.showInstallBanner = true;
}

installPWA() {
  this.deferredPrompt?.prompt();
}
```

---

## 7. Design System (Tailwind)

### Paleta de Cores
```css
/* styles.css */
:root {
  --color-primary: #6C3FC5;
  --color-primary-light: #8B5CF6;
  --color-bg: #FFFFFF;
  --color-bg-card: #F4F4F8;
  --color-text: #1E1E2E;
  --color-text-secondary: #6B7280;
  --color-success: #10B981;
  --color-danger: #EF4444;
  --color-warning: #F59E0B;
}

.dark {
  --color-bg: #1E1E2E;
  --color-bg-card: #2A2A3E;
  --color-text: #F4F4F8;
  --color-text-secondary: #9CA3AF;
}
```

No `tailwind.config.js`:
```js
theme: {
  extend: {
    colors: {
      primary: { DEFAULT: '#6C3FC5', light: '#8B5CF6' },
      surface: { DEFAULT: '#F4F4F8', dark: '#2A2A3E' },
    }
  }
}
```

### Dark Mode Toggle
```typescript
// theme.service.ts
export class ThemeService {
  private readonly darkClass = 'dark';

  isDark$ = new BehaviorSubject<boolean>(
    localStorage.getItem('theme') === 'dark'
  );

  constructor() {
    this.isDark$.subscribe(dark => {
      document.documentElement.classList.toggle(this.darkClass, dark);
      localStorage.setItem('theme', dark ? 'dark' : 'light');
    });
  }

  toggle() { this.isDark$.next(!this.isDark$.value); }
}
```

### Bottom Navigation
5 ícones fixos na parte inferior (mobile-first):
1. **Dashboard** (house) → `/dashboard`
2. **Despesas** (wallet) → `/groups/{id}/expenses`
3. **Tarefas** (checklist) → `/groups/{id}/tasks`
4. **Cobranças** (cash) → `/groups/{id}/payments`
5. **Perfil** (person) → `/profile`

Em telas >= 768px: sidebar esquerda ao invés de bottom nav.

### FAB (Floating Action Button)
- Posição: bottom-right, 80px acima do bottom nav
- Ícone: `+`
- Expansão: mostrar 2 opções — "Nova Despesa" | "Nova Tarefa"
- `position: fixed; bottom: 6rem; right: 1.5rem;`

### Skeleton / Loading
Usar `skeleton.component.ts` para:
- Listas: 5 linhas com shimmer animation
- Cards: blocos retangulares pulsando
- Detalhes: placeholder de texto e imagem

### Empty State
Componente reutilizável com:
```html
<div class="flex flex-col items-center justify-center py-16 text-center">
  <img [src]="illustration" class="w-48 h-48 mb-4">
  <h3 class="text-lg font-semibold">{{ title }}</h3>
  <p class="text-sm text-secondary mt-1">{{ message }}</p>
</div>
```

---

## 8. Páginas — Especificação Detalhada

### 8.1 Onboarding (3 slides)
- Slide 1: "Divida contas com facilidade" — ilustração de pessoas dividindo dinheiro
- Slide 2: "Organize tarefas da casa" — ilustração de Kanban/checklist
- Slide 3: "Histórico completo" — ilustração de gráficos
- Botão "Pular" no canto superior direito
- Bolinhas de progresso (3 dots)
- Botão "Próximo" / "Começar" no último slide → `/login`

### 8.2 Login / Register
- **Login**: email + password + "Entrar" + link "Criar conta"
- **Register**: name + email + password + confirm_password + "Criar conta" + link "Já tenho conta"
- Validação: email format, password min 6, confirm match
- Após sucesso: salva token, redireciona para `/groups` (seleção de grupo) ou `/dashboard/{groupId}`

### 8.3 Dashboard (`/dashboard/:groupId`)
- Header: nome do grupo + botão configurações (admin apenas)
- **Card Resumo Financeiro**:
  - Total de gastos do mês: `R$ 4.250,00`
  - Meu saldo: "Você deve R$ 150,00" ou "Você tem R$ 230,00 a receber"
  - Pendências: "3 contas pendentes"
- **Seção Últimas Despesas** (5 itens): descrição, valor, quem pagou, data
- **Seção Tarefas Urgentes** (5 itens): título, responsável, dias restantes (vermelho se atrasado)
- **FAB** para nova despesa / tarefa

### 8.4 Despesas (`/groups/:groupId/expenses`)
- **Lista Cronológica**: card por despesa com avatar de quem pagou, descrição, valor, categoria (ícone), data, status (em aberto / pago)
- **Pull-to-refresh**: recarregar lista
- **Pagination**: scroll infinito ou "Carregar mais" + `PaginatedResponse.total`
- **Filtros**: por mês/ano, categoria, status
- **Formulário Rápido** (< 30s):
  1. Descrição (text)
  2. Valor (R$ — máscara de moeda)
  3. Data (date picker, default hoje)
  4. Categoria (dropdown: luz, água, aluguel, mercado, internet, outros)
  5. Divisão: "Igual entre todos" (default) | "Personalizada" (digitar valor por membro)
  6. Participantes: checkboxes dos moradores
  7. É parcelado? switch + número de parcelas
  8. Salvar
- **Detalhe da Despesa**: valor total, splits por pessoa com status (pago/pendente), botão "Marcar como pago" para o próprio split

### 8.5 Tarefas — Kanban (`/groups/:groupId/tasks`)
- 3 colunas com drag & drop (Angular CDK `CdkDropList`):
  - **Pendente** (status: `pending`) — laranja
  - **Em Andamento** (status: `in_progress`) — azul
  - **Concluída** (status: `completed`) — verde
- **Card de tarefa**: título, avatar do responsável, data de vencimento, badge de atraso (se > due_date)
- **Mover card**: `PATCH /api/tasks/occurrences/{id}/complete` ou `.../discard`
- **Nova tarefa**: FAB → modal com título, responsável (dropdown), categoria, data de vencimento, repetição (nunca/diário/semanal/mensal)

### 8.6 Cobranças (`/groups/:groupId/payments`)
- **Matriz "Quem deve para quem"**: tabela simplificada:
  ```
  | Devedor   | Credor     | Valor   | Status    | Ação     |
  |-----------|-----------|---------|-----------|----------|
  | Maria     | João      | R$ 50   | pendente  | [Cobrar] |
  ```
- **Botão "Cobrar"**: marca payment como `confirmed` + envia notificação via Kafka
- **Histórico de pagamentos**: lista cronológica de todos os pagamentos do grupo
- **Novo pagamento**: FAB → formulário: pagador, recebedor, valor, data, observação

### 8.7 Perfil (`/profile`)
- Foto + nome + email
- Editar perfil (name, email, avatar)
- Alterar senha
- **Seção Admin** (se role === "owner" || "admin"):
  - Gerenciar moradores: listar, adicionar por ID, remover, mudar papel (admin/member)
  - Editar nome/moeda do grupo
  - Desativar/ativar grupo
- Tema claro/escuro (toggle)
- Sair (logout)

### 8.8 Notificações (`/notifications`)
- Lista cronológica de notificações
- Badge de não lidas no bottom nav
- Swipe para marcar como lida
- Botão "Marcar todas como lidas"
- Ícone por tipo: `expense_added` (💰), `payment_confirmed` (✅), `task_completed` (📋)

---

## 9. Guards e Rotas

### `app.routes.ts`
```typescript
export const routes: Routes = [
  { path: '', redirectTo: '/onboarding', pathMatch: 'full' },
  { path: 'onboarding', component: OnboardingComponent },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  {
    path: '',
    canActivate: [AuthGuard],
    children: [
      { path: 'groups', component: GroupListComponent },
      { path: 'groups/new', component: GroupCreateComponent },
      { path: 'dashboard/:groupId', component: DashboardComponent },
      { path: 'groups/:groupId/expenses', component: ExpenseListComponent },
      { path: 'groups/:groupId/expenses/new', component: ExpenseFormComponent },
      { path: 'groups/:groupId/expenses/:id', component: ExpenseDetailComponent },
      { path: 'groups/:groupId/tasks', component: TaskKanbanComponent },
      { path: 'groups/:groupId/payments', component: PaymentBalanceComponent },
      { path: 'groups/:groupId/settings', component: GroupSettingsComponent },
      { path: 'invites', component: InviteListComponent },
      { path: 'notifications', component: NotificationListComponent },
      { path: 'profile', component: ProfileComponent },
    ]
  },
  { path: '**', redirectTo: '/onboarding' },
];
```

---

## 10. UX / Heurísticas de Nielsen (Obrigatório)

### Prevenção de Erros
- **Máscara de moeda**: enquanto digita, formatar como `R$ 1.234,56` (pipe `currencyMask`)
- **Máscara de data**: input type="date" ou máscara `__/__/____`
- **Confirmação**: modal antes de excluir despesa, remover morador, sair do grupo
- **Validação inline**: erro abaixo do campo enquanto digita (não só ao submeter)

### Visibilidade do Status do Sistema
- **Skeleton** enquanto carrega listas (não spinner circular)
- **Toast** de sucesso/erro após salvar (duração 3s, posição top)
- **Badge** de "salvando..." no botão de submit (disabled + spinner)
- **Pull-to-refresh** nas listas principais

### Reconhecimento vs Memorização
- **Avatar obrigatório** em cards de despesa, tarefa, membro — a foto do usuário, não apenas o nome
- **Avatar Group**: fotos sobrepostas dos participantes de uma despesa
- **Ícones por categoria**: 💡 luz, 💧 água, 🏠 aluguel, 🛒 mercado, 🌐 internet
- **Cores por status**: vermelho (atrasado), amarelo (pendente), verde (concluído/pago)

---

## 11. Responsividade (Mobile-First)

| Breakpoint | Largura | Layout |
|---|---|---|
| Mobile | < 480px | Bottom nav, cards full-width |
| Tablet | 480-768px | Bottom nav, grid 2 colunas |
| Desktop | > 768px | Sidebar + grid 3 colunas |

### Mobile-first CSS (Tailwind)
```html
<!-- Bottom nav visível só em mobile, sidebar em desktop -->
<nav class="fixed bottom-0 w-full md:hidden">...</nav>
<aside class="hidden md:flex md:flex-col md:w-64">...</aside>

<!-- Grid adaptativo -->
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
```

---

## 12. Dependências npm

```json
{
  "dependencies": {
    "@angular/common": "^19.0.0",
    "@angular/core": "^19.0.0",
    "@angular/router": "^19.0.0",
    "@angular/cdk": "^19.0.0",
    "@angular/service-worker": "^19.0.0",
    "rxjs": "~7.8.0"
  },
  "devDependencies": {
    "tailwindcss": "^4.0.0",
    "@tailwindcss/postcss": "^4.0.0",
    "typescript": "~5.7.0"
  }
}
```

---

## 13. Exemplo de App Config

### `app.config.ts`
```typescript
export const appConfig: ApplicationConfig = {
  providers: [
    provideHttpClient(
      withInterceptors([
        authInterceptor,
        refreshInterceptor,
      ]),
    ),
    provideServiceWorker('ngsw-worker.js', {
      enabled: isDevMode(),
      registrationStrategy: 'registerWhenStable:30000',
    }),
    provideRouter(routes),
  ],
};
```

---

## 14. Observações Importantes

1. **Refresh Token**: o cookie `refresh_token` é httpOnly — o JavaScript não consegue ler. O interceptor de refresh deve apenas chamar `POST /api/auth/refresh` e o cookie é enviado automaticamente pelo browser.
2. **Pagination**: todas as listas aceitam `?limit=XX&offset=XX` e retornam `{ data: [...], total: N }`. Use `total` para calcular número de páginas.
3. **Split Tags**: são etiquetas de rateio predefinidas (ex: "50% para cada" ou "apenas João e Maria"). Use no formulário de despesa como atalho de divisão.
4. **Drag & Drop Kanban**: ao soltar um card em outra coluna, chamar `PATCH /api/tasks/occurrences/{id}/complete` (se for "Concluída") ou `.../discard` (se "Descartada"). Tarefas "Pendente" e "Em Andamento" são gerenciadas pelo frontend (não têm endpoint específico — são todas `pending` no backend até serem completadas).
5. **File Upload**: enviar como `FormData` com campo `file`. Máx 10MB. Tipos: jpeg, png, gif, pdf.
