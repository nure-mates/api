BEGIN;

CREATE TABLE "tracks" (
    "id" bigserial primary key,
    track_url varchar not null,
    added_on timestamp not null default now(),
    room_id int not null,
    added_by bigserial not null,
        CONSTRAINT added_by_fk
        FOREIGN KEY(added_by) REFERENCES users(id)
        ON DELETE SET NULL,
    CONSTRAINT room_id_fk
    FOREIGN KEY(room_id) REFERENCES rooms(id)
    ON DELETE CASCADE
);

END;