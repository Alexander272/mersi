-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.realms 
ADD COLUMN IF NOT EXISTS return_notice boolean DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.realms 
DROP COLUMN IF EXISTS return_notice;
-- +goose StatementEnd
