CREATE TABLE IF NOT EXISTS auth_user
(
    id        BIGSERIAL PRIMARY KEY,
    email     VARCHAR     NOT NULL,
    pass_hash VARCHAR     NOT NULL,
    UNIQUE (email)
);
