ALTER TABLE verification_codes
    ADD COLUMN attempts INT NOT NULL DEFAULT 0 AFTER expires_at;
