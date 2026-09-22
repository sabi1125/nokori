CREATE TABLE monthly_analyses (
    analysis_id       CHAR(36) NOT NULL PRIMARY KEY,
    budget_cycle_id   CHAR(36) NOT NULL,
    ai_analysis       TEXT     NOT NULL,
    recommendations   TEXT     NOT NULL,
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_monthly_analyses_budget_cycle
        FOREIGN KEY (budget_cycle_id) REFERENCES budget_cycles (budget_cycle_id)
        ON DELETE CASCADE,
    UNIQUE KEY uq_monthly_analyses_budget_cycle (budget_cycle_id)
);
