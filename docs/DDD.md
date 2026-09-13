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

Four pieces — the AI provider is not an afterthought bolted onto the
side; the monthly analysis and cycle comparison are both core MVP
features, so Claude is a first-class part of this system's shape, not an
external nice-to-have:

- **Client** — the Flutter iOS app. Owns all user interaction: capture
  (text/photo/voice), the live budget meter, history, settings.
- **Backend API** — a Go service. Owns account/auth, budget cycle state,
  expense persistence, and (for users on the app's own AI access) the
  proxy to the AI provider.
- **Database** — Postgres. The system of record for accounts, budget
  cycles, and expenses.
- **AI provider (Claude API)** — generates the monthly spending analysis
  and the cycle comparison (PRD 5.6, 5.7) — two distinct features, each
  with its own hosted-key cap. Reached one of two ways depending on the
  user: through the backend, which proxies the call and enforces the
  relevant cap (hosted-key users), or directly from the client, which
  talks to the provider on its own (bring-your-own-key users). It's the
  one piece of the system Nokori doesn't operate itself, but it's still
  a dependency the design has to account for deliberately — prompt/response shape,
  failure handling, latency — not just "call an API and hope."

## 2. Why a backend is required

Once there's an account system and per-user budget cycles that need to
persist and be trustworthy across devices and reinstalls, the data
can't live purely on-device. The backend's job in v1 is to fully and
properly own everything that needs a trusted, persistent source of
truth — auth, cycle bookkeeping, expense persistence, AI proxying for
hosted-key users. This is an MVP, not a proof of concept: the backend
should do that job completely and correctly from the start, not be
trimmed down just to look small. It's still scoped to Nokori's own
domain rather than growing into a general-purpose platform — but within
that domain, "MVP" describes the feature set, not the amount of care
put into building it. See section 5 for how that same standard applies
to infrastructure.

## 3. Core flows (shape only)

### 3.1 Authentication

Email/password on the client → backend issues a session/token →
subsequent requests are authenticated with it. This is the only auth
method — no OAuth, from Apple or anyone else (see decisions.md: since
there's no third-party login of any kind, Apple's Sign in with Apple
requirement never actually applies).

### 3.2 Expense capture

Text, photo, and voice inputs all converge on the client into one
question: "what expense record does this represent?" Photo and voice
need an extraction step (on-device OCR / on-device speech-to-text)
before that question can be answered by either simple parsing or an AI
call. In every case, the user confirms the result before it's written.
Once confirmed, the client sends a finished expense record to the
backend, which persists it. Whether that expense affects the cycle's
tracked (discretionary) spend depends on its category (PRD 5.3): a
discretionary-category expense counts toward it, a basic-needs-category
expense doesn't, since that amount is already accounted for in the
budget calculation (PRD 5.2) and counting it again would double it.

### 3.3 Budget cycle

The backend is the source of truth for "what cycle are we in and how
much is left." Cycle timing is computed from the user's salary date; the
cycle's budget is either a suggested figure derived deterministically
from salary and basic needs, or that figure as edited by the user (see
PRD 5.2 for the formula — no AI involvement here). Both are re-derived
(not manually advanced) whenever the current date crosses into a new
cycle. Unspent budget never carries across that boundary — each new
cycle's tracked spend starts at zero against its own budget figure, not
the previous cycle's leftover. The client reads this state to drive the
live meter; it should never compute the authoritative remaining balance
itself.

### 3.4 AI monthly analysis

On request, the backend (or client, for bring-your-own-key users)
gathers the current cycle's expense data and sends it to the AI
provider for a narrative summary. For hosted-key users, the backend is
responsible for enforcing the one-per-cycle limit before making the
call, since that's also where the cost is incurred.

### 3.5 AI cycle comparison

A separate flow from 3.4, with its own hosted-key cap (PRD 5.7) — this
one can be triggered mid-cycle, not just once per cycle. On request, the
backend gathers spend-to-date for the current cycle and spend *up to
the same elapsed point* in the previous cycle (not the previous cycle's
final total — a mid-cycle partial total needs a like-for-like
comparison, not a full-cycle one), broken down by category, and sends
both to the AI provider for a comparison. The response returns both the
narrative and the underlying category numbers to the client — the
client doesn't compute the comparison numbers itself, same principle as
the live meter.

### 3.6 Expense history

The client can request the expense list for the current cycle or any
past one. "Archived" (PRD 5.5) isn't a distinct data-movement step —
there's no migration, no separate cold-storage table, nothing that
happens at the moment a cycle ends. A cycle is "archived" purely because
its window has closed (today's date has moved past it); the row and its
expenses just sit where they always were, fully queryable, indefinitely.
The backend's job is serving that list back — paginated per cycle — with
the exact same code path for a live cycle and a cycle from six months
ago, not two different ones.

## 4. Client/backend boundary — what stays on-device

To keep the app fast, cheap to run, and free of holding sensitive raw
input it has no product reason to keep, the following stay client-side
rather than round-tripping to the server:

- Receipt text extraction (on-device OCR).
- Speech-to-text transcription (on-device).
- AI calls for bring-your-own-key users (client talks to the provider
  directly; the backend never sees that traffic or that key).

None of the raw material behind these — the receipt photo, the raw OCR
text, the audio recording, the raw speech transcript — is ever sent to
or stored by the backend, or persisted on the client past the
confirmation step. The only thing that reaches the database is the
final, user-confirmed expense record (merchant/title, amount, category,
date, note) — the same shape regardless of which of the three capture
paths produced it. This is both a cost/scope boundary and a deliberate
privacy stance: there's no reason for Nokori's database to ever hold a
photo of a receipt or a recording of a voice, only the expense it
resolved to.

## 5. Infrastructure trajectory

Nokori's MVP is not a demo or a throwaway proof of concept — per the
PRD's success criteria, it's meant to carry a real, full salary cycle
of the developer's actual expenses with no other app or spreadsheet as
a fallback. That sets a floor even at MVP size: no losing expense data,
no downtime that blocks logging a spend mid-day. Railway is the MVP
hosting choice because it's the fastest path to a real running
instance, not because durability or correctness are optional at this
stage — those are the "must-have, now" half of production quality;
redundancy, horizontal scale, and formal observability tooling are the
"deliberately later" half. The intended direction for this system, once
the product itself is proven out through daily real use, is a properly
designed production environment: containerized, deployed on real cloud
infrastructure, with the full non-functional characteristics a
production system is expected to have — redundancy, observability, a
real CI/CD pipeline, and database design that holds up under realistic
load rather than MVP-scale data. The specific target infrastructure and
the path to it is a decision to be made deliberately later, with its
own design work, not defaulted into this document.

## 6. Implementation decisions deferred past this document

None of these are being silently skipped — each has an explicit,
deliberate disposition, they're just not planning-document decisions:

- **Exact schema for accounts / cycles / expenses, and indexing
  strategy.** Not yet decided — this is real backend design work
  (table/column layout, keys, which columns need a database index for
  fast lookups) that happens once implementation starts, not before.
- **API contract between client and backend** (REST shape, error
  format). To be worked out together once real implementation begins —
  deliberately not needed before then.
- **Concurrency handling for the live meter** (e.g. a photo-derived
  expense landing after a manual entry already updated the balance). On
  reflection, this isn't really a standing architectural risk that needs
  active design work — for a single-user app, it comes down to writing
  the expense-insert code path correctly once (an atomic balance update,
  or deriving the remaining balance from a `SUM()` over expense rows
  instead of maintaining a separately-mutated counter). Normal
  implementation correctness, not something to keep watching.
- **Where the AI usage counts are tracked and how they're reset on
  cycle rollover** — both the monthly-analysis count (PRD 5.6) and the
  cycle-comparison count (PRD 5.7), which are separate counters against
  separate caps. Folds into the schema design above — concretely, this
  is "which table/columns hold the counts." If they live as columns on
  the cycle row itself (one row per cycle, per PRD 5.2), the "reset on
  rollover" problem disappears for free: a new cycle is a new row, which
  starts both counters at zero without any extra reset logic. Still left
  to schema design to confirm, but likely not a hard problem once there.
