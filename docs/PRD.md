[toc]

# Nokori (残り) — Product Requirements Document

Status: Draft v0.1
Owner: Sabir
Platform: iOS (Flutter)

## 1. One-line pitch

Nokori is a personal expense tracker built around the way money actually
arrives: not the calendar month, but your salary date. It shows you, at a
glance, what's left until the next one.

## 2. Problem

Most budgeting apps reset on the 1st of the calendar month. That doesn't
match how people actually get paid, especially where salary dates vary
person to person (e.g. the 15th, the 25th, end-of-month). The result is a
budget that never quite lines up with reality, so people stop trusting it
and stop using it.

Separately, logging expenses is friction. Apps that require a rigid manual
entry flow for every purchase get abandoned within weeks. The entry point
needs to be as close to zero-effort as possible.

## 3. Target user (v1)

The primary user for v1 is the developer themself — this is a
dogfooding-first product. The design should hold up for a single, real,
daily user before any thought is given to a wider audience.

The one exception is the shared read-only view (5.8): it's built as part
of MVP scope, but there's no specific second person lined up to use it
yet. It ships code-complete against its spec, not proven through actual
daily two-person use the way the rest of MVP is (see 7).

## 4. Core value proposition

- A budget cycle anchored to *your* salary date, not the calendar month —
  and a budget that's grounded in your actual salary and fixed costs,
  not a number typed in on a guess.
- A live, visible "what's left" meter that moves every time you log a
  spend — the core emotional hook of the app.
- Logging an expense should take one of three near-zero-effort paths:
  type it, photograph the receipt, or say it out loud.
- Once a month, a reflective AI summary of spending behavior — not just a
  restated total, but a read on patterns.
- On demand, an AI comparison of this cycle against the last one — a
  live check-in, not just an end-of-cycle look back.

## 5. MVP scope

### 5.1 Account

- Email/password, self-built — the only sign-in method. No Sign in with
  Apple, no Google, no OAuth of any kind (see decisions.md: Apple only
  requires Sign in with Apple when the app offers a third-party login as
  an option, which this deliberately never does).
- Basic account settings: change salary date, update salary and
  basic-needs amounts (recomputes the suggested budget — see 5.2),
  delete account, manage AI key (below).

### 5.2 Budget cycle

- User sets a recurring salary date (day of month). On setup, and again
  whenever it changes, they also enter their salary and their basic
  needs for the cycle — a fixed, app-defined set of categories (rent,
  electricity/water/gas, transportation, food) with a user-entered
  amount each.
- From salary and basic needs, the app computes a suggested cycle
  budget with a fixed formula, adapted from the well-known 50/30/20
  personal-budgeting rule: `salary − basic needs − (20% of salary set
  aside as savings)` — the actual entered basic needs stand in for the
  rule's assumed 50% "needs" bucket, and what's left after the 20%
  savings cut becomes the budget (see decisions.md for the reasoning).
  This is a deterministic calculation — no AI call is involved.
- The user can edit the suggested number before it becomes the cycle's
  real budget, matching the app's rule that nothing gets silently
  applied from a computed guess (see 5.3). If the edit would eat into
  the savings cut or fall short of the entered basic needs, the app
  shows a caution before letting them confirm — a warning, not a block.
- The suggested budget is only recomputed when salary or basic needs
  actually change, not automatically every cycle.
- If the computed suggestion is zero or negative (basic needs plus the
  savings cut consume the entire salary or more), the app still shows
  that number, with a caution — it's a real signal about the user's
  situation, not something to hide or silently auto-correct. The user
  can still confirm it or edit it.
- The app computes the current cycle window (last salary date → next
  salary date) and tracks total spend against the (possibly
  user-adjusted) budget within it.
- Unspent budget does not carry over between cycles — each new cycle
  starts from its own freshly computed (or manually set) budget, not the
  previous cycle's leftover.
- When a new cycle starts, the previous cycle is archived and a new one
  begins automatically, carrying forward the same salary/basic-needs
  figures — and therefore the same suggested budget — unless the user
  changes them.

### 5.3 Expense capture

Three input paths, all converging into the same underlying expense
record (amount, merchant, category, date, note):

- **Text**: quick manual entry form.
- **Photo**: photograph a receipt; on-device text recognition extracts
  the raw text, which is then structured into an expense record.
- **Voice**: on-device speech-to-text transcribes what the user says;
  the transcript is structured into an expense record the same way as
  the photo path.

In both the photo and voice paths, the user confirms/edits the
extracted result before it's saved — the app should never silently log
a guessed number.

Categories split into two kinds: **basic-needs categories** (the same
fixed list used in 5.2 — rent, electricity/water/gas, transportation,
food) and **discretionary categories** (everything else). This isn't
just labeling — see 5.4 for why it matters to the meter.

### 5.4 Live budget meter

- Home screen shows remaining budget for the current cycle, updating
  immediately as expenses are logged.
- Only discretionary-category expenses (5.3) count against the
  remaining balance shown here. Basic-needs-category expenses are still
  logged for the record, but don't move the meter, since that amount is
  already subtracted upfront by the budget formula (5.2) — counting it
  again here would double it.
- If total discretionary spend exceeds the budget, the meter goes
  negative rather than clamping at zero — an honest signal of
  overspending, consistent with the pacing indicator below.
- Visual indication of pace (e.g. spending faster, on track, or slower
  than the cycle position would suggest) is a committed MVP
  requirement, not a stretch goal — it's what makes the meter a
  behavior-changing signal rather than just a running total.

### 5.5 Expense history

- Chronological list of all logged expenses in the current cycle, with
  access to past (archived) cycles.

### 5.6 AI monthly analysis

- Once per cycle, the user can request an AI-generated summary of their
  spending behavior for that cycle (patterns, notable changes vs. their
  own history — not compared to any external population).
- Users on the app's own AI access get one analysis per cycle.
- Users who supply their own AI provider key are not limited by the
  app and are billed directly by their provider.

### 5.7 AI cycle comparison

A separate, more advanced feature from 5.6 — that's a single-cycle
reflection; this puts the current cycle up against the previous one.

- Can be requested at any time, including mid-cycle — not gated to
  once-per-cycle like 5.6. The point is being able to check "how am I
  doing vs. last cycle" while there's still time to act on it.
- Always compares against the immediately previous cycle (no picker for
  older cycles in v1).
- A mid-cycle request compares spend-to-date this cycle against spend
  *at the same point* in the previous cycle (e.g. day 9 of this cycle
  vs. day 9 of last cycle) — never against the previous cycle's full
  total, which would make the current cycle look artificially good
  simply because it isn't over yet.
- Output is both an AI-written comparison *and* the category-by-category
  numbers it's based on, shown together — not narrative alone.
- Users on the app's own AI access get a separate allowance from 5.6's
  once-per-cycle cap, since this is a different feature with a different
  usage pattern. Starting cap: up to 3 requests per cycle — a rough
  starting point, not a researched number, worth revisiting once there's
  real usage to look at.
- Users who supply their own AI provider key are not limited by the app
  and are billed directly by their provider, same as 5.6.

### 5.8 Shared read-only view

A user can invite another existing Nokori account to view their current
cycle, read-only.

- Visible to the invited viewer: the live meter (remaining budget,
  pacing) only — not itemized expense history, not account settings,
  not the salary/basic-needs figures behind the budget. (This is a
  starting default, not researched against what's actually useful once
  there's a real second user to ask — worth revisiting then.)
- One-directional per invite: inviting someone to view your cycle
  doesn't grant you the reverse. Mutual visibility is just two separate
  invites.
- Either side can revoke the connection at any time.
- The viewer has no write access whatsoever — can't log an expense on
  the other person's behalf, edit anything, or get notified of their
  activity. Pure read access to the meter, nothing else.

## 6. Explicit non-goals for v1

To keep the MVP shippable, the following are deliberately out of scope
and deferred to later phases:

- Bank account linking / automatic transaction import — not really a
  "defer to later" item like the rest of this list. The aggregators
  that make this feasible (Plaid, etc.) require a business agreement;
  personal/individual use isn't eligible. Barring the project becoming
  a registered business, this stays out of scope indefinitely rather
  than "coming eventually."
- Multi-currency support.
- Full shared/family budgets — joint editing, a combined pooled budget,
  anything beyond a one-directional read-only view. (The read-only view
  itself is in MVP scope — see 5.8. This non-goal is about anything past
  that.)
- Android or any non-iOS platform.
- Per-item price history / barcode-based product tracking.
- Push notifications / proactive nudges (may follow shortly after MVP,
  not blocking it).

## 7. Success criteria for MVP

The MVP is "done" when the developer can, for one full real salary
cycle, exclusively use Nokori (no spreadsheet, no other app) to:

1. Set up their salary date, salary, and basic needs, and confirm the
   resulting suggested budget.
2. Log every real expense that cycle via text, photo, or voice.
3. Watch the live meter accurately reflect what's left at any point.
4. Get one AI analysis at the end of the cycle that says something
   genuinely useful about their spending.
5. Use the cycle comparison at least once mid-cycle and get a comparison
   against last cycle that's actually worth checking again.

The shared read-only view (5.8) is a deliberate exception to the above:
"done" for it means it works correctly against its spec (invite, view,
revoke) — there's no second real person to validate it through actual
daily use yet, unlike everything else on this list.

## 8. Future direction (post-MVP, not committed)

- Production-grade backend infrastructure (this is being treated as a
  first-class target, not an afterthought — see the design document for
  the intended trajectory).
- Per-item cost-per-use / price-history tracking for recurring
  purchases.
- Richer shared budgets — mutual/joint visibility, combined budgets —
  beyond the one-directional read-only view already in MVP (5.8).
- Richer, more proactive AI behavior (mid-cycle warnings, not just
  end-of-cycle summaries).
