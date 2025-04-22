-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.sections
(
    id uuid NOT NULL,
    realm_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    position integer,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT sections_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.sections
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.sections;
-- +goose StatementEnd
