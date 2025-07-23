-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.filters
(
    id uuid NOT NULL,
    sso_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    compare_type text COLLATE pg_catalog."default" NOT NULL,
    value text COLLATE pg_catalog."default" NOT NULL,
    section_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT filters_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.filters
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.filters;
-- +goose StatementEnd
