-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.context_menu
(
    id uuid NOT NULL,
    "position" integer NOT NULL,
    section_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    label text COLLATE pg_catalog."default" DEFAULT ''::text,
    rule_item_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT context_menu_pkey PRIMARY KEY (id),
    CONSTRAINT context_menu_rule_item_id_fkey FOREIGN KEY (rule_item_id)
        REFERENCES public.rule_item (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT context_menu_section_id_fkey FOREIGN KEY (section_id)
        REFERENCES public.sections (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.context_menu
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.context_menu;
-- +goose StatementEnd
