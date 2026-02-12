-- +goose Up
-- +goose StatementBegin
-- 1. Индекс для instruments (фильтрация по section_id и status)
CREATE INDEX IF NOT EXISTS idx_instruments_section_status
ON instruments (section_id, status);

-- 2. Индекс для verification (для получения последней записи по instrument_id)
CREATE INDEX IF NOT EXISTS idx_verification_instrument_date
ON verifications (instrument_id, date DESC, created_at DESC);

-- 3. Индекс для verification_docs (для получения документа по verification_id)
CREATE INDEX IF NOT EXISTS idx_verification_docs_verification_id
ON verification_docs (verification_id);

-- 4. Индекс для repair (для получения последней записи по instrument_id)
CREATE INDEX IF NOT EXISTS idx_repair_instrument_period
ON repair (instrument_id, period_start DESC);

-- 5. Индекс для preservation (для получения последней записи по instrument_id)
CREATE INDEX IF NOT EXISTS idx_preservation_instrument_date
ON preservation (instrument_id, date_start DESC);

-- 6. Индекс для transfer_to_save (для получения последней записи по instrument_id)
CREATE INDEX IF NOT EXISTS idx_transfer_to_save_instrument_date
ON transfer_to_save (instrument_id, date_start DESC);

-- 7. Индекс для transfer_to_dep (для получения последней записи по instrument_id)
CREATE INDEX IF NOT EXISTS idx_transfer_to_dep_instrument_date
ON transfer_to_department (instrument_id, date DESC);

-- 8. Индекс для write_off (для получения последней записи по instrument_id)
CREATE INDEX IF NOT EXISTS idx_write_off_instrument_date
ON write_off (instrument_id, date DESC);

-- 9. Индекс для location (для получения последней записи по instrument_id)
CREATE INDEX IF NOT EXISTS idx_location_instrument_dates
ON locations (instrument_id, date_of_issue DESC, created_at DESC);

-- 10. Индекс для employee (поиск по person_id)
CREATE INDEX IF NOT EXISTS idx_employee_id
ON employee (id);

-- 11. Индекс для department (поиск по department_id и last_place_id)
CREATE INDEX IF NOT EXISTS idx_department_id
ON departments (id);

-- 12. Индекс для location.last_place_id (JOIN с department)
CREATE INDEX IF NOT EXISTS idx_location_last_place_id
ON locations (last_place_id);

-- 13. Индекс для location.person_id (JOIN с employee)
CREATE INDEX IF NOT EXISTS idx_location_person_id
ON locations (person_id);

-- 14. Индекс для location.department_id (JOIN с department)
CREATE INDEX IF NOT EXISTS idx_location_department_id
ON locations (department_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_location_department_id;
DROP INDEX IF EXISTS idx_location_person_id;
DROP INDEX IF EXISTS idx_location_last_place_id;
DROP INDEX IF EXISTS idx_department_id;
DROP INDEX IF EXISTS idx_employee_id;
DROP INDEX IF EXISTS idx_location_instrument_dates;
DROP INDEX IF EXISTS idx_write_off_instrument_date;
DROP INDEX IF EXISTS idx_transfer_to_dep_instrument_date;
DROP INDEX IF EXISTS idx_transfer_to_save_instrument_date;
DROP INDEX IF EXISTS idx_preservation_instrument_date;
DROP INDEX IF EXISTS idx_repair_instrument_period;
DROP INDEX IF EXISTS idx_verification_docs_verification_id;
DROP INDEX IF EXISTS idx_verification_instrument_date;
DROP INDEX IF EXISTS idx_instruments_section_status;

-- +goose StatementEnd
