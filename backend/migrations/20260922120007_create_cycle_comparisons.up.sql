CREATE TABLE cycle_comparisons (
    comparison_id        CHAR(36) NOT NULL PRIMARY KEY,
    budget_cycle_id      CHAR(36) NOT NULL,
    comparison_text      TEXT     NOT NULL,
    category_breakdown   JSON     NOT NULL,
    comparison_count     INT      NOT NULL,
    created_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cycle_comparisons_budget_cycle
        FOREIGN KEY (budget_cycle_id) REFERENCES budget_cycles (budget_cycle_id)
        ON DELETE CASCADE,
    KEY idx_cycle_comparisons_budget_cycle (budget_cycle_id)
);
