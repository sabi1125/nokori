ALTER TABLE verification_codes
    ADD COLUMN resend_email_id VARCHAR(255) NULL AFTER code;
