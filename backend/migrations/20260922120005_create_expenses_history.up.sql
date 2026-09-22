CREATE TABLE expenses_history (
    expenses_history_id  CHAR(36)     NOT NULL PRIMARY KEY,
    budget_cycle_id      CHAR(36)     NOT NULL,
    spent                INT          NOT NULL,
    item_name            VARCHAR(255) NOT NULL,
    category             VARCHAR(255) NOT NULL,
    expense_date         DATE         NOT NULL,
    created_at           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_expenses_history_budget_cycle
        FOREIGN KEY (budget_cycle_id) REFERENCES budget_cycles (budget_cycle_id)
        ON DELETE CASCADE,
    KEY idx_expenses_history_budget_cycle (budget_cycle_id)
);
