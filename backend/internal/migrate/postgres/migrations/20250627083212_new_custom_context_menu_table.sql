-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.custom_context_menu
(
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    tools_menu_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT custom_context_menu_pkey PRIMARY KEY (id),
    CONSTRAINT custom_context_menu_tools_menu_id_fkey FOREIGN KEY (tools_menu_id)
        REFERENCES public.tools_menu (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.custom_context_menu
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.custom_context_menu;
-- +goose StatementEnd
