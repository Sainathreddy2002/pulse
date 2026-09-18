ALTER TABLE outbox
    ADD COLUMN attempts INT NOT NULL DEFAULT 0;

UPDATE outbox SET status = 'pending' WHERE status = 'processing';
