CREATE TABLE settings (
    settings_id          CHAR(36) NOT NULL PRIMARY KEY,
    user_id              CHAR(36) NOT NULL,
    default_salary_date  INT      NOT NULL,
    created_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_settings_user
        FOREIGN KEY (user_id) REFERENCES users (user_id)
        ON DELETE CASCADE,
    CONSTRAINT chk_settings_salary_date
        CHECK (default_salary_date BETWEEN 1 AND 31),
    UNIQUE KEY uq_settings_user_id (user_id)
);
