# Database schema

```mermaid
erDiagram
    USER ||--o| SETTINGS : has
    USER ||--o{ VERIFICATION_CODE : has
    USER ||--o{ BUDGET_CYCLE : has
    BUDGET_CYCLE ||--o{ EXPENSES_HISTORY : has
    BUDGET_CYCLE ||--o{ MONTHLY_ANALYSIS : has
    BUDGET_CYCLE ||--o{ CYCLE_COMPARISON : has

    USER {
        int id PK
        string first_name
        string last_name
        string email
        string password_hash
        bool verified
        datetime created_at
        datetime updated_at
    }

    VERIFICATION_CODE {
        int id PK
        int user_id FK
        string code
        datetime expires_at
        datetime created_at
    }

    SETTINGS {
        int id PK
        int user_id FK
        int default_salary_date
        datetime created_at
        datetime updated_at
    }

    BUDGET_CYCLE {
        int id PK
        int user_id FK
        date start_month
        date end_month
        int salary
        json basic_needs
        int budget
        datetime created_at
        datetime updated_at
    }

    EXPENSES_HISTORY {
        int id PK
        int budget_cycle_id FK
        int spent
        string item_name
        string category
        date expense_date
        datetime created_at
    }

    MONTHLY_ANALYSIS {
        int id PK
        int budget_cycle_id FK
        text ai_analysis
        text recommendations
        datetime created_at
    }

    CYCLE_COMPARISON {
        int id PK
        int budget_cycle_id FK
        text comparison_text
        json category_breakdown
        int comparison_count
        datetime created_at
    }
```

## Notes

- `SETTINGS` is 1:1 with `USER` — just `default_salary_date`. Salary
  and basic needs live on `BUDGET_CYCLE` instead (see decisions.md
  "Basic needs: ship a default set, not a fixed list").
- `VERIFICATION_CODE` rows are deleted by app logic on use or expiry —
  no `used` flag needed (see decisions.md "Email verification").
- `BUDGET_CYCLE.basic_needs` is JSON: a variable, user-editable list of
  `{name, amount}` items, not a fixed set of columns.
- `MONTHLY_ANALYSIS` has no cap-counter column: capped at 1/cycle, so
  "does a row exist for this `budget_cycle_id`" is the check.
  `CYCLE_COMPARISON.comparison_count` exists because its cap is 3/cycle,
  where row-presence alone can't tell you which count you're at.
- No `EXPENSES_HISTORY`/`MONTHLY_ANALYSIS`/`CYCLE_COMPARISON` table
  carries its own `user_id` — each reaches the user through
  `budget_cycle_id` instead.
