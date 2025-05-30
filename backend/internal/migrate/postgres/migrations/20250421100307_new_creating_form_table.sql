-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.creating_form
(
    id uuid NOT NULL,
    section_id uuid NOT NULL,
    step integer NOT NULL,
    step_name text COLLATE pg_catalog."default" NOT NULL,
    field text COLLATE pg_catalog."default" NOT NULL,
    field_name text COLLATE pg_catalog."default" NOT NULL,
    path text COLLATE pg_catalog."default" DEFAULT ''::text,
    type text COLLATE pg_catalog."default" NOT NULL,
    is_required boolean DEFAULT true,
    position integer NOT NULL,
    min integer,
    max integer,
    multiline boolean DEFAULT false,
    rows integer,
    disableFuture boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT creating_form_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.creating_form
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.creating_form;
-- +goose StatementEnd
