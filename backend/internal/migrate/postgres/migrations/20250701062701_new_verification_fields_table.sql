-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.verification_fields
(
    id uuid NOT NULL,
    section_id uuid NOT NULL,
    field text COLLATE pg_catalog."default" NOT NULL,
    label text COLLATE pg_catalog."default" DEFAULT ''::text,
    type text COLLATE pg_catalog."default" DEFAULT 'text'::text,
    "position" integer,
    "group" text COLLATE pg_catalog."default" DEFAULT 'form'::text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT verification_fields_pkey PRIMARY KEY (id),
    CONSTRAINT verification_fields_section_id_fkey FOREIGN KEY (section_id)
        REFERENCES public.sections (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.verification_fields
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.verification_fields;
-- +goose StatementEnd
