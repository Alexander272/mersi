-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users
(
    id uuid NOT NULL,
    sso_id text COLLATE pg_catalog."default" NOT NULL,
    username text COLLATE pg_catalog."default" DEFAULT ''::text,
    first_name text COLLATE pg_catalog."default" DEFAULT ''::text,
    last_name text COLLATE pg_catalog."default" DEFAULT ''::text,
    email text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT users_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users
-- +goose StatementEnd
