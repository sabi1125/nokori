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

- `open-api/` — the API contract, written before the backend implements
  it, no dependency on `backend/` (see docs/decisions.md — "API
  documentation, spec-first" and "OpenAPI spec: split by epic").
  - `spec/` — **the files you edit.** `openapi.yaml` is the table of
    contents: every URL listed, each `$ref`-ing a short descriptive key
    in `paths/<epic>.yaml` (`paths/system.yaml` for non-epic endpoints
    like `/health`). Reused pieces live in `components/schemas.yaml`
    (data shapes, incl. the shared `Error`) and
    `components/responses.yaml` (whole error responses). `$ref`s resolve
    relative to the file they're in — from `paths/`, components are
    `../components/...`.
  - `public/` — what nginx serves: `index.html` (Scalar via CDN) and
    `openapi.yaml`, which is **generated** from `spec/` (gitignored,
    never edit it). Scalar can't follow `$ref`s across files, hence the
    bundle.
  - `Dockerfile` — for deploys: `bundler` stage (Redocly CLI + watcher
    script, no spec; also what `open-api-bundler` runs in dev), `bundle`
    stage (bundles `spec/` — a broken spec fails the build on purpose),
    then nginx. In dev, `open-api` is plain `nginx:alpine` serving the
    bind mount, so a broken spec never blocks startup.
- `backend/` — **exists.** Go backend (CodeSeed scaffold), its own
  Dockerfile, local port **8080**. Migrations in `backend/migrations/`.
  Built *to* whatever the spec specifies.
- `frontend/` — **exists** (Flutter skeleton only). Built and run via
  Xcode, no Dockerfile.
- `docs/` — PRD, DDD, decisions.md, and diagrams (ER, state, etc.) — can
  be split into subfolders for cleanliness as it grows.
- Root `docker-compose.yml` — `open-api` (port **3000**),
  `open-api-bundler` (dev watcher), `mysql` (3306), `backend` (8080).
  Add services here rather than a separate compose file per service.

**Running `open-api` for local dev**: `docker compose up open-api` from
the repo root (add `--build` after changing `open-api/Dockerfile` or the
bundler script). This also starts `open-api-bundler` (nginx waits for
its first bundle),
which rebuilds `public/openapi.yaml` within ~1s of any save under
`spec/` — refresh the browser to see it. If the page stops updating,
check `docker compose logs open-api-bundler` for `bundle FAILED`. Note
the bundler only catches YAML/`$ref` breakage; a typo'd OpenAPI keyword
(e.g. `request:` instead of `requestBody:`) is silently ignored by
Scalar — `npx @redocly/cli lint open-api/spec/openapi.yaml` catches
those.

`open-api/public/` is bind-mounted as a *directory*, not individual
files. Keep it that way: a single-*file* bind mount breaks the moment
the file is replaced via atomic rename (most editors, `sed -i`, and the
bundler itself) — the mount stays pinned to the old inode and silently
keeps serving stale content.

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
- **Every title starts with its kind: `[Epic]`, `[Docs]` or `[Impl]`.**
  The type label backs it up: `[Docs]` tickets carry `docs`; `[Impl]`
  tickets carry `feature`, `chore` or `bug`.
- **Design comes before implementation, and the tree shows it:**
  `[Epic]` → `[Docs]` design ticket → the `[Impl]` tickets that build on
  that design. Backend work sits under a `[Docs] API spec: …` ticket
  (done when the contract is in `openapi.yaml`); frontend work sits under
  a `[Docs] UI design: …` ticket (done when the sketch is in `docs/ui/`).
  Every `[Impl]` ticket has a `[Docs]` parent. To move an existing child
  to a new parent, add `-F replace_parent=true` to the sub_issues call.
- Work outside any PRD feature (setup, tooling, docs upkeep) goes under
  `[Epic] Project setup & housekeeping` (#84).
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

## Code reviews

This is a personal project. Speed matters more than polish.

### Blocking: flag these, and I fix them before moving on
- Bugs that crash the app or break a core feature
- Data loss or corruption
- Security holes (exposed secrets, injection, auth bypass)
- Anything that would be painful to undo later (data formats, schema, public APIs)

### Non-blocking: don't hold me up over these
- Naming, style, formatting
- "Cleaner" or "more idiomatic" ways to write working code
- Missing tests, unless the code is risky
- Small performance wins
- Edge cases that are unlikely in my real usage

### Format
- Blocking issues go first. For each one: where it is, what breaks, and a hint on where to look.
- Put non-blocking notes in a short "Later" list of 5 items at most, one line each.
- If nothing is blocking, say "Nothing blocking" plainly so I can keep going.
- Don't repeat a "Later" item I've already chosen to skip.

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
