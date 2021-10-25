BEGIN;

ALTER TABLE labels ADD COLUMN bundle_id uuid;

ALTER TABLE labels
    ADD CONSTRAINT labels_bundle_id_fk
        FOREIGN KEY (bundle_id) REFERENCES bundles(id);

COMMIT;
