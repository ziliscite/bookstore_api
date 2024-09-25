CREATE TABLE Sessions (
    id            VARCHAR(255) PRIMARY KEY NOT NULL,
    user_email    VARCHAR(255)             NOT NULL,
    refresh_token VARCHAR(512)             NOT NULL,
    is_revoked    BOOLEAN   DEFAULT FALSE  NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW(),
    expires_at    TIMESTAMP NOT NULL
)