BEGIN;

CREATE TABLE "tracks" (
    "id" bigserial primary key,
    track_url varchar not null,
    added_on timestamp not null default now(),
    added_by bigserial not null,
        CONSTRAINT added_by_fk
        FOREIGN KEY(added_by) REFERENCES users(id)
        ON DELETE SET NULL
);

END;