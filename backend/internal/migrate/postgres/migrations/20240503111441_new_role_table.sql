-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.roles
(
    id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" DEFAULT ''::text,
    level integer NOT NULL,
    extends uuid[],
    is_show boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT roles_pkey PRIMARY KEY (id)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.roles
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.roles;
-- +goose StatementEnd
