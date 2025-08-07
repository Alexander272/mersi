-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.employee
(
    id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    department_id uuid NOT NULL,
    most_id text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    is_lead boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    sso_id text COLLATE pg_catalog."default" DEFAULT ''::text,
    channel_id text COLLATE pg_catalog."default" DEFAULT ''::text,
    notes text COLLATE pg_catalog."default" DEFAULT ''::text,
    CONSTRAINT employee_pkey PRIMARY KEY (id),
    CONSTRAINT employee_department_id_fkey FOREIGN KEY (department_id)
        REFERENCES public.departments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.employee
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.employee;
-- +goose StatementEnd
