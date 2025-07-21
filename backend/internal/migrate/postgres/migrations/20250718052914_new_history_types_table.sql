-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.history_types
(
    id uuid NOT NULL,
    section_id uuid NOT NULL,
    "group" text COLLATE pg_catalog."default" NOT NULL,
    label text COLLATE pg_catalog."default" NOT NULL,
    "position" integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT history_types_pkey PRIMARY KEY (id),
    CONSTRAINT history_types_section_id_fkey FOREIGN KEY (section_id)
        REFERENCES public.sections (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.history_types
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.history_types;
-- +goose StatementEnd
