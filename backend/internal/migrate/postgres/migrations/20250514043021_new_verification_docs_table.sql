-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.verification_docs
(
    id uuid NOT NULL,
    verification_id uuid NOT NULL,
    doc_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT verification_docs_pkey PRIMARY KEY (id),
    CONSTRAINT verification_docs_verification_id_fkey FOREIGN KEY (verification_id)
        REFERENCES public.verifications (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT verification_docs_doc_id_fkey FOREIGN KEY (doc_id)
        REFERENCES public.documents (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.verification_docs
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.verification_docs;
-- +goose StatementEnd
