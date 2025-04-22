-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.columns
(
    id uuid NOT NULL,
    section_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    field text COLLATE pg_catalog."default" NOT NULL,
    position integer NOT NULL,
    type text COLLATE pg_catalog."default" NOT NULL,
    width integer DEFAULT 200,
    parent_id uuid DEFAULT NULL::uuid,
    allow_sort boolean DEFAULT true,
    allow_filter boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT columns_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.columns
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.columns;
-- +goose StatementEnd
