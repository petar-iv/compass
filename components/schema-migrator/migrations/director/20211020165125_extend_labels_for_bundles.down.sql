BEGIN;

ALTER TABLE labels DROP CONSTRAINT labels_bundle_id_fk;
ALTER TABLE labels DROP COLUMN bundle_id;

COMMIT;
