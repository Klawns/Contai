# Contas a Pagar e Contas a Receber

## Objetivo

Documentar o desenho funcional e tecnico dos modulos de Contas a Pagar e Contas a Receber do Contai.

Estes modulos devem representar compromissos financeiros futuros, separados das movimentacoes financeiras reais. Cadastrar, editar ou cancelar uma conta a pagar ou a receber nao altera saldo. O saldo so muda quando uma conta a pagar e marcada como paga ou uma conta a receber e marcada como recebida, pois nesse momento o backend cria uma transacao real e aplica o efeito no saldo da conta.

## Conceito principal

O sistema passa a ter dois tipos de registro financeiro:

- Compromisso financeiro: lancamento planejado, futuro ou pendente, usado para previsao, organizacao e cobranca operacional.
- Transacao financeira: movimentacao real ja ocorrida, usada para saldo, dashboard financeiro atual e relatorios financeiros consolidados.

Contas a pagar e contas a receber pertencem ao primeiro grupo. Receitas, despesas e transferencias atuais em `backend/internal/transactions` continuam sendo transacoes reais.

## Dominio recomendado

Crie um dominio compartilhado para compromissos, por exemplo `backend/internal/commitments`.

O compromisso deve ter um tipo persistido:

- `payable`: conta a pagar.
- `receivable`: conta a receber.

Status persistidos:

- `pending`: pendente.
- `paid`: pago, apenas para `payable`.
- `received`: recebido, apenas para `receivable`.
- `canceled`: cancelado.

`overdue` ou `vencido` nao deve ser persistido. Esse estado deve ser calculado quando o compromisso esta `pending` e a data de vencimento/previsao ja passou.

## Modelo de dados

Campos recomendados para persistencia:

- `id`: identificador do compromisso.
- `userID`: usuario dono do compromisso.
- `type`: `payable` ou `receivable`.
- `description`: descricao obrigatoria.
- `amount`: valor positivo em centavos.
- `dueDate`: vencimento da conta a pagar ou data prevista da conta a receber.
- `accountID`: conta usada para pagamento ou recebimento.
- `categoryID`: categoria financeira.
- `note`: observacao opcional.
- `status`: status persistido do ciclo de vida.
- `recurrence`: configuracao opcional de recorrencia (usando goroutines e channels).
- `settledAt`: data efetiva de pagamento ou recebimento.
- `settlementTransactionID`: transacao real gerada na quitacao.
- `canceledAt`: data de cancelamento, quando aplicavel.
- `createdAt`: data de criacao.
- `updatedAt`: data da ultima atualizacao.

Recorrencia pode entrar no desenho funcional inicial, mas a primeira implementacao pode limitar-se a compromissos unicos e deixar a geracao automatica recorrente para uma etapa posterior.

## Regras de negocio

- Cadastro de compromisso pendente nao altera saldo.
- Edicao de compromisso pendente nao altera saldo.
- Cancelamento de compromisso pendente nao altera saldo.
- Valor deve ser positivo.
- Descricao, valor, data, conta e categoria sao obrigatorios.
- Conta e categoria devem pertencer ao usuario autenticado.
- Conta e categoria precisam estar ativas.
- Conta a pagar deve usar categoria do tipo `expense`.
- Conta a receber deve usar categoria do tipo `income`.
- Conta a pagar gera transacao `expense` ao ser paga.
- Conta a receber gera transacao `income` ao ser recebida.
- Compromisso `paid`, `received` ou `canceled` nao pode ser editado pelo fluxo comum.
- Conta paga nao pode ser paga novamente.
- Conta recebida nao pode ser recebida novamente.
- Conta cancelada nao pode ser quitada.
- Compromisso quitado deve manter `settlementTransactionID`.
- Se a criacao da transacao falhar, o compromisso deve continuar no estado anterior.
- Se a atualizacao do compromisso falhar depois da criacao da transacao, a transacao e o saldo devem ser revertidos pela transacao de banco.
- Transacoes geradas por compromisso nao devem ser editadas ou removidas diretamente pelo fluxo comum de movimentacoes.
- Edicao ou exclusao de transacao gerada por compromisso deve exigir fluxo explicito de estorno, reabertura ou ajuste controlado.

## Casos de uso

Casos de uso comuns:

- Criar compromisso.
- Listar compromissos.
- Editar compromisso pendente.
- Cancelar compromisso pendente.
- Marcar conta a pagar como paga.
- Marcar conta a receber como recebida.

Filtros recomendados para listagem:

- Tipo: `payable` ou `receivable`.
- Status persistido: `pending`, `paid`, `received`, `canceled`.
- Status efetivo: incluir `overdue` como calculado.
- Periodo por `dueDate`.
- Conta.
- Categoria.

## Fluxo de pagamento ou recebimento

Ao marcar uma conta a pagar como paga ou uma conta a receber como recebida, execute tudo dentro de `UnitOfWork`.

Sequencia recomendada:

1. Abrir transacao de banco.
2. Buscar e bloquear o compromisso para evitar quitacao concorrente.
3. Validar que o compromisso pertence ao usuario autenticado.
4. Validar que o compromisso esta `pending`.
5. Validar que conta e categoria existem, pertencem ao usuario e estao ativas.
6. Validar tipo da categoria: `expense` para pagar, `income` para receber.
7. Criar a transacao real.
8. Aplicar o efeito da transacao no saldo da conta.
9. Atualizar o compromisso para `paid` ou `received`.
10. Gravar `settledAt` e `settlementTransactionID`.
11. Confirmar a transacao de banco.

Se qualquer etapa falhar, a unidade de trabalho deve reverter a operacao inteira.

## Integracao com transacoes

O fluxo atual de transacoes deve aceitar origem opcional:

- `manual`: transacao criada diretamente pelo usuario.
- `payable`: transacao gerada por conta a pagar.
- `receivable`: transacao gerada por conta a receber.

Campos recomendados em transacoes:

- `originType`: `manual`, `payable` ou `receivable`.
- `originID`: identificador do compromisso que gerou a transacao, quando houver.

Transacoes manuais continuam podendo ser criadas pelos endpoints atuais. Transacoes com origem `payable` ou `receivable` devem ser protegidas contra edicao ou remocao direta para evitar inconsistencia entre saldo, transacao e compromisso.

## Contratos HTTP

Endpoints de contas a pagar:

- `GET /api/payables`
- `POST /api/payables`
- `PATCH /api/payables/:id`
- `PATCH /api/payables/:id/pay`
- `PATCH /api/payables/:id/cancel`

Endpoints de contas a receber:

- `GET /api/receivables`
- `POST /api/receivables`
- `PATCH /api/receivables/:id`
- `PATCH /api/receivables/:id/receive`
- `PATCH /api/receivables/:id/cancel`

Os endpoints separados deixam o contrato mais claro para o frontend, mas internamente podem usar o mesmo dominio `commitments`.

### Criacao de conta a pagar

`POST /api/payables`

Payload:

```json
{
  "description": "Aluguel",
  "amount": 180000,
  "dueDate": "2026-06-10T12:00:00-03:00",
  "accountId": "account-id",
  "categoryId": "category-id",
  "note": "Pagamento mensal",
  "recurrence": null
}
```

Regras especificas:

- A categoria deve ser `expense`.
- A quitacao gera transacao `expense`.

### Criacao de conta a receber

`POST /api/receivables`

Payload:

```json
{
  "description": "Cliente ACME",
  "amount": 320000,
  "dueDate": "2026-06-15T12:00:00-03:00",
  "accountId": "account-id",
  "categoryId": "category-id",
  "note": "Parcela 1/3",
  "recurrence": null
}
```

Regras especificas:

- A categoria deve ser `income`.
- A quitacao gera transacao `income`.

### Pagamento de conta a pagar

`PATCH /api/payables/:id/pay`

Payload:

```json
{
  "amount": 180000,
  "paidAt": "2026-06-09T10:30:00-03:00",
  "accountId": "account-id",
  "categoryId": "category-id",
  "note": "Pago com desconto"
}
```

O payload permite ajustar valor, data efetiva, conta, categoria e observacao antes da criacao da transacao real.

### Recebimento de conta a receber

`PATCH /api/receivables/:id/receive`

Payload:

```json
{
  "amount": 320000,
  "receivedAt": "2026-06-15T14:00:00-03:00",
  "accountId": "account-id",
  "categoryId": "category-id",
  "note": "Recebido por PIX"
}
```

O payload permite ajustar valor, data efetiva, conta, categoria e observacao antes da criacao da transacao real.

## Frontend

Adicionar telas dentro de Planejamento:

- `Contas a pagar`.
- `Nova conta a pagar`.
- `Editar conta a pagar`.
- `Contas a receber`.
- `Nova conta a receber`.
- `Editar conta a receber`.

O fluxo atual de receitas e despesas continua representando lancamentos imediatos. Os novos modulos representam lancamentos planejados.

## Listas

As listas de contas a pagar e contas a receber devem conter:

- Seletor de mes ou periodo.
- Filtro por status.
- Filtro por conta e categoria, quando houver espaco adequado.
- Cards de resumo.
- Lista de compromissos.
- Estado vazio.
- Estado de carregamento.
- Estado de erro.

Resumo recomendado para contas a pagar:

- Total pendente.
- Total vencido.
- Total pago no periodo.
- Proximos vencimentos.

Resumo recomendado para contas a receber:

- Total previsto.
- Total vencido.
- Total recebido no periodo.
- Proximos recebimentos.

Acoes por item:

- Editar, quando `pending`.
- Cancelar, quando `pending`.
- Marcar como pago, para `payable` pendente.
- Marcar como recebido, para `receivable` pendente.
- Ver transacao gerada, quando quitado.

## Formularios

Formulario de Contas a Pagar:

- Valor.
- Descricao.
- Vencimento.
- Categoria de despesa.
- Conta de pagamento.
- Observacao.
- Recorrencia opcional.

Formulario de Contas a Receber:

- Valor.
- Descricao.
- Data prevista.
- Categoria de receita.
- Conta de recebimento.
- Observacao.
- Recorrencia opcional.

Validacoes de frontend devem ser superficiais e alinhadas ao backend. O frontend nao deve assumir regras de dominio que dependem de banco, permissao ou consistencia transacional.

## Confirmacao de quitacao

Marcar como pago ou recebido deve abrir uma confirmacao antes de criar a transacao.

A confirmacao deve permitir ajustar:

- Conta.
- Categoria.
- Valor.
- Data efetiva.
- Observacao.

Ao confirmar:

- O frontend chama `PATCH /api/payables/:id/pay` ou `PATCH /api/receivables/:id/receive`.
- O backend gera a transacao real.
- O compromisso muda para `paid` ou `received`.
- A movimentacao aparece no fluxo atual de Transacoes.
- Queries de compromissos, transacoes, contas, dashboard e resumos relacionados devem ser invalidadas.

## Dashboard e relatorios

Dashboard e relatorios financeiros atuais devem continuar calculando saldo real a partir de transacoes reais.

Novas secoes podem ser adicionadas sem misturar previsao com saldo real:

- A pagar.
- A receber.
- Vencidos.
- Proximos vencimentos.
- Previsao de entrada.
- Previsao de saida.

Relatorios podem ganhar secoes de compromissos pendentes e vencidos, mas os relatorios financeiros consolidados devem continuar baseados nas transacoes efetivadas.

## Backend: pontos de implementacao

- Criar `backend/internal/commitments` seguindo a organizacao dos dominios atuais.
- Criar entidade de persistencia e migracao/automigrate correspondente.
- Criar repositorio com busca por ID, listagem filtrada, criacao, atualizacao e bloqueio para quitacao.
- Criar servico de aplicacao com os casos de uso de criacao, listagem, edicao, cancelamento e quitacao.
- Reutilizar `UnitOfWork` existente para quitacao.
- Reutilizar validacoes de conta e categoria ja existentes ou criar portas especificas para consulta.
- Alterar transacoes para registrar origem opcional.
- Bloquear edicao/remocao direta de transacoes com origem diferente de `manual`.
- Registrar rotas autenticadas em `backend/internal/server/routes.go`.
- Atualizar injecao de dependencias em `backend/internal/server/dependencies.go`.

## Testes recomendados

Testes de dominio:

- Criacao valida de conta a pagar.
- Criacao valida de conta a receber.
- Valor invalido.
- Data obrigatoria.
- Categoria incompatibilizada pelo caso de uso.
- Status invalido para transicao.
- Cancelamento de pendente.
- Quitacao de pendente.
- Tentativa de quitacao duplicada.
- Tentativa de editar compromisso cancelado ou quitado.

Testes de servico:

- Criar compromisso sem alterar saldo.
- Editar compromisso sem alterar saldo.
- Cancelar compromisso sem alterar saldo.
- Pagar conta a pagar criando transacao `expense`.
- Receber conta a receber criando transacao `income`.
- Atualizar saldo da conta na quitacao.
- Fazer rollback quando a criacao da transacao falhar.
- Fazer rollback quando a atualizacao do compromisso falhar.
- Rejeitar conta ou categoria de outro usuario.
- Rejeitar conta ou categoria inativa.
- Rejeitar categoria de tipo incorreto.

Testes de integracao/HTTP:

- Validar contratos de `GET`, `POST`, `PATCH`, `pay`, `receive` e `cancel`.
- Validar autenticacao obrigatoria.
- Validar filtros por tipo, status, periodo, conta e categoria.
- Validar retorno de status efetivo `overdue` quando aplicavel.
- Validar que transacoes geradas por compromisso nao podem ser editadas ou excluidas diretamente.

Aceite manual de frontend:

- Criar conta a pagar.
- Editar conta a pagar pendente.
- Cancelar conta a pagar pendente.
- Marcar conta a pagar como paga.
- Criar conta a receber.
- Editar conta a receber pendente.
- Cancelar conta a receber pendente.
- Marcar conta a receber como recebida.
- Filtrar por status, periodo, conta e categoria.
- Confirmar que cadastro, edicao e cancelamento nao alteram saldo.
- Confirmar que pagar/receber altera saldo e gera transacao.
- Confirmar que dashboard nao mistura valores previstos com saldo real.

## Checklist de aceite

- A documentacao existe em `docs/contas-a-pagar-receber/README.md`.
- O desenho separa compromissos financeiros de transacoes reais.
- Cadastro, edicao e cancelamento de compromissos nao alteram saldo.
- Pagamento e recebimento geram transacao real dentro de `UnitOfWork`.
- Contas a pagar usam categorias `expense` e geram transacoes `expense`.
- Contas a receber usam categorias `income` e geram transacoes `income`.
- `overdue` e documentado como status calculado, nao persistido.
- Transacoes geradas por compromisso possuem origem e nao podem ser editadas ou removidas diretamente.
- Endpoints de payables e receivables estao documentados.
- Telas e fluxos principais do frontend estao documentados.
- Testes de dominio, servico, HTTP e aceite manual estao documentados.
