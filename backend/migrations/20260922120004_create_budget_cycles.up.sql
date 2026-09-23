CREATE TABLE budget_cycles (
    budget_cycle_id  CHAR(36) NOT NULL PRIMARY KEY,
    user_id          CHAR(36) NOT NULL,
    start_month      DATE     NOT NULL,
    end_month        DATE     NOT NULL,
    salary           INT      NOT NULL,
    basic_needs      JSON     NOT NULL,
    budget           INT      NOT NULL,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_budget_cycles_user
        FOREIGN KEY (user_id) REFERENCES users (user_id)
        ON DELETE CASCADE,
    KEY idx_budget_cycles_user_window (user_id, start_month, end_month)
);
