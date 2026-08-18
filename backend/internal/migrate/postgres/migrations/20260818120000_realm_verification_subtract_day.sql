-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.realms
ADD COLUMN IF NOT EXISTS verification_subtract_day boolean DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.realms
DROP COLUMN IF EXISTS verification_subtract_day;
-- +goose StatementEnd
