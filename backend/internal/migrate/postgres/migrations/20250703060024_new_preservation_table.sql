-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.preservation
(
    id uuid NOT NULL,
    instrument_id uuid NOT NULL,
    date_start integer NOT NULL,
    date_end integer DEFAULT 0,
    notes_start text COLLATE pg_catalog."default" DEFAULT ''::text,
    notes_end text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT preservation_pkey PRIMARY KEY (id),
    CONSTRAINT preservation_instrument_id_fkey FOREIGN KEY (instrument_id)
        REFERENCES public.instruments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.preservation
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.preservation;
-- +goose StatementEnd
