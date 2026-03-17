-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.department_access
(
    id uuid NOT NULL,
    department_id uuid NOT NULL,
    sso_id text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT department_access_pkey PRIMARY KEY (id),
    CONSTRAINT department_access_department_id_fkey FOREIGN KEY (department_id)
        REFERENCES public.departments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.department_access
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.department_access;
-- +goose StatementEnd
