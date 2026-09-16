CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL ,
    caused_by BIGINT NOT NULL REFERENCES users(id),
    received_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_read BOOLEAN NOT NULL DEFAULT FALSE
);