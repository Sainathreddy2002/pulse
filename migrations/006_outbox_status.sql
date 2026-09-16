ALTER TABLE outbox
    ADD COLUMN status TEXT NOT NULL DEFAULT 'pending';
