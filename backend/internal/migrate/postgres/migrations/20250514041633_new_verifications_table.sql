-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.verifications
(
    id uuid NOT NULL,
    instrument_id uuid NOT NULL,
    register_link text COLLATE pg_catalog."default" DEFAULT ''::text,
    status text COLLATE pg_catalog."default" DEFAULT ''::text,
    date integer DEFAULT 0,
    next_date integer DEFAULT 0,
    notes text COLLATE pg_catalog."default" DEFAULT ''::text,
    not_verified boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT verifications_pkey PRIMARY KEY (id),
    CONSTRAINT verifications_instrument_id_fkey FOREIGN KEY (instrument_id)
        REFERENCES public.instruments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.verifications
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.verifications;
-- +goose StatementEnd
