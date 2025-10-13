-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.si_status
(
    id uuid NOT NULL,
    section_id uuid NOT NULL,
    "position" integer NOT NULL,
    value text COLLATE pg_catalog."default" NOT NULL,
    label text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT si_status_pkey PRIMARY KEY (id),
    CONSTRAINT si_status_section_id_fkey FOREIGN KEY (section_id)
        REFERENCES public.sections (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.si_status
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.si_status;
-- +goose StatementEnd
