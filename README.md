# Nokori (残り)

Somewhere between one salary date and the next, you stop trusting your
own budget — the app resets on the 1st, your paycheck never agreed to
that. So you guess, and go back to not really knowing.

Nokori — Japanese for "the remainder" — starts its clock on your salary
date, not the calendar. Savings come first: a real cut, carved out
before you ever see a spendable number, so what's left is what's
actually yours to spend. It updates live as you log every expense,
honestly, even when it goes negative — and once a cycle closes, Claude
reads back what actually happened: not just a total, a read on your
patterns, and how this cycle stacked up against the last.

## Stack

- **Backend**: Go, [Echo](https://echo.labstack.com/), GORM, MySQL,
  go-migrate for schema migrations.
- **Frontend**: Flutter, iOS only.
- **API docs**: hand-authored OpenAPI, rendered via
  [Scalar](https://scalar.com/), designed before the backend implements
  it.

Full detail: [`docs/DDD.md`](docs/DDD.md).

## Repo layout

- `open-api/` — hand-authored `openapi.yaml` + a Scalar-based viewer,
  served by its own Dockerfile. No dependency on `backend/`.
- `backend/` — Go service. Owns accounts, budget cycles, expense
  persistence, and AI-provider proxying for hosted-key users.
- `frontend/` — Flutter iOS client.
- `docs/` — [`PRD.md`](docs/PRD.md) (product scope),
  [`DDD.md`](docs/DDD.md) (system shape), and
  [`decisions.md`](docs/decisions.md) (the actual source of truth for
  *why* — check this before assuming anything described elsewhere is
  settled; it wins over any summary of it).

## Local dev

From the repo root:

```
docker compose up
```

Brings up:
- **open-api docs** — http://localhost:3000
- **MySQL** — localhost:3306

**Backend**:
```
cd backend
cp .env.example .env   # fill in DB credentials matching your local MySQL
go run ./cmd/backend
```
Currently exposes `GET /health` only — real endpoints land as the
backend tickets in GitHub Issues close.

**Frontend**:
```
cd frontend
flutter pub get
flutter run
```
Local builds are installed via Xcode during active development (no CI
yet); free Apple-ID signing expires every 7 days.

## Contributing

Solo project, but if you're picking this up on another machine: see
[`CLAUDE.md`](CLAUDE.md) for ticket conventions, steering docs, and how
this repo expects work to be tracked.
