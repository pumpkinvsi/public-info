-- =============================================================
-- 000004_init_outbox.up.sql
-- Initializes the outbox table.
-- =============================================================

CREATE TABLE IF NOT EXISTS outbox (
    id         SERIAL      PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    payload    JSONB       NOT NULL,
    processed  BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS outbox_unsent_idx ON outbox (id) WHERE sent_at IS NULL;