-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.rule_item
(
    id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    method text COLLATE pg_catalog."default" DEFAULT 'read'::text,
    description text COLLATE pg_catalog."default" DEFAULT ''::text,
    is_show boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT rule_item_pkey PRIMARY KEY (id)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.rule_item
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.rule_item;
-- +goose StatementEnd
