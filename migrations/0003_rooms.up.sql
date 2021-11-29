BEGIN;

CREATE TABLE rooms (
    "id" bigserial primary key,
    name varchar(40),
    host_id varchar
);

CREATE TABLE users_rooms (
                             id bigserial primary key,
                             user_id bigserial not null,
                             room_id bigserial not null,
                             CONSTRAINT user_id_fk
                                 FOREIGN KEY(user_id) REFERENCES users(id)
                                     ON DELETE SET NULL,
                             CONSTRAINT room_id_fk
                                 FOREIGN KEY(room_id) REFERENCES rooms(id)
                                     ON DELETE CASCADE

);

END;