BEGIN;

ALTER TABLE webhooks
    DROP COLUMN IF EXISTS proxy_url;

ALTER TABLE fetch_requests
    DROP COLUMN IF EXISTS proxy_url;


COMMIT;
