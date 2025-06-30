-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.tools_menu
(
    id uuid NOT NULL,
    "position" integer NOT NULL,
    section_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    label text COLLATE pg_catalog."default" DEFAULT ''::text,
    rule_item_id uuid NOT NULL,
    can_be_favorite boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT tools_menu_pkey PRIMARY KEY (id),
    CONSTRAINT tools_menu_rule_item_id_fkey FOREIGN KEY (rule_item_id)
        REFERENCES public.rule_item (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT tools_menu_section_id_fkey FOREIGN KEY (section_id)
        REFERENCES public.sections (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tools_menu
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.tools_menu;
-- +goose StatementEnd
