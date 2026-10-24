CREATE TABLE sessions (
    session_id            CHAR(36)  NOT NULL PRIMARY KEY,
    user_id               CHAR(36)  NOT NULL,
    token_hash            CHAR(64)  NOT NULL,
    previous_token_hash   CHAR(64)  NULL,
    previous_valid_until  DATETIME  NULL,
    expires_at            DATETIME  NOT NULL,
    created_at            DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_sessions_token_hash (token_hash),
    UNIQUE KEY uq_sessions_previous_token_hash (previous_token_hash),
    CONSTRAINT fk_sessions_user
        FOREIGN KEY (user_id) REFERENCES users (user_id)
        ON DELETE CASCADE
);
