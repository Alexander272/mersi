-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.instruments
(
    id uuid NOT NULL,
    section_id uuid NOT NULL,
    user_id uuid,
    position integer,
    name text COLLATE pg_catalog."default" NOT NULL,
    date_of_receipt integer NOT NULL,
    type text COLLATE pg_catalog."default" DEFAULT ''::text,
    factory_number text COLLATE pg_catalog."default" NOT NULL,
    measurement_limits text COLLATE pg_catalog."default" DEFAULT ''::text,
    accuracy text COLLATE pg_catalog."default" DEFAULT ''::text,
    state_register text COLLATE pg_catalog."default" DEFAULT ''::text,
    country_of_produce text COLLATE pg_catalog."default" DEFAULT ''::text,
    manufacturer text COLLATE pg_catalog."default" DEFAULT ''::text,
    responsible text COLLATE pg_catalog."default" DEFAULT ''::text,
    inventory text COLLATE pg_catalog."default" DEFAULT ''::text,
    year_of_issue integer DEFAULT 0,
    inter_verification_interval integer DEFAULT 0,
    act_of_entering text COLLATE pg_catalog."default" DEFAULT ''::text,
    act_of_entering_id uuid,
    notes text COLLATE pg_catalog."default" DEFAULT ''::text,
    status text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted timestamp with time zone,
    CONSTRAINT instruments_pkey PRIMARY KEY (id),
    CONSTRAINT instruments_section_id_fkey FOREIGN KEY (section_id)
        REFERENCES public.sections (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.instruments
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.instruments;
-- +goose StatementEnd
