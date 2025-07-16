-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.transfer_to_department
(
    id uuid NOT NULL,
    instrument_id uuid NOT NULL,
    date integer NOT NULL,
    notes text COLLATE pg_catalog."default" DEFAULT ''::text,
    doc_id uuid,
    doc_name text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT transfer_to_department_pkey PRIMARY KEY (id),
    CONSTRAINT transfer_to_department_instrument_id_fkey FOREIGN KEY (instrument_id)
        REFERENCES public.instruments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transfer_to_department
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.transfer_to_department;
-- +goose StatementEnd
