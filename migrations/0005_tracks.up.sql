BEGIN;

CREATE TABLE tracks (
    "id" bigserial primary key,
    track_url varchar,
    added_by bigserial
        CONSTRAINT added_by_fk
        FOREIGN KEY(added_by) REFERENCES users(id)
        ON DELETE SET NULL,
);

END;