BEGIN;

CREATE TABLE "users"
(
    "id"            bigserial primary key,
    "first_name"    character varying,
    "last_name"     character varying,
    "email"         character varying,
    "created_at"    timestamp NOT NULL,
    "updated_at"    timestamp NOT NULL
);

COMMIT;
