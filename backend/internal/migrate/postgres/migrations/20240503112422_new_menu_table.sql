-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.menu
(
    id uuid NOT NULL,
    role_id uuid NOT NULL,
    menu_item_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT menu_pkey PRIMARY KEY (id),
    CONSTRAINT menu_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES public.roles (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT menu_menu_item_id_fkey FOREIGN KEY (menu_item_id)
        REFERENCES public.menu_item (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.menu
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.menu;
-- +goose StatementEnd
