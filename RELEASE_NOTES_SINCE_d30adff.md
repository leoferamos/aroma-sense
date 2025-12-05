# Release Notes - Desde Commit d30adff (feat/enhanced-user-deactivation)

## Visão Geral
Este release implementa o sistema completo de contestações de desativação de conta para compliance LGPD, incluindo backend, frontend admin e melhorias na UX. Também inclui correções e refatores menores.

## Novas Funcionalidades (Features)
- **Sistema de Contestações de Desativação de Conta (LGPD)**:
  - Novo modelo `UserContestation` para rastrear contestações.
  - Migração de banco: tabela `user_contestations` com campos para status, razão, timestamps e auditoria.
  - Serviço `UserContestationService` para criar, listar e aprovar/rejeitar contestações.
  - Handler `AdminContestationHandler` com endpoints para admins gerenciarem contestações.
  - Rotas admin: `/admin/contestations` (listar), `/admin/contestations/:id/approve`, `/admin/contestations/:id/reject`.
  - Integração com `LgpdService` para criar contestações via endpoint `/users/me/contest`.
  - Reativação automática: ao aprovar contestação, o usuário é reativado (remove suspensão) via `AdminUserService.AdminReactivateUser`.
  - Documentação Swagger completa para endpoints de contestações.

- **Página Admin de Contestações**:
  - Nova página `/admin/contestations` com tabela de contestações pendentes.
  - Funcionalidades: listar, aprovar/rejeitar com notas de revisão.
  - Integração com layout admin existente.

- **Melhorias no Frontend**:
  - Validação de mínimo 10 caracteres no formulário de contestação (botão desabilitado até cumprir).
  - Normalização de resposta da API para contestações (trata `data: null` como array vazio, evita ErrorBoundary).

## Correções (Bug Fixes)
- **AdminUserService**: Adicionada validação para `suspensionUntil` não poder ser no passado (commit d30adff).
- **Resposta de API**: Normalização para evitar crashes quando `data` é `null` em listagens de contestações.

## Melhorias Técnicas (Improvements)
- **Integração de Serviços**: `UserContestationService` agora aceita `UserRepository` e `AdminUserService` para resolver admin IDs e reativar usuários.
- **Auditoria**: Contestações aprovadas/rejeitadas registram `reviewed_by` com ID numérico do admin.
- **Documentação**: Swagger docs atualizados com endpoints de contestações e definições de segurança.

## Refatores (Refactors)
- **Notifier**: Removidos comentários desnecessários em `NewNotifier`.
- **Docs**: Removida menção a "predictable cost" no documento de arquitetura RAG do chat.

## Arquivos Alterados (Resumo)
- **Backend**: 18 arquivos modificados (principalmente adições para contestações: model, repo, service, handler, router, docs, migrations).
- **Frontend**: 4 arquivos modificados (nova página admin, validação no overlay, serviços).
- **Docs**: 3 arquivos modificados (Swagger, arquitetura).
- Total: 24 arquivos, +989 linhas, -84 linhas.

## Notas para Deploy
- Executar migração: `20251204170000_create_user_contestations_table.up.sql`.
- Verificar permissões admin para novos endpoints.
- Testar fluxo completo: desativar usuário → contestar → admin aprovar/rejeitar → reativação automática.

## Próximos Passos Sugeridos
- Implementar marcação de contestações pendentes como "superseded" em re-suspensões subsequentes.
- Adicionar notificações por email para resultados de contestações.
- Testes unitários para novos serviços e handlers.