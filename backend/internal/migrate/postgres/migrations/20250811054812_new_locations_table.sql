-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.locations
(
    id uuid NOT NULL,
    instrument_id uuid NOT NULL,
    date_of_receiving timestamp without time zone DEFAULT '0001-01-01'::DATE,
    date_of_issue timestamp without time zone DEFAULT '0001-01-01'::DATE,
    status text COLLATE pg_catalog."default" NOT NULL,
    has_confirmed boolean DEFAULT false,
    need_confirmed boolean DEFAULT false,
    person text COLLATE pg_catalog."default" DEFAULT ''::text,
    person_id uuid,
    place text COLLATE pg_catalog."default" DEFAULT ''::text,
    department_id uuid,
    last_place text COLLATE pg_catalog."default" DEFAULT ''::text,
    last_place_id text COLLATE pg_catalog."default" DEFAULT ''::text,
    user_id uuid,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT locations_pkey PRIMARY KEY (id),
    CONSTRAINT locations_instrument_id_fkey FOREIGN KEY (instrument_id)
        REFERENCES public.instruments (id) MATCH SIMPLE
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.locations
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.locations;
-- +goose StatementEnd
