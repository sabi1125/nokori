[toc]

# Nokori (残り) — Decisions Log

Status: Living document (started 2026-09-14)
Companion to: PRD.md, DDD.md

A running record of real engineering/product decisions made during the
build — what we chose, what else we considered, and why. The PRD and DDD
describe the current state of the plan; this document is chronological
and additive instead — an entry isn't rewritten to erase its original
reasoning just because a later decision changes course. A reversal gets
its own new entry that references the one it replaces.

Each entry: context, the realistic options actually on the table as a
pros/cons table, what was picked, and why. The "why" is the point of
this document — the PRD/DDD already say *what* we're doing.

## Index

| Date | Decision | Chosen | Status |
|---|---|---|---|
| 2026-09-14 | [Auth build order](#auth) | Email/password first, Sign in with Apple added before App Store submission | Superseded — see below |
| 2026-09-14 | [Auth: no OAuth at all](#auth-final) | Self-built email/password is the only auth method, permanently — no Sign in with Apple, no Google | Decided |
| 2026-09-14 | [Cycle budget formula](#budget-formula) | Deterministic 50/30/20-derived formula, no AI call | Decided |
| 2026-09-14 | [Basic needs vs. logged expenses](#double-count) | Category split; basic-needs expenses excluded from the meter | Decided |
| 2026-09-14 | [Cycle rollover](#rollover) | No rollover — every cycle starts fresh | Decided |
| 2026-09-14 | [Overspend display](#overspend) | Meter can go negative | Decided |
| 2026-09-14 | [Non-positive suggested budget](#non-positive) | Show it as-is, with a caution | Decided |
| 2026-09-14 | [AI cycle comparison](#comparison) | Separate feature from the monthly analysis: on-demand, own hosted-key cap, always vs. previous cycle at the same elapsed point | Decided |
| 2026-09-14 | [Database engine](#database) | MySQL | Decided (corrects an earlier mischaracterization — see entry) |
| 2026-09-14 | [Backend framework](#framework) | Echo | Decided |
| 2026-09-14 | [API documentation](#api-docs) | OpenAPI via swaggo + echo-swagger | Decided |
| 2026-09-14 | [Diagram tooling](#diagrams) | Mermaid | Decided |
| 2026-09-14 | [Shared read-only view](#sharing) | Part of MVP: one-directional, meter-only, revocable — no second user to validate against yet | Decided |
| 2026-09-14 | [Frontend build/distribution](#frontend-build) | Build and install locally via Xcode for now; CI deferred | Decided |

---

<a id="auth"></a>
## 2026-09-14 — Auth: build our own email/password system first; add Sign in with Apple later, not now

**Status:** Decided for the build/dev phase. Has an open follow-up below.

### Context

The PRD already commits to Sign in with Apple + email/password as the
MVP's two auth methods (PRD 5.1) — that part isn't in question. This
decision is about *sequencing*: which one actually gets built first, and
whether Sign in with Apple is worth setting up right now.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| A. Build Sign in with Apple now, per the PRD's MVP list | Matches the PRD's stated MVP auth from day one. Offloads password storage, hashing, and reset flows entirely to Apple. | Requires an active, paid Apple Developer Program membership ($99/yr) before the app does anything real yet. Also: Apple's guidelines require offering Sign in with Apple alongside *any* third-party OAuth login — relevant if another social login is ever added (PRD already rules this out for v1). |
| **B. Build email/password ourselves first; add Sign in with Apple later** ✅ chosen | Zero cost to start — no Developer Program membership needed to begin backend work or run the app during development. Doesn't block anything else in the build. | Real auth security work (password hashing, session/token issuance, reset flow) lands on us instead of Apple. Objectively more to build than option A — paying and using Apple's is the technically easier path, cost aside. |

### Decision

Option B. Sign in with Apple gets added before App Store submission — at
that point it stops being optional anyway, since shipping to the App
Store requires a paid Developer Program membership regardless of auth
method, so the cost becomes unavoidable then, just not before it has to
be.

### Rationale

The point isn't "avoid this cost forever" — it's not paying it before
it's actually needed. Email/password is enough to build and test
everything else against during development.

### Open follow-up — resolved

This was flagged as an assumption rather than confirmed: was this a
*build-order* call only, or did it actually change what ships in MVP?
**Resolved, and it turned out to be the latter** — see
[Auth: no OAuth at all](#auth-final) below, which supersedes this entry's
"add Sign in with Apple later" plan entirely.

---

<a id="auth-final"></a>
## 2026-09-14 — Auth, superseding the above: no OAuth at all, ever — self-built email/password is permanent

**Status:** Decided. Supersedes [the auth build-order entry](#auth) above
— that entry planned to add Sign in with Apple before App Store
submission; this entry replaces that plan.

### Context

Re-examining *why* Sign in with Apple seemed required at all: Apple's
Guideline 4.8 only mandates it when an app offers a third-party login
(Google, Facebook, etc.) as an option — it's an "if you offer one
third-party option, offer Apple's too" rule, not a blanket requirement.
Since Google/other social logins were already ruled out for v1 (see
PRD 6), and the app only ever offers its own first-party email/password
system, the rule that justified building Sign in with Apple in the first
place never actually applies.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| A. Add Sign in with Apple before App Store submission, as previously planned | Slightly more familiar sign-in option for a wider future audience, if this ever grows past dogfooding. | Solves a compliance requirement that doesn't actually exist for this app. Pure added cost (Developer Program membership tied to a specific purpose) and auth-code surface area for a constraint that was misread. |
| **B. Self-built email/password only, permanently — no OAuth from anyone** ✅ chosen | Own auth "beats" every OAuth-related Apple requirement by never triggering them — no third-party login exists to require parity with. One auth system to build, secure, and maintain, not two. | Full auth security (hashing, sessions, reset flow) is entirely on us, indefinitely — there's no fallback to a provider's security team ever, not just during a dev phase. |

### Decision

Option B. This replaces the "build Sign in with Apple before shipping"
plan from the earlier entry — there is no plan to add it. PRD 5.1
updated accordingly.

### Rationale

The earlier entry treated the Apple Developer Program cost as something
to defer, assuming Sign in with Apple would eventually be required
regardless. That assumption was wrong: the requirement is conditional on
offering a third-party login at all, and this app deliberately never
does. Once that's true, self-built auth isn't a stopgap — it's simply
correct, and there's nothing to add later.

---

<a id="budget-formula"></a>
## 2026-09-14 — Cycle budget: derive it from salary + basic needs via a fixed formula, not an AI call

**Status:** Decided.

### Context

The PRD originally had the user just type a budget number for the
cycle, informed by nothing the app actually knows. On reflection, that's
a real gap: if the app doesn't know what the user earns, it can't judge
whether a "budget" is realistic. The fix: have the user enter their
salary and their basic (fixed) needs — rent, electricity/water/gas,
transportation, food, a fixed app-defined list — and derive the cycle's
budget from that. Open question: *how* to turn salary + basic needs into
a budget number.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| A. Keep it fully manual (original PRD) | Simplest possible; no new data to collect or model. | Budget stays untethered from anything real — weakens the core meter and the later AI analysis (comparing spend against a number that was never grounded in income is a weaker signal). |
| B. An AI call reasons about salary + basic needs and proposes a budget | Could layer in qualitative advice alongside a number. | Overkill for a solved arithmetic problem — adds cost, latency, and a second AI touchpoint to design (a new quota question) for a number a formula gets exactly as well. |
| **C. A fixed, deterministic formula, adapted from the 50/30/20 rule** ✅ chosen | Instant, free, no new AI touchpoint or quota question. Based on an established, widely-cited personal-budgeting heuristic (50% needs / 30% wants / 20% savings) rather than an invented number — fully explainable. `budget = salary − basic needs − (20% of salary set aside as savings)`; actual entered needs stand in for the rule's assumed 50%, and what's left after the 20% savings cut is the spending budget. | It's a heuristic, not personalized — doesn't account for debt, dependents, or local cost of living the way a more elaborate approach eventually might. |

### Decision

Option C. No AI call is involved in producing the suggested budget.

Also decided alongside it:
- The suggested number is **editable, not authoritative** — same rule as
  expense capture: nothing gets silently applied. If the edit would eat
  into the 20% savings cut or fall short of basic needs, the app shows a
  caution before letting the user confirm (a warning, not a hard block).
- It's **only recomputed when salary or basic needs actually change** —
  not automatically re-run every cycle.

### Rationale

The 50/30/20 split is a genuine, widely-used financial-planning
heuristic, not something invented for this app — it gives the suggested
budget a defensible basis instead of being an arbitrary guess or an
unexplainable AI output. Reserving the AI call for the monthly analysis
(where narrative reasoning is actually the point) keeps that feature's
cost/quota story simple.

---

<a id="double-count"></a>
## 2026-09-14 — Basic needs vs. logged expenses: avoid double-counting against the meter

**Status:** Decided.

### Context

Once the cycle budget already has basic needs (rent, utilities,
transport, food) subtracted out via the formula above, a real bug
appears: if the user also logs a rent payment as a normal expense (text/
photo/voice — nothing stops them), the live meter drains by that amount
*twice* — once implicitly via the formula, once explicitly via the
logged expense. This directly breaks the app's core feature.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| A. Keep basic needs settings-only, never logged as expenses | Simplest — no taxonomy complexity, no exclusion logic needed. | No historical record of exactly when/how much was paid for rent etc. Also doesn't actually *prevent* double-counting, just relies on the user never logging them anyway — nothing in the UI stops it. |
| **B. Split expense categories into basic-needs vs. discretionary; exclude basic-needs-category expenses from the live meter** ✅ chosen | Keeps full expense history, including basic needs. Correctness enforced structurally (by category), not by user memory. Minimal added complexity — one exclusion filter on the meter calculation. | `category` stops being purely descriptive — it now drives real accounting behavior, which has to be enforced consistently on both client and backend. |

### Decision

Option B. Expense categories split into the fixed basic-needs list (same
list used for the budget formula: rent, electricity/water/gas,
transportation, food) and everything else (discretionary). Logging a
basic-needs-category expense keeps a record in history but does not
move the live meter, since that amount is already subtracted upfront by
the budget formula. Only discretionary-category expenses count against
the remaining balance.

### Rationale

The meter is the core hook of the app — its correctness matters more
than keeping the category field a simple free-form label. Enforcing the
exclusion structurally means it can't quietly break just because
future-you forgets not to log rent as a normal expense.

---

<a id="rollover"></a>
## 2026-09-14 — Cycle rollover: unspent budget does not carry over

**Status:** Decided.

### Context

When a cycle ends with money left over, does that leftover carry into
the next cycle's budget, or does the next cycle just start fresh? Not
addressed anywhere until now, despite being core to what the live meter
actually means.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| A. Carry unspent budget into the next cycle | Rewards underspending; arguably more "correct" in a strict accounting sense. | Turns each cycle's budget into a running, compounding number instead of a clean, independent figure. Complicates what "what's left" even means (this cycle, or cumulative?). Real scope add nobody asked for. |
| **B. No rollover — every cycle starts fresh from its own computed/set budget** ✅ chosen | Keeps "what's left" meaning exactly what it says, scoped to the current cycle. Matches the existing model where each cycle is already a clean, independently archived unit (DDD 3.6). Zero added complexity. | Doesn't reward disciplined underspending within the app itself — though it's still fully visible in expense history how much was left over. |

### Decision

Option B.

### Rationale

Matches the mental model already established for cycles elsewhere (each
one is a clean, self-contained window). Rollover is a real feature with
real design questions of its own (does it compound indefinitely? cap?
decay?) that nobody has actually asked for — better added deliberately
later if it turns out to matter than defaulted into the MVP.

---

<a id="overspend"></a>
## 2026-09-14 — Overspend: the meter can go negative

**Status:** Decided.

### Context

If actual discretionary spend exceeds the budget, what does the meter
show? Not previously addressed.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| A. Clamp the meter at zero once spend exceeds budget | Simpler-looking UI; "0" reads as an obvious floor. | Actively hides the real number — throws away exactly the information that matters most (how far over you actually are), which undercuts the whole point of a *live, honest* meter. |
| **B. Let the meter go negative** ✅ chosen | Honest, and free — `remaining = budget − spent` already produces the right answer with no special-casing. Consistent with the pacing indicator (PRD 5.4), which exists precisely to flag this kind of situation. | Needs a bit of UI thought for how a negative "what's left" reads at a glance — a design detail, not an architecture one. |

### Decision

Option B.

### Rationale

Consistent with the meter's whole reason for existing: it's supposed to
tell the truth about where you stand, not a flattering version of it. No
extra backend logic required — negative numbers are just what the
existing formula naturally produces once spend passes budget.

---

<a id="non-positive"></a>
## 2026-09-14 — Non-positive suggested budget: show it as-is, with a caution

**Status:** Decided.

### Context

If basic needs plus the 20% savings cut consume all or more of the
user's salary, the formula (see [Cycle budget formula](#budget-formula))
produces a suggested budget of zero or less. What should the app do with
that result?

### Options considered

| Option | Pros | Cons |
|---|---|---|
| A. Prevent this state — e.g. cap the savings cut, or refuse to compute a suggestion until basic needs are lower than salary | User never sees a confusing zero/negative "budget." | Actively hides a genuinely important fact about the user's finances — exactly the kind of silent-guess behavior the app's own rules elsewhere already reject. |
| **B. Show the computed number as-is (zero or negative), with a caution — the same mechanism already planned for risky edits** ✅ chosen | Reuses an already-planned UI pattern (the caution from the budget-formula decision above). Surfaces a real, important signal instead of masking it. No special-case logic — the formula just runs and the existing warning path catches the result. | Can be a jarring first-run experience if it happens — arguably the point. |

### Decision

Option B.

### Rationale

Consistent with the "never silently apply or hide a guess" principle
already established for expense capture and budget edits — extending
the same caution mechanism to the initial suggestion is simpler than
inventing separate hide/prevent logic, and more honest.

---

<a id="comparison"></a>
## 2026-09-14 — AI cycle comparison: a separate, more advanced feature from the monthly analysis

**Status:** Decided.

### Context

The monthly analysis (PRD 5.6) reflects on a single cycle. A different,
more advanced idea came up separately: let the AI compare the current
cycle against the previous one — not just "what did I spend on," but
"how does this stack up." Four things needed pinning down: what it
outputs, when it can be requested, whether it shares the monthly
analysis's hosted-key cap, and which cycle(s) it compares against.

### Output shape

| Option | Pros | Cons |
|---|---|---|
| **A. AI narrative + the raw category numbers shown together** ✅ chosen | Lets the reader see the story and the data it's built from, not just take the AI's word for it. | Slightly more to build than narrative alone (has to surface the numbers, not just the text). |
| B. Pure AI narrative only | Simpler, matches the existing monthly analysis's output shape. | A comparison is exactly the kind of claim you want to be able to check against the actual numbers. |

### Timing

| Option | Pros | Cons |
|---|---|---|
| **A. Any time, including mid-cycle** ✅ chosen | Lets you check "how am I doing vs. last cycle" while there's still time to act on it — more useful than a look back after the fact. | A mid-cycle request is comparing a partial cycle to a full one unless handled carefully (see Decision below). |
| B. Once per cycle, at the end — same rule as the monthly analysis | Simpler to reason about and to quota. | Loses the "catch it while it's still actionable" value that's the whole point of a comparison. |

### Hosted-key cap

| Option | Pros | Cons |
|---|---|---|
| **A. Separate allowance from the monthly analysis's cap** ✅ chosen | Different feature, different usage pattern (repeatable mid-cycle vs. once at the end) — bundling quotas would conflate two distinct things, same reasoning as the earlier budget-formula decision. | Two caps to track and explain instead of one. |
| B. Shares the monthly analysis's per-cycle cap | One combined AI-usage limit, simpler to explain. | Using the comparison a couple of times mid-cycle would use up the once-per-cycle monthly analysis too — punishes exactly the check-in behavior this feature exists to encourage. |

### Comparison scope

| Option | Pros | Cons |
|---|---|---|
| **A. Always the immediately previous cycle** ✅ chosen | Simplest — no cycle-picker UI, no extra backend query flexibility needed. | Can't compare against an arbitrary older cycle (e.g. 3 cycles ago) in v1. |
| B. User can pick any past cycle | More flexible. | Needs a picker UI and more backend query surface for a v1 feature that doesn't need it yet. |

### Decision

All four "chosen" options above. One more detail decided alongside them,
not originally asked but necessary to make "any time, mid-cycle" actually
correct: **a mid-cycle request compares spend-to-date this cycle against
spend up to the same elapsed point in the previous cycle** — e.g. day 9
of this cycle vs. day 9 of last cycle — never against the previous
cycle's full final total, which would always make the current cycle look
artificially good simply because it isn't over yet.

For the hosted-key cap's actual number: **starting at 3 requests per
cycle**, a rough placeholder rather than a researched figure — worth
revisiting once there's real usage to look at (same spirit as the 20%
savings figure elsewhere in this doc: a reasonable starting point, not
a permanent constant).

### Rationale

Treating this as genuinely separate from the monthly analysis — in
output shape, cap, and timing — reflects that it's a different kind of
feature with a different job: the monthly analysis is a reflective,
once-a-cycle read; the comparison is a check-in tool meant to be used
while the cycle is still unfolding. Giving it the same constraints as
the monthly analysis would have quietly made it useless for that job.

---

<a id="database"></a>
## 2026-09-14 — Database engine: MySQL (correcting an earlier mischaracterization)

**Status:** Decided.

### Context

The DDD had said "Postgres" since before this decisions log existed —
inherited from an earlier session's recap file and written into the DDD
as if settled. It was never actually a decision Sabir made and confirmed
in this project; it got carried forward and treated as locked purely
because it was already written down. When backend stack questions came
up for real, it turned out Sabir's own project scaffold (CodeSeed)
defaults every Go project to MySQL — a direct conflict with what the DDD
claimed, and the prompt to actually settle this properly instead of
patching around it.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| A. Keep Postgres, patch CodeSeed's generated `database.go` to use `gorm.io/driver/postgres` instead of MySQL | Matches what the DDD happened to already say. Small mechanical patch (different import, different DSN format) — not a structural problem. | Solves for a decision that was never actually made deliberately — matching a stale doc for its own sake, not for a reason. |
| **B. MySQL, matching CodeSeed's default as-is** ✅ chosen | No patching needed after every `codeseed init`. Matches the database Sabir already has real, repeated experience with across his other projects (same tool, same engine, same muscle memory). Nothing about Nokori's actual data (accounts, cycles, expenses) needs a Postgres-specific feature. | Diverges from what earlier planning material said — hence this entry, to make the correction explicit rather than silent. |

### Decision

Option B. The DDD's "Database" line is corrected to MySQL.

### Rationale

There was no real technical case for Postgres here — it was inertia from
a document written before this decisions log existed, not a considered
choice. Once actually examined, matching Sabir's own tooling and
experience is the more defensible call, and nothing about this project's
data model needs anything MySQL doesn't already provide.

---

<a id="framework"></a>
## 2026-09-14 — Backend framework: Echo

**Status:** Decided.

### Context

Gin was floated as worth considering. CodeSeed's golang templates
(`router.go.tmpl`, `main.go.tmpl`) hardcode Echo (`labstack/echo`),
though — this is less a fork between equally-available options and more
"does the existing tool's default get used or fought."

### Options considered

| Option | Pros | Cons |
|---|---|---|
| **A. Echo, matching CodeSeed's default** ✅ chosen | Zero setup cost — comes free on every `codeseed init` and every `codeseed create`. Echo and Gin are close enough in capability and performance that there's no real technical reason to pay a switching cost here. | None specific to this project. |
| B. Gin | Also a mature, well-supported Echo alternative, if there were ever a concrete reason to prefer it. | Would mean forking CodeSeed's templates or manually rewriting the generated router/main files after every scaffold command — ongoing friction for no identified benefit. |

### Decision

Option A.

### Rationale

There's no concrete requirement pulling toward Gin specifically — this
is a case where using the tool's default beats introducing friction
against it for a preference with no material backing.

---

<a id="api-docs"></a>
## 2026-09-14 — API documentation: OpenAPI via swaggo + echo-swagger

**Status:** Decided.

### Context

The backend needs some form of API documentation as it's built —
question was format and tooling.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| **A. OpenAPI, generated from handler annotations via swaggo (`swaggo/swag` + `swaggo/echo-swagger`)** ✅ chosen | Spec lives next to the handler code as comments, so it's structurally harder for it to silently drift out of sync. Produces a browsable Swagger UI — useful for exercising the API by hand before the Flutter client exists to do it. Industry-standard format if anything ever needs to consume the spec later (client codegen, external tooling). | One more dependency; annotation comments can get verbose for complex request/response shapes. |
| B. Hand-maintained markdown docs | Consistent with the existing PRD/DDD/decisions.md style; no new dependency. | Has to be manually kept in sync with every route change, with nothing structurally enforcing that — the exact kind of drift the other option avoids by construction. |

### Decision

Option A.

### Rationale

Given the backend is being built with real engineering rigor as a
stated goal (see PRD's target-user framing and the DDD's "MVP, not a
POC" stance), documentation that's structurally tied to the code is
worth the one extra dependency over documentation that relies on
remembering to update it.

---

<a id="diagrams"></a>
## 2026-09-14 — Diagram tooling: Mermaid

**Status:** Decided.

### Context

Diagrams will be needed for schema (ER diagrams), system shape, and flow
documentation as the build progresses — question was Mermaid vs.
PlantUML.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| **A. Mermaid** ✅ chosen | Renders natively in GitHub, GitLab, and Claude Code artifacts straight from a fenced code block in the same markdown files already being written — no separate render step or hosting. Already the tool Sabir uses for CodeSeed's own command-flow docs, so it's an existing habit, not a new one. Covers everything needed here: ER diagrams, flowcharts, sequence diagrams. | Less mature for strict UML fidelity than PlantUML in some diagram types. |
| B. PlantUML | More mature/precise UML semantics for complex diagrams. | Needs a separate renderer (local Java tool or a public server) to actually produce an image — friction this project's all-markdown doc workflow doesn't otherwise have anywhere. |

### Decision

Option A.

### Rationale

Consistency with tooling already in use, plus zero-friction rendering
directly inside the markdown files this project already lives in, beats
PlantUML's extra UML precision for what these diagrams actually need to
do here.

---

<a id="sharing"></a>
## 2026-09-14 — Shared read-only view: part of MVP, one-directional, meter-only

**Status:** Decided.

### Context

"Shared budgets" had been sitting as a flat non-goal since early
planning. Revisited: Sabir wants something between nothing and full
shared budgets — a read-only view of another Nokori user's spending.
Three things needed pinning down: whether a specific second person
exists to build this for, whether it ships in MVP or after, and (not
originally asked, but necessary to actually spec it) what exactly the
viewer gets to see.

### Whether a specific second user exists

| Option | Pros | Cons |
|---|---|---|
| A. Yes, a specific real person would use it soon | Real usage to validate against, same as everything else in MVP. | N/A — not the case here. |
| **B. No specific second user yet — forward-looking** ✅ chosen | Doesn't block the feature on finding someone to test with. | This feature can't be validated by real daily use the way the rest of MVP is (PRD 3, PRD 7) — it ships correct-to-spec, not proven-through-use. Worth being explicit about that gap rather than pretending it's tested the same way. |

### MVP timing

| Option | Pros | Cons |
|---|---|---|
| **A. Part of MVP** ✅ chosen | Gets the invite/grant mechanism built now, while the account/auth model is still being shaped — likely cheaper than retrofitting it later. | Adds real scope (an invite/grant flow, cross-account read access) to a v1 that otherwise has no other user in the picture at all. |
| B. Right after MVP | Keeps v1 scope matching what actually gets validated (solo use only). | Retrofits a cross-account concept onto a data model that was designed single-user-only, potentially more expensive later than building it in from the start. |

### What's visible to the viewer (not originally asked — a default proposed here)

| Option | Pros | Cons |
|---|---|---|
| **A. Live meter only (remaining budget, pacing) — no itemized history, no settings** ✅ chosen (starting default) | Conservative by default, consistent with the app's existing privacy stance (DDD 4) of not exposing more than necessary. Simplest to build — one read endpoint, no history pagination or filtering logic for a second account. | May turn out to be too little to actually be useful once there's a real second user — untested guess, flagged as such rather than presented as researched. |
| B. Live meter + full expense history | More genuinely useful for a partner actually trying to understand spending together. | More surface area (itemized data, categories) exposed to a second account with no real user to weigh that trade-off against yet. |

### Decision

Part of MVP (PRD 5.8). One-directional per invite (inviting someone to
view you doesn't grant the reverse — mutual visibility is just two
invites). Revocable at any time. Read access to the live meter only, no
write access of any kind. The "meter only, not full history" scope is a
starting default, explicitly flagged as worth revisiting once there's an
actual second person to ask.

### Rationale

Building the grant mechanism now, while the account model is still
young, is likely cheaper than bolting cross-account access onto a
single-user-shaped schema later — even without a concrete second user
yet. Keeping what's visible minimal by default matches the privacy
posture already established elsewhere in this project (DDD 4) and gives
the cheapest possible version to actually build and revise later.

---

<a id="frontend-build"></a>
## 2026-09-14 — Frontend build/distribution: local Xcode builds for now, CI deferred

**Status:** Decided.

### Context

Railway (the backend host) doesn't build or distribute iOS apps — that
requires Xcode, which is macOS-only, and a completely different pipeline
from a backend deploy. Needed to decide how builds get onto a phone
during active development.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| **A. Build and install locally via Xcode** ✅ chosen | Zero new tooling — Sabir has a Mac. Fastest possible loop during active development, where rebuilds happen constantly anyway. | Free-Apple-ID signing expires every 7 days, forcing a rebuild/reinstall on that cadence even with no code changes — a real but minor annoyance once development slows down. |
| B. Set up a CI pipeline now (Codemagic / Xcode Cloud / GitHub Actions + Fastlane) | Automated builds from day one; no manual rebuild-on-expiry cycle. | Real upfront setup cost for a problem that isn't actually blocking anything yet — premature given the app doesn't exist as code yet either. |

### Decision

Option A for now. Revisit a real CI pipeline once local builds actually
become the bottleneck (e.g. once daily use matters more than active
development, and the 7-day expiry becomes a real annoyance rather than a
non-issue).

### Rationale

There's no reason to pay a CI setup cost before it's actually needed —
same reasoning already used for the Apple Developer Program membership
in the auth decision. Local builds are free and immediate given Sabir
already has a Mac.
