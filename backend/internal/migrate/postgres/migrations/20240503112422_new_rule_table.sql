-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.rule
(
    id uuid NOT NULL,
    role_id uuid NOT NULL,
    rule_item_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT rule_pkey PRIMARY KEY (id),
    CONSTRAINT rule_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES public.roles (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT rule_rule_item_id_fkey FOREIGN KEY (rule_item_id)
        REFERENCES public.rule_item (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.rule
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.rule;
-- +goose StatementEnd
