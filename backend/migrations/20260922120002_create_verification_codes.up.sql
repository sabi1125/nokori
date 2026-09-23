CREATE TABLE verification_codes (
    verification_code_id  CHAR(36)     NOT NULL PRIMARY KEY,
    user_id               CHAR(36)     NOT NULL,
    code                  VARCHAR(255) NOT NULL,
    expires_at            DATETIME     NOT NULL,
    created_at            DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_verification_codes_user
        FOREIGN KEY (user_id) REFERENCES users (user_id)
        ON DELETE CASCADE
);
