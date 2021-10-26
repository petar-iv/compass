BEGIN;

DROP INDEX IF EXISTS labels_tenant_id_key_coalesce_idx;
CREATE UNIQUE INDEX labels_tenant_id_key_coalesce_coalesce1_coalesce2_idx ON labels (tenant_id, key,
                                                                                     coalesce(app_id, '00000000-0000-0000-0000-000000000000'),
                                                                                     coalesce(runtime_id, '00000000-0000-0000-0000-000000000000'),
                                                                                     coalesce(labels.runtime_context_id, '00000000-0000-0000-0000-000000000000'));

ALTER TABLE labels DROP CONSTRAINT labels_bundle_id_fk;
ALTER TABLE labels DROP COLUMN bundle_id;

COMMIT;
