-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.sorting
(
    id uuid NOT NULL,
    sso_id uuid NOT NULL,
    section_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    order_type text COLLATE pg_catalog."default" DEFAULT 'ASC'::text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT sorting_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.sorting
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.sorting;
-- +goose StatementEnd
