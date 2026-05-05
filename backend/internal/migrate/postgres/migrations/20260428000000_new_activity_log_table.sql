-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.activity_log
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    table_name text NOT NULL,
    record_id uuid NOT NULL,
    record_name text NOT NULL,
    action text NOT NULL,
    field_name text,
    old_value jsonb,
    new_value jsonb,
    user_id uuid NOT NULL,
    user_name text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT activity_log_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.activity_log
    OWNER to postgres;

CREATE INDEX IF NOT EXISTS idx_activity_log_table_record ON public.activity_log (table_name, record_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_user ON public.activity_log (user_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_created ON public.activity_log (created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.activity_log;
-- +goose StatementEnd
