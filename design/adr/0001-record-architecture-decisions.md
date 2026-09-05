# ADR-0001: Record Architecture Decisions

- Status: Approved
- Date: 2026-09-05
- Author(s): Equipe de Desenvolvimento
- Deciders: Mantenedores do The Legend of the Go Dragon

## 1. Context
Projetos de software desenvolvidos colaborativamente por pessoas e múltiplos agentes de IA sofrem frequentemente com context loss e erosão arquitetural. Quando decisões prévias não estão registradas de forma imutável e estruturada, refatorações subsequentes correm o risco de reintroduzir problemas já resolvidos, quebrar contratos e violar diretrizes de design.

O projeto necessitava de um mecanismo padrão, leve e auditável para documentar escolhas técnicas fundamentais e suas restrições.

## 2. Decision
Adotamos o padrão de Architecture Decision Records (ADR) seguindo a especificação `adr-template` (da Danicat Skills Catalog). 

Todos os registros devem:
1. Ser gravados em formato Markdown em `design/adr/NNNN-short-descriptive-title.md`.
2. Seguir a estrutura fixa: Contexto, Decisão, Consequências e Conformidade/Verificação.
3. Ser imutáveis após aprovação. Se uma decisão mudar, cria-se um novo ADR com status superseded no original.
4. Manter tom objetivo e factual, documentando alternativas descartadas e trade-offs.

Alternativas consideradas:
- **Manter apenas documentação descritiva em `docs/architecture.md`**: Rejeitado porque documentos estáticos misturam o estado presente com o histórico do "porquê", ocultando as restrições históricas e debates técnicos.
- **Utilizar issues/PRs no GitHub**: Rejeitado porque a busca no histórico de issues tem alta fricção para agentes locais e desenvolvedores offline, além de não compor a árvore de arquivos versionada no Git.

## 3. Consequences
### Positivas
- Registro claro da evolução arquitetural do projeto.
- Alinhamento explícito entre engenheiros humanos e agentes de IA sobre as restrições do sistema.
- Facilidade para onboarding e auditorias técnicas.

### Negativas / Restrições
- Overhead de governança para manter e redigir ADRs sempre que uma decisão estrutural for tomada.
- Obrigatoriedade de manter o índice atualizado.

## 4. Compliance and Verification
- Agentes de IA são instruídos via `AGENTS.md` a consultar as decisões documentadas antes de introduzir mudanças estruturais.
- Revisões de PR (code reviews) devem exigir a criação de um ADR para quaisquer alterações fundamentais em frameworks, modelos de concorrência ou camadas de armazenamento.
