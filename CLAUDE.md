# Nokori

Personal iOS expense tracker (Flutter client, Go backend). Full context
lives in `docs/`:

- `docs/decisions.md` — the actual source of truth for *why* things are
  the way they are. Append-only log, chronological. Check this before
  asserting any technical or product choice is settled — it wins over
  any summary or memory of it.
- `docs/DDD.md` — current system shape and tech stack.
- `docs/PRD.md` — current product scope.

## Ticket conventions

- One issue template: `.github/ISSUE_TEMPLATE/ticket.yml`. No split by
  bug/feature or frontend/backend — those are just labels, not different
  question sets.
- Labels: type (`bug`, `feature`, `chore`, `docs`) and area (`backend`,
  `frontend`, `ai`, `infra`).
- **Keep tickets short.** State the goal and the done-when condition
  plainly. Don't restate context that's already in `docs/` — link to the
  relevant `decisions.md` entry instead of re-explaining the reasoning
  inline. A ticket that needs paragraphs is a sign the reasoning belongs
  in decisions.md, not in the ticket.
- **Follow-up**: record progress as comments on the issue as you go
  (what was tried, what broke) — that comment trail is the record of the
  trial-and-error process, not a separate thing to maintain. Close by
  referencing the ticket in the resolving commit/PR (`Fixes #12`), not by
  closing it manually after the fact. Group related tickets under a
  milestone by phase (e.g. "Backend skeleton," "Auth") if useful — no
  obligation to use milestones for everything.

## Steering documentation

This file (and any nested `CLAUDE.md` added under a subdirectory later,
e.g. `backend/CLAUDE.md`, `frontend/CLAUDE.md`) is the steering
mechanism for this project — Claude Code loads it automatically in every
session, on any machine. There's no separate steering system; this and
`docs/` are it.

**Rule: when something structurally new is added, update the relevant
doc in the same change, not as a follow-up.**
- A new product/architecture decision (a new feature, a tech swap, a
  scope change) → an entry in `docs/decisions.md`, with reasoning.
- A new standing convention (how tickets get filed, how a directory is
  organized, a rule for how something gets built) → this file, or a
  nested `CLAUDE.md` once `backend/` or `frontend/` exist and warrant
  their own scoped conventions.
- Don't create nested `CLAUDE.md` files pre-emptively before there's
  real code to document — an empty steering doc is worse than none,
  it's just another place to forget to fill in.
