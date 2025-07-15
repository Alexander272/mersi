-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.transfer_to_save
(
    id uuid NOT NULL,
    instrument_id uuid NOT NULL,
    date_start integer NOT NULL,
    notes_start text COLLATE pg_catalog."default" DEFAULT ''::text,
    date_end integer DEFAULT 0,
    notes_end text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT transfer_to_save_pkey PRIMARY KEY (id),
    CONSTRAINT transfer_to_save_instrument_id_fkey FOREIGN KEY (instrument_id)
        REFERENCES public.instruments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transfer_to_save
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.transfer_to_save;
-- +goose StatementEnd
