-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.realms
(
    id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    realm text COLLATE pg_catalog."default" NOT NULL,
    is_active boolean DEFAULT true,
    notification_channel text COLLATE pg_catalog."default" DEFAULT ''::text,
    expiration_notice boolean DEFAULT false,
    has_locations boolean DEFAULT false,
    location_type text COLLATE pg_catalog."default" DEFAULT 'department'::text,
    need_confirmed boolean NOT NULL DEFAULT true,
    has_employees boolean NOT NULL DEFAULT true,
    has_responsible boolean NOT NULL DEFAULT true,
    has_commissioning_cert boolean DEFAULT false,
    has_preservations boolean DEFAULT false,
    has_transfer boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT realms_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.realms
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.realms;
-- +goose StatementEnd
