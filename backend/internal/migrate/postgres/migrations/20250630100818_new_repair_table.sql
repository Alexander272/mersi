-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.repair
(
    id uuid NOT NULL,
    instrument_id uuid NOT NULL,
    defect text COLLATE pg_catalog."default" DEFAULT ''::text,
    work text COLLATE pg_catalog."default" DEFAULT ''::text,
    period_start integer DEFAULT 0,
    period_end integer DEFAULT 0,
    description text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT repair_pkey PRIMARY KEY (id),
    CONSTRAINT repair_instrument_id_fkey FOREIGN KEY (instrument_id)
        REFERENCES public.instruments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.repair
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.repair;
-- +goose StatementEnd
