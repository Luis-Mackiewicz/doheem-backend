# Use Cases — Sistema de Gestão de República
**Versão 1.0**

---

## Atores

| Ator | Descrição |
|---|---|
| **Administrador** | Morador com perfil elevado. Gerencia o grupo, despesas, tarefas e configurações. |
| **Morador** | Membro ativo do grupo. Registra despesas, visualiza saldos e conclui tarefas. |
| **Sistema** | Executa ações automáticas: notificações, geração de parcelas, lembretes, exclusões. |

---

## 1. Gestão de Grupos

---

### UC-01 – Criar Grupo

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Usuário autenticado no sistema. |
| **Pós-condição** | Grupo criado e ativo, administrador definido como líder. |

**Fluxo principal:**
1. Administrador informa nome da república.
2. Sistema valida que o nome foi preenchido.
3. Sistema cria o grupo com o administrador como único membro ativo.
4. Grupo fica disponível para receber convites.

**Fluxos alternativos:**
- **1a.** Nome não informado → sistema bloqueia criação e exibe mensagem de campo obrigatório.

---

### UC-02 – Convidar Morador

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Grupo ativo com menos de 30 moradores. |
| **Pós-condição** | Link/QR Code de convite gerado e disponível para envio. |

**Fluxo principal:**
1. Administrador solicita geração de convite.
2. Sistema gera link/QR Code de uso único com validade de 48 horas.
3. Administrador compartilha o convite com o novo morador.
4. Novo morador acessa o link e entra no grupo.
5. Sistema invalida o link automaticamente após o primeiro acesso.

**Fluxos alternativos:**
- **2a.** Grupo já possui 30 moradores ativos → sistema bloqueia geração do convite e exibe mensagem de limite atingido.
- **3a.** Administrador revoga o convite antes do uso → link é invalidado imediatamente.

**Fluxos de exceção:**
- **4a.** Convite expirado (mais de 48h) → sistema informa expiração e orienta o administrador a gerar novo convite.

---

### UC-03 – Remover Morador

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Morador ativo no grupo. |
| **Pós-condição** | Morador removido; referenciado no histórico como "Morador removido [nome]". |

**Fluxo principal:**
1. Administrador seleciona o morador a ser removido.
2. Sistema verifica se o morador possui saldo zerado e sem tarefas pendentes.
3. Sistema solicita confirmação da ação.
4. Administrador confirma.
5. Morador é removido do grupo.
6. Sistema mantém o nome do morador no histórico como "Morador removido [nome]".

**Fluxos de exceção:**
- **2a.** Morador possui dívidas pendentes ou tarefas em aberto → sistema exibe resumo das pendências e bloqueia a remoção.

---

### UC-04 – Transferir Cargo de Administrador

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Grupo com ao menos 2 moradores ativos. |
| **Pós-condição** | Novo administrador definido; ator anterior passa a ter perfil de Morador. |

**Fluxo principal:**
1. Administrador seleciona um morador ativo para assumir o cargo.
2. Sistema solicita confirmação.
3. Administrador confirma.
4. Sistema transfere o cargo e rebaixa o ator anterior para perfil Morador.

---

### UC-05 – Sair do Grupo (Morador)

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador |
| **Pré-condição** | Morador ativo no grupo. |
| **Pós-condição** | Morador removido; histórico preservado com referência ao nome. |

**Fluxo principal:**
1. Morador solicita saída do grupo.
2. Sistema verifica saldo zerado e ausência de tarefas pendentes.
3. Sistema solicita confirmação.
4. Morador confirma.
5. Morador é removido do grupo.

**Fluxos de exceção:**
- **2a.** Morador possui dívidas ou tarefas pendentes → sistema exibe resumo das pendências e bloqueia a saída.

---

### UC-06 – Sair do Grupo (Administrador)

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Administrador ativo no grupo. |
| **Pós-condição** | Administrador removido ou grupo excluído permanentemente. |

**Fluxo principal (grupo com outros membros):**
1. Administrador solicita saída.
2. Sistema exige transferência do cargo antes de prosseguir.
3. Administrador executa UC-04.
4. Administrador realiza UC-05 como Morador.

**Fluxo alternativo (último membro do grupo):**
1. Administrador solicita saída.
2. Sistema exibe aviso explícito de que o grupo e todos os dados serão excluídos permanentemente.
3. Sistema exige confirmação digitada pelo administrador.
4. Administrador digita a confirmação.
5. Sistema exclui o grupo e todos os dados vinculados (despesas, tarefas, histórico).

---

### UC-07 – Editar Informações do Grupo

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Grupo ativo. |
| **Pós-condição** | Informações do grupo atualizadas. |

**Fluxo principal:**
1. Administrador acessa as configurações do grupo.
2. Administrador edita os campos desejados (nome, moeda, configurações de notificação).
3. Sistema valida os dados informados.
4. Sistema salva as alterações.

---

## 2. Divisão de Despesas

---

### UC-08 – Registrar Despesa

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador / Administrador |
| **Pré-condição** | Grupo ativo com ao menos 1 morador. |
| **Pós-condição** | Despesa registrada; saldos dos moradores envolvidos atualizados; notificações enviadas. |

**Fluxo principal:**
1. Ator preenche: descrição, valor total, data de competência, responsável pelo pagamento, categoria e modo de rateio.
2. Ator seleciona o modo de rateio (UC-09).
3. Sistema valida que a soma do rateio é igual ao valor total.
4. Sistema registra a despesa.
5. Sistema atualiza os saldos dos moradores envolvidos em tempo real.
6. Sistema notifica todos os moradores incluídos no rateio com o valor da sua parte.

**Fluxos alternativos:**
- **1a.** Ator informa data de vencimento → sistema agenda notificação para X dias antes (padrão: 3 dias).
- **1b.** Ator registra despesa como parcelada → sistema executa UC-10.
- **3a.** Soma do rateio difere do valor total → sistema bloqueia o registro e exibe a diferença.
- **3b.** Há diferença de centavos por arredondamento → sistema atribui a diferença automaticamente ao responsável pelo pagamento.

---

### UC-09 – Selecionar Modo de Rateio

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador / Administrador |
| **Pré-condição** | Registro de despesa em andamento (UC-08). |
| **Pós-condição** | Modo de rateio definido e partes calculadas. |

**Fluxo principal — Dividir entre todos:**
1. Ator seleciona "Dividir entre todos".
2. Sistema divide o valor igualmente entre todos os moradores ativos no momento do registro.

**Fluxo alternativo A — Dividir entre alguns:**
1. Ator seleciona "Dividir entre alguns membros".
2. Ator seleciona ao menos 2 moradores.
3. Sistema divide o valor igualmente entre os selecionados.

**Fluxo alternativo B — Valores personalizados:**
1. Ator seleciona "Valores personalizados".
2. Ator define valor ou percentual para cada morador.
3. Sistema valida que a soma é exatamente igual ao total.

**Fluxo alternativo C — Etiqueta de rateio:**
1. Ator seleciona uma etiqueta pré-configurada pelo administrador.
2. Sistema aplica os moradores vinculados à etiqueta e divide igualmente entre eles.

---

### UC-10 – Registrar Despesa Parcelada

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador / Administrador |
| **Pré-condição** | Registro de despesa em andamento (UC-08). |
| **Pós-condição** | Parcelas geradas automaticamente com vencimentos mensais. |

**Fluxo principal:**
1. Ator indica que a despesa é parcelada.
2. Ator informa número de parcelas e data da primeira.
3. Sistema gera as parcelas subsequentes com vencimento mensal.
4. Quando a data não existir no mês seguinte (ex.: dia 31), o sistema ajusta para o último dia do mês.
5. Cada parcela é tratada como cobrança independente no saldo dos moradores.

---

### UC-11 – Editar Despesa

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Despesa registrada no sistema. |
| **Pós-condição** | Despesa atualizada; registro de auditoria gerado; moradores notificados. |

**Fluxo principal:**
1. Administrador localiza a despesa e solicita edição.
2. Administrador altera os campos desejados.
3. Sistema valida os dados.
4. Sistema salva a alteração e registra em auditoria: campo alterado, valor anterior, valor novo, usuário e timestamp.
5. Sistema notifica todos os moradores afetados pela despesa.

**Fluxos de exceção:**
- **1a.** Despesa não pode ser excluída — apenas editada. Botão de exclusão não está disponível.

---

### UC-12 – Criar Etiqueta de Rateio

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Grupo ativo com ao menos 2 moradores. |
| **Pós-condição** | Etiqueta salva e disponível para uso no registro de despesas. |

**Fluxo principal:**
1. Administrador acessa o gerenciamento de etiquetas.
2. Administrador define um nome para a etiqueta e seleciona os moradores vinculados.
3. Sistema valida que ao menos 2 moradores foram selecionados.
4. Sistema salva a etiqueta.

---

### UC-13 – Registrar Pagamento de Dívida

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador / Administrador |
| **Pré-condição** | Morador possui saldo negativo (dívida) com outro morador. |
| **Pós-condição** | Pagamento registrado com status "Aguardando confirmação". |

**Fluxo principal:**
1. Ator informa: quem pagou, para quem, valor e data.
2. Ator pode anexar até 3 comprovantes (JPG, JPEG, PNG ou PDF; máximo 5MB cada).
3. Sistema registra o pagamento com status **"Aguardando confirmação"**.
4. Sistema notifica o morador recebedor para confirmar o pagamento.

**Fluxos alternativos:**
- **1a.** Pagamento parcial → sistema registra o valor informado e atualiza o saldo proporcionalmente após confirmação.
- **2a.** Arquivo fora do formato ou acima de 5MB → sistema rejeita o arquivo e exibe mensagem de erro explícita.

---

### UC-14 – Confirmar Recebimento de Pagamento

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador (recebedor) / Administrador |
| **Pré-condição** | Pagamento com status "Aguardando confirmação". |
| **Pós-condição** | Saldos de ambas as partes atualizados definitivamente. |

**Fluxo principal:**
1. Recebedor recebe notificação de pagamento pendente.
2. Recebedor acessa o pagamento e confirma o recebimento.
3. Sistema atualiza os saldos definitivamente.
4. Sistema registra a confirmação em auditoria com timestamp.

**Fluxos alternativos:**
- **1a.** Prazo de 72h expirado sem confirmação → sistema notifica pagador e recebedor.
- **1b.** Pagador cancela o pagamento antes da confirmação → pagamento é removido e saldos não são alterados.

**Fluxos de exceção:**
- Após confirmação, o pagamento não pode ser desfeito — apenas consultado no histórico.

---

### UC-15 – Consultar Histórico de Despesas e Pagamentos

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador / Administrador |
| **Pré-condição** | Membro ativo do grupo. |
| **Pós-condição** | Histórico exibido conforme filtros aplicados. |

**Fluxo principal:**
1. Ator acessa o histórico do grupo.
2. Ator aplica filtros opcionais: período, categoria, morador, status (pago/pendente).
3. Sistema exibe os registros correspondentes.

---

### UC-16 – Visualizar Saldo Individual

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador / Administrador |
| **Pré-condição** | Membro ativo do grupo. |
| **Pós-condição** | Saldo atualizado exibido ao ator. |

**Fluxo principal:**
1. Ator acessa a tela de saldos.
2. Sistema exibe o saldo consolidado do ator no grupo.
3. Sistema exibe a sugestão de saldo simplificado (menor conjunto de transações para zerar todas as dívidas do grupo).

---

## 3. Gestão de Tarefas

---

### UC-17 – Criar Tarefa

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Grupo ativo com ao menos 1 morador. |
| **Pós-condição** | Tarefa criada e atribuída ao morador responsável. |

**Fluxo principal:**
1. Administrador informa: descrição, responsável, data limite e categoria.
2. Sistema valida os campos obrigatórios.
3. Sistema cria a tarefa e notifica o morador responsável.

**Fluxos alternativos:**
- **1a.** Tarefa configurada como recorrente → administrador define periodicidade (diária, semanal, mensal). Sistema passa a gerenciar a geração automática das próximas ocorrências.

---

### UC-18 – Editar ou Reatribuir Tarefa

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Tarefa existente no grupo. |
| **Pós-condição** | Tarefa atualizada; registro de auditoria gerado. |

**Fluxo principal:**
1. Administrador seleciona a tarefa e edita os campos (descrição, data limite, categoria, responsável).
2. Sistema valida os dados.
3. Sistema salva a alteração e registra em auditoria com usuário e timestamp.

---

### UC-19 – Concluir Tarefa

| Campo | Descrição |
|---|---|
| **Ator principal** | Morador responsável / Administrador |
| **Pré-condição** | Tarefa atribuída e em aberto. |
| **Pós-condição** | Tarefa marcada como concluída; próxima ocorrência gerada (se recorrente). |

**Fluxo principal:**
1. Ator acessa a tarefa e marca como concluída.
2. Sistema registra quem realizou a conclusão e o timestamp.
3. Se a tarefa for recorrente, sistema gera automaticamente a próxima ocorrência com nova data limite.

**Fluxos de exceção:**
- **1a.** Morador que não é o responsável tenta concluir → sistema bloqueia a ação (somente o responsável ou o administrador podem concluir).

---

### UC-20 – Descartar Ocorrência de Tarefa Recorrente

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Tarefa recorrente com ocorrência em aberto ou atrasada. |
| **Pós-condição** | Ocorrência atual descartada; próxima ocorrência gerada. |

**Fluxo principal:**
1. Administrador seleciona a ocorrência atrasada ou em aberto.
2. Administrador solicita descarte da ocorrência.
3. Sistema descarta a ocorrência atual e gera a próxima automaticamente.

---

### UC-21 – Encerrar Recorrência de Tarefa

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Tarefa configurada como recorrente. |
| **Pós-condição** | Recorrência encerrada; nenhuma nova ocorrência é gerada. |

**Fluxo principal:**
1. Administrador acessa a tarefa recorrente.
2. Administrador solicita encerramento da recorrência.
3. Sistema confirma a ação e interrompe a geração de novas ocorrências.

---

## 4. Notificações

---

### UC-22 – Notificar Moradores sobre Nova Despesa

| Campo | Descrição |
|---|---|
| **Ator principal** | Sistema |
| **Gatilho** | Registro de nova despesa (UC-08). |
| **Pós-condição** | Todos os moradores do rateio notificados com o valor da sua parte. |

**Fluxo principal:**
1. Nova despesa é registrada.
2. Sistema identifica os moradores incluídos no rateio.
3. Sistema envia notificação individual para cada morador com o valor da sua parte.

---

### UC-23 – Enviar Lembrete de Dívida Pendente

| Campo | Descrição |
|---|---|
| **Ator principal** | Sistema |
| **Gatilho** | Dívida pendente há mais de X dias (configurável; padrão: 3 dias). |
| **Pós-condição** | Morador notificado sobre a dívida pendente. |

**Fluxo principal:**
1. Sistema verifica diariamente moradores com dívidas pendentes além do prazo configurado.
2. Sistema envia lembrete ao morador devedor.
3. Sistema repete o envio respeitando o intervalo mínimo de 3 dias entre lembretes.
4. Sistema interrompe os lembretes após 5 envios ou quando a dívida for quitada.

---

### UC-24 – Notificar Vencimento de Despesa

| Campo | Descrição |
|---|---|
| **Ator principal** | Sistema |
| **Gatilho** | Data de vencimento de despesa se aproximando (X dias antes; padrão: 3). |
| **Pós-condição** | Administrador e moradores envolvidos notificados. |

**Fluxo principal:**
1. Sistema verifica diariamente despesas com vencimento próximo.
2. Sistema envia notificação ao administrador e aos moradores envolvidos no rateio.

---

### UC-25 – Notificar Prazo de Tarefa

| Campo | Descrição |
|---|---|
| **Ator principal** | Sistema |
| **Gatilho** | 24 horas antes do prazo limite de uma tarefa. |
| **Pós-condição** | Morador responsável notificado. |

**Fluxo principal:**
1. Sistema verifica tarefas com prazo nas próximas 24 horas.
2. Sistema notifica o morador responsável pela tarefa.

---

### UC-26 – Alertar Tarefa Vencida

| Campo | Descrição |
|---|---|
| **Ator principal** | Sistema |
| **Gatilho** | Tarefa não concluída após a data limite. |
| **Pós-condição** | Tarefa sinalizada como atrasada; administrador notificado. |

**Fluxo principal:**
1. Sistema identifica tarefas com data limite ultrapassada e não concluídas.
2. Sistema marca a tarefa como **atrasada**.
3. Sistema notifica o administrador imediatamente.

---

## 5. Regras Gerais

---

### UC-27 – Consultar Auditoria

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Pré-condição** | Grupo ativo com histórico de ações. |
| **Pós-condição** | Registros de auditoria exibidos. |

**Fluxo principal:**
1. Administrador acessa o log de auditoria.
2. Sistema exibe registros de ações relevantes: criação/edição de despesas, conclusão de tarefas, pagamentos — com usuário e timestamp.

---

### UC-28 – Reativar Grupo Inativo

| Campo | Descrição |
|---|---|
| **Ator principal** | Administrador |
| **Gatilho** | Notificação de inatividade (30 dias antes da exclusão automática). |
| **Pós-condição** | Grupo reativado; contador de inatividade reiniciado. |

**Fluxo principal:**
1. Sistema envia notificação ao administrador 30 dias antes da exclusão por inatividade (12 meses sem atividade registrada).
2. Administrador realiza qualquer ação registrada no grupo (despesa, tarefa ou pagamento).
3. Sistema reativa o grupo e reinicia o contador de inatividade.

**Fluxos de exceção:**
- **2a.** Administrador não realiza nenhuma ação no prazo → sistema exclui o grupo e todos os dados permanentemente, em conformidade com a LGPD.