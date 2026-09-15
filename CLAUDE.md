# Nokori

Personal iOS expense tracker (Flutter client, Go backend). Full context
lives in `docs/`:

- `docs/decisions.md` — the actual source of truth for *why* things are
  the way they are. Append-only log, chronological. Check this before
  asserting any technical or product choice is settled — it wins over
  any summary or memory of it.
- `docs/DDD.md` — current system shape and tech stack.
- `docs/PRD.md` — current product scope.

## Repo layout

Documented here so any session (any machine) knows where things land,
rather than guessing:

- `open-api/` — **exists.** `open-api/public/` holds the served files:
  hand-authored `openapi.yaml` and `index.html` (Scalar via CDN script).
  `open-api/Dockerfile` just copies `public/` into nginx. Written before
  backend implementation, no dependency on `backend/` existing (see
  docs/decisions.md — "API documentation, spec-first"). Fill in
  `paths`/`components.schemas` in `openapi.yaml` as the API contract
  gets designed — that design work happens in this file, not in code.
  Has one real path so far: `GET /health`.
- `backend/` — not created yet. Go backend, its own Dockerfile, local
  port **8080** (already assumed by `open-api/public/openapi.yaml`'s
  `servers` block). Built *to* whatever the spec specifies once that's
  designed.
- `frontend/` — not created yet. Flutter app, its own Dockerfile.
- `docs/` — PRD, DDD, decisions.md, and diagrams (ER, state, etc.) — can
  be split into subfolders for cleanliness as it grows.
- Root `docker-compose.yml` — **exists**, with the `open-api` service
  (port **3000**). Add `backend`/`frontend` services here as those roots
  get built, rather than a separate compose file per service.

**Running `open-api` for local dev**: `docker compose up --build` from
the repo root (drop `--build` after the first run, or whenever
`open-api/Dockerfile` itself hasn't changed). `open-api/public/` is
bind-mounted as a *directory*, not individual files — editing
`openapi.yaml` or `index.html` shows up on a browser refresh with no
rebuild or restart needed. This only works because it's a directory
mount: a single-*file* bind mount breaks the moment the file is edited
by anything that saves via atomic rename (which is most editors,
including VS Code's default, and `sed -i`) — the mount stays pinned to
the old inode and silently keeps serving stale content. Don't go back to
per-file mounts here even though it looks more precise; it isn't.

Build each root when its own work actually starts, not pre-emptively —
same reasoning as not pre-creating empty nested `CLAUDE.md` files below.

## Ticket conventions

- One issue template: `.github/ISSUE_TEMPLATE/ticket.yml`. No split by
  bug/feature or frontend/backend — those are just labels, not different
  question sets.
- Labels: type (`bug`, `feature`, `chore`, `docs`) and area (`backend`,
  `frontend`, `ai`, `infra`, `open-api`).
- **Every ticket gets at least one type label and one area label, and is
  assigned to `sabi1125` (Sabir) — no unlabeled or unassigned tickets.**
- **Feature-sized work gets one parent ticket per PRD feature (labeled
  `epic`, no type/area label needed on the parent itself), with the
  actual backend/frontend implementation tickets attached as GitHub
  native sub-issues** (`gh api repos/{owner}/{repo}/issues/{parent}/sub_issues -X POST -F sub_issue_id={child's numeric id}`
  — note: `sub_issue_id` is the child issue's internal `id` field, not
  its issue number). This is also the answer to "how do I find just the
  parent tickets": GitHub has no built-in search qualifier for
  sub-issues (`has:sub-issues` looks plausible but silently does
  nothing) — filter with `label:epic` instead.
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
