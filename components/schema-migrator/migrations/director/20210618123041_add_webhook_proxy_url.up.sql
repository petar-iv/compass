BEGIN;

ALTER TABLE webhooks
    ADD COLUMN proxy_url VARCHAR(255);

ALTER TABLE fetch_requests
    ADD COLUMN proxy_url VARCHAR(255);

COMMIT;
