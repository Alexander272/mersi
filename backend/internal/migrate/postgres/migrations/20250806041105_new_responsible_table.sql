-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.responsible
(
    id uuid NOT NULL,
    department_id uuid NOT NULL,
    sso_id text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT responsible_pkey PRIMARY KEY (id),
    CONSTRAINT responsible_department_id_fkey FOREIGN KEY (department_id)
        REFERENCES public.departments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.responsible
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.responsible;
-- +goose StatementEnd
