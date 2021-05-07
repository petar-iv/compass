BEGIN;

ALTER TABLE bundle_instance_auths
    DROP COLUMN ready,
    DROP COLUMN created_at,
    DROP COLUMN updated_at,
    DROP COLUMN deleted_at,
    DROP COLUMN error;

ALTER TABLE webhooks
    DROP COLUMN result_template jsonb;

ALTER TABLE webhooks
ALTER COLUMN type TYPE VARCHAR(255);

DROP TYPE webhook_type;

CREATE TYPE webhook_type AS ENUM (
    'CONFIGURATION_CHANGED',
    'REGISTER_APPLICATION',
    'UNREGISTER_APPLICATION',
    'OPEN_RESOURCE_DISCOVERY'
    );

ALTER TABLE webhooks
ALTER COLUMN type TYPE webhook_type USING (type::webhook_type);

COMMIT;
