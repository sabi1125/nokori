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

## 4. Core value proposition

- A budget cycle anchored to *your* salary date, not the calendar month.
- A live, visible "what's left" meter that moves every time you log a
  spend — the core emotional hook of the app.
- Logging an expense should take one of three near-zero-effort paths:
  type it, photograph the receipt, or say it out loud.
- Once a month, a reflective AI summary of spending behavior — not just a
  restated total, but a read on patterns.

## 5. MVP scope

### 5.1 Account

- Sign in with Apple (required — the only social login offered, which
  keeps the app compliant without needing a second OAuth provider).
- Email/password as a fallback sign-in method.
- Basic account settings: change salary date, change budget amount,
  delete account, manage AI key (below).

### 5.2 Budget cycle

- User sets a recurring salary date (day of month) and a budget amount
  for that cycle.
- The app computes the current cycle window (last salary date → next
  salary date) and tracks total spend against the budget within it.
- When a new cycle starts, the previous cycle is archived and a new one
  begins automatically.

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

### 5.4 Live budget meter

- Home screen shows remaining budget for the current cycle, updating
  immediately as expenses are logged.
- Visual indication of pace (e.g. spending faster than the cycle
  position would suggest) is in scope for v1 if time allows; not a
  blocker for shipping.

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

## 6. Explicit non-goals for v1

To keep the MVP shippable, the following are deliberately out of scope
and deferred to later phases:

- Bank account linking / automatic transaction import.
- Multi-currency support.
- Shared or family/household budgets.
- Android or any non-iOS platform.
- Per-item price history / barcode-based product tracking.
- Push notifications / proactive nudges (may follow shortly after MVP,
  not blocking it).

## 7. Success criteria for MVP

The MVP is "done" when the developer can, for one full real salary
cycle, exclusively use Nokori (no spreadsheet, no other app) to:

1. Set up their salary date and budget.
2. Log every real expense that cycle via text, photo, or voice.
3. Watch the live meter accurately reflect what's left at any point.
4. Get one AI analysis at the end of the cycle that says something
   genuinely useful about their spending.

## 8. Future direction (post-MVP, not committed)

- Production-grade backend infrastructure (this is being treated as a
  first-class target, not an afterthought — see the design document for
  the intended trajectory).
- Per-item cost-per-use / price-history tracking for recurring
  purchases.
- Shared budgets for two people.
- Richer, more proactive AI behavior (mid-cycle warnings, not just
  end-of-cycle summaries).
