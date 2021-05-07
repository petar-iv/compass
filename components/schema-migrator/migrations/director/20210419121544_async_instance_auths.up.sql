BEGIN;

ALTER TABLE bundle_instance_auths
    ADD COLUMN ready bool DEFAULT TRUE,
    ADD COLUMN created_at timestamp,
    ADD COLUMN updated_at timestamp,
    ADD COLUMN deleted_at timestamp,
    ADD COLUMN error jsonb;

ALTER TABLE webhooks
    ADD COLUMN result_template jsonb;

ALTER TABLE webhooks
ALTER COLUMN type TYPE VARCHAR(255);

DROP TYPE webhook_type;

CREATE TYPE webhook_type AS ENUM (
    'CONFIGURATION_CHANGED',
    'REGISTER_APPLICATION',
    'UNREGISTER_APPLICATION',
    'OPEN_RESOURCE_DISCOVERY',
    'BUNDLE_INSTANCE_AUTH_CREATION',
    'BUNDLE_INSTANCE_AUTH_DELETION'
);

ALTER TABLE webhooks
ALTER COLUMN type TYPE webhook_type USING (type::webhook_type);

COMMIT;
