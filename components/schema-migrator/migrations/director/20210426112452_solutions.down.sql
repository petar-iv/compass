BEGIN;

ALTER TABLE labels
    DROP COLUMN solution_id UUID;

ALTER TABLE labels
    DROP CONSTRAINT solution_id_fk FOREIGN KEY (solution_id) REFERENCES solutions(id) ON DELETE CASCADE;

ALTER TABLE labels
DROP CONSTRAINT valid_refs;

ALTER TABLE labels
    ADD CONSTRAINT valid_refs
        CHECK (app_id IS NOT NULL OR runtime_id IS NOT NULL OR labels.runtime_context_id IS NOT NULL);

DROP INDEX IF EXISTS labels_tenant_id_key_coalesce_coalesce1_coalesce2_idx;
CREATE UNIQUE INDEX ON labels (tenant_id, key, coalesce(app_id, '00000000-0000-0000-0000-000000000000'), coalesce(runtime_id, '00000000-0000-0000-0000-000000000000'), coalesce(labels.runtime_context_id, '00000000-0000-0000-0000-000000000000'));

DROP TABLE solutions;

COMMIT;
