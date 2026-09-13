[toc]

# Nokori (残り) — Design Direction Document

Status: Draft v0.1
Companion to: PRD.md

This document describes the shape of the system at a high level — the
major pieces, how they talk to each other, and the trajectory the
architecture is meant to grow along. It intentionally stops short of
schemas, API contracts, indexing strategy, concurrency implementation,
and infrastructure manifests — those are real engineering decisions
that belong to whoever is building the thing, made and recorded as the
build happens, not prescribed in advance from a planning document.

## 1. System shape

Three pieces:

- **Client** — the Flutter iOS app. Owns all user interaction: capture
  (text/photo/voice), the live budget meter, history, settings.
- **Backend API** — a Go service. Owns account/auth, budget cycle state,
  expense persistence, and (for users on the app's own AI access) the
  proxy to the AI provider.
- **Database** — Postgres. The system of record for accounts, budget
  cycles, and expenses.

A fourth party sits outside the system: the AI provider (Claude API),
reached either through the backend (hosted-key users) or directly from
the client (bring-your-own-key users).

## 2. Why a backend is required

Once there's an account system and per-user budget cycles that need to
persist and be trustworthy across devices and reinstalls, the data
can't live purely on-device. The backend's job in v1 is deliberately
narrow — auth, cycle bookkeeping, expense storage, AI proxying — not a
general-purpose platform. It should feel small and obviously correct at
MVP size, with room designed in to grow into the production-grade
system described in section 5.

## 3. Core flows (shape only)

### 3.1 Authentication

Sign in with Apple or email/password on the client → backend issues a
session/token → subsequent requests are authenticated with it. Apple
Sign In is not optional per platform policy, and is the primary path.

### 3.2 Expense capture

Text, photo, and voice inputs all converge on the client into one
question: "what expense record does this represent?" Photo and voice
need an extraction step (on-device OCR / on-device speech-to-text)
before that question can be answered by either simple parsing or an AI
call. In every case, the user confirms the result before it's written.
Once confirmed, the client sends a finished expense record to the
backend, which persists it and updates the cycle's running total.

### 3.3 Budget cycle

The backend is the source of truth for "what cycle are we in and how
much is left." It's computed from the user's salary date and budget
amount, and re-derived (not manually advanced) whenever the current
date crosses into a new cycle. The client reads this state to drive the
live meter; it should never compute the authoritative remaining balance
itself.

### 3.4 AI monthly analysis

On request, the backend (or client, for bring-your-own-key users)
gathers the current cycle's expense data and sends it to the AI
provider for a narrative summary. For hosted-key users, the backend is
responsible for enforcing the one-per-cycle limit before making the
call, since that's also where the cost is incurred.

## 4. Client/backend boundary — what stays on-device

To keep the backend's surface area small and to keep the app fast and
cheap to run, the following stay client-side rather than round-tripping
to the server:

- Receipt text extraction (on-device OCR).
- Speech-to-text transcription (on-device).
- AI calls for bring-your-own-key users (client talks to the provider
  directly; the backend never sees that traffic or that key).

## 5. Infrastructure trajectory

MVP ships on whatever gets a real, working, end-to-end app in front of
the developer fastest — a small managed host is enough to prove the
product. That is a starting point, not the destination. The intended
direction for this system, once the product itself is proven out, is a
properly designed production environment: containerized, deployed on
real cloud infrastructure, with the non-functional characteristics a
production system is expected to have — redundancy, observability,
a real CI/CD pipeline, and database design that holds up under
realistic load rather than MVP-scale data. The specific target
infrastructure and the path to it is a decision to be made deliberately
later, with its own design work, not defaulted into this document.

## 6. Open questions to resolve during build

These are flagged here so they aren't lost, but are intentionally left
unanswered — they're implementation decisions, not planning decisions:

- Exact schema for accounts / cycles / expenses, and indexing strategy.
- API contract between client and backend (REST shape, error format).
- Concurrency handling for the live meter when multiple expenses land
  in quick succession (e.g. photo processing finishing after a manual
  entry was already made).
- Where the monthly-analysis-count is tracked and how it's reset
  exactly on cycle rollover.
- Testing strategy for the backend (unit, integration, what's covered
  where).
