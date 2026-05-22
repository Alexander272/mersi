-- ======================================================================
-- Миграция данных: SI accounting → mersi_v2
-- Запускать на БД mersi_v2 (целевая).
-- Использует dblink — не создаёт постоянных объектов FDW.
--
-- Идемпотентно: при повторном запуске существующие записи пропускаются.
-- Если нужно обновить существующие — замените DO NOTHING на DO UPDATE.
-- ======================================================================
-- Порядок действий:
--   1. Установить расширение: CREATE EXTENSION IF NOT EXISTS dblink;
--   2. Заменить параметры подключения к старой БД в dblink_connect
--   3. Выполнить скрипт
--   4. Выполнить dblink_disconnect (если не в той же сессии)
-- ======================================================================

-- Подключение к старой БД (замените параметры)
SELECT dblink_connect('old', 'host=localhost port=5432 dbname=si_accounting user=postgres password=password');

-- ======================================================================
-- 1. Departments
-- ======================================================================
INSERT INTO public.departments (id, name, leader_id, created_at, channel_id, realm_id)
SELECT id, name, leader_id, created_at, channel_id, realm_id
FROM dblink('old', 'SELECT id, name, leader_id, created_at, channel_id, realm_id FROM public.departments')
AS t (id uuid, name text, leader_id uuid, created_at timestamptz, channel_id uuid, realm_id uuid)
ON CONFLICT (id) DO NOTHING;

-- ======================================================================
-- 2. Employee
-- ======================================================================
INSERT INTO public.employee (id, name, department_id, most_id, is_lead, created_at, sso_id, channel_id, notes)
SELECT id, name, department_id, most_id, is_lead, created_at,
       COALESCE(NULLIF(sso_id, ''), '') AS sso_id,
       ''::text AS channel_id,
       COALESCE(notes, '') AS notes
FROM dblink('old', 'SELECT id, name, department_id, most_id, sso_id, is_lead, created_at, notes FROM public.employee')
AS t (id uuid, name text, department_id uuid, most_id text, sso_id text, is_lead boolean, created_at timestamptz, notes text)
ON CONFLICT (id) DO NOTHING;

-- ======================================================================
-- 3. Accesses
--     user_id в старой и новой БД разные (internal UUID).
--     Матчим по sso_id (Keycloak UUID): old.users → public.users.
--     Если пользователь не найден в новой БД — строка пропускается.
-- ======================================================================
BEGIN;
CREATE TEMP TABLE _old_users ON COMMIT DROP AS
SELECT id, sso_id
FROM dblink('old', 'SELECT id, sso_id FROM public.users')
AS t (id uuid, sso_id text);

INSERT INTO public.accesses (id, realm_id, user_id, sso_id, role_id, created_at)
SELECT a.id, a.realm_id, nu.id AS user_id, nu.sso_id, a.role_id, a.created_at
FROM dblink('old', 'SELECT id, realm_id, user_id, role_id, created_at FROM public.accesses')
AS a (id uuid, realm_id uuid, user_id uuid, role_id uuid, created_at timestamptz)
JOIN _old_users ou ON ou.id = a.user_id
JOIN public.users nu ON nu.sso_id = ou.sso_id
ON CONFLICT (id) DO NOTHING;
COMMIT;

-- ======================================================================
-- 4. Responsible
-- ======================================================================
INSERT INTO public.responsible (id, department_id, sso_id, created_at)
SELECT id, department_id, sso_id, created_at
FROM dblink('old', 'SELECT id, department_id, sso_id, created_at FROM public.responsible')
AS t (id uuid, department_id uuid, sso_id text, created_at timestamptz)
ON CONFLICT (id) DO NOTHING;

-- ======================================================================
-- 5. Instruments
--     section_id: подбирается по realm_id инструмента
--     year_of_issue / inter_verification_interval: text → integer
-- ======================================================================
INSERT INTO public.instruments (
    id, section_id, user_id, position, name, date_of_receipt,
    type, factory_number, measurement_limits, accuracy, state_register,
    country_of_produce, manufacturer, responsible, inventory,
    year_of_issue, inter_verification_interval,
    act_of_entering, act_of_entering_id, notes, status,
    created_at, updated_at, deleted
)
SELECT
    sub.id, sub.section_id,
    '328f6b3d-6211-4013-bedf-2272349b1f44'::uuid AS user_id,
    ROW_NUMBER() OVER (PARTITION BY sub.section_id ORDER BY sub.created_at) AS position,
    sub.name, sub.created_at AS date_of_receipt,
    sub.type, sub.factory_number, sub.measurement_limits, sub.accuracy,
    sub.state_register,
    ''::text AS country_of_produce,
    sub.manufacturer,
    ''::text AS responsible,
    ''::text AS inventory,
    CASE WHEN sub.year_of_issue ~ '^\d+$' THEN sub.year_of_issue::integer ELSE 0 END AS year_of_issue,
    CASE WHEN sub.inter_verification_interval ~ '^\d+$' THEN sub.inter_verification_interval::integer ELSE 0 END AS inter_verification_interval,
    ''::text AS act_of_entering,
    '00000000-0000-0000-0000-000000000000'::uuid AS act_of_entering_id,
    sub.notes, sub.status,
    sub.created_at, sub.created_at AS updated_at,
    NULL::timestamptz AS deleted
FROM (
    SELECT
        i.id,
        COALESCE(
            (SELECT s.id FROM public.sections s WHERE s.realm_id = i.realm_id LIMIT 1),
            (SELECT id FROM public.sections LIMIT 1)
        ) AS section_id,
        i.name, i.type, i.factory_number, i.measurement_limits, i.accuracy,
        i.state_register, i.manufacturer, i.year_of_issue,
        i.inter_verification_interval, i.notes, i.status, i.created_at
    FROM dblink('old', 'SELECT id, name, type, factory_number, measurement_limits, accuracy, state_register, manufacturer, year_of_issue, inter_verification_interval, notes, status, created_at, realm_id FROM public.instruments')
    AS i (id uuid, name text, type text, factory_number text, measurement_limits text, accuracy text, state_register text, manufacturer text, year_of_issue text, inter_verification_interval text, notes text, status text, created_at timestamptz, realm_id uuid)
) sub
ON CONFLICT (id) DO NOTHING;

-- ======================================================================
-- 6. Locations (бывший si_movement_history)
--     date_of_receiving / date_of_issue: integer (unix) → timestamp
-- ======================================================================
INSERT INTO public.locations (
    id, instrument_id, date_of_receiving, date_of_issue, status,
    has_confirmed, need_confirmed, person, person_id, place,
    department_id, last_place, last_place_id, user_id, created_at
)
SELECT
    m.id,
    m.instrument_id,
    CASE WHEN m.date_of_receiving > 0
        THEN to_timestamp(m.date_of_receiving)::timestamp
        ELSE '0001-01-01'::DATE
    END AS date_of_receiving,
    CASE WHEN m.date_of_issue > 0
        THEN to_timestamp(m.date_of_issue)::timestamp
        ELSE '0001-01-01'::DATE
    END AS date_of_issue,
    m.status,
    m.has_confirmed,
    m.need_confirmed,
    COALESCE(m.person, '') AS person,
    m.person_id,
    COALESCE(m.place, '') AS place,
    m.department_id,
    COALESCE(m.last_place, '') AS last_place,
    COALESCE(m.last_place_id, '') AS last_place_id,
    '328f6b3d-6211-4013-bedf-2272349b1f44'::uuid AS user_id,
    m.created_at
FROM dblink('old', 'SELECT id, instrument_id, date_of_receiving, date_of_issue, status, has_confirmed, need_confirmed, person, person_id, place, department_id, last_place, last_place_id, created_at FROM public.si_movement_history')
AS m (id uuid, instrument_id uuid, date_of_receiving integer, date_of_issue integer, status text, has_confirmed boolean, need_confirmed boolean, person text, person_id uuid, place text, department_id uuid, last_place text, last_place_id text, created_at timestamptz)
ON CONFLICT (id) DO NOTHING;

-- ======================================================================
-- 7. Verifications (бывший verification_history)
--     date / next_date: integer (unix) → timestamp
-- ======================================================================
INSERT INTO public.verifications (
    id, instrument_id, register_link, status, date, next_date,
    notes, not_verified, created_at, updated_at
)
SELECT
    v.id,
    v.instrument_id,
    v.register_link,
    v.status,
    CASE WHEN v.date > 0
        THEN to_timestamp(v.date)::timestamp
        ELSE '0001-01-01'::DATE
    END AS date,
    CASE WHEN v.next_date > 0
        THEN to_timestamp(v.next_date)::timestamp
        ELSE '0001-01-01'::DATE
    END AS next_date,
    v.notes,
    COALESCE(v.not_verified, false) AS not_verified,
    v.created_at,
    v.created_at AS updated_at
FROM dblink('old', 'SELECT id, instrument_id, register_link, status, date, next_date, notes, not_verified, created_at FROM public.verification_history')
AS v (id uuid, instrument_id uuid, register_link text, status text, date integer, next_date integer, notes text, not_verified boolean, created_at timestamptz)
ON CONFLICT (id) DO NOTHING;

-- ======================================================================
-- 8. Documents (verification_id уходит в verification_docs)
-- ======================================================================
INSERT INTO public.documents (
    id, label, size, path, type, instrument_id, belongs, user_id, created_at
)
SELECT
    d.id,
    d.label,
    d.size,
    d.path,
    d.type,
    d.instrument_id,
    ''::text AS belongs,
    '328f6b3d-6211-4013-bedf-2272349b1f44'::uuid AS user_id,
    d.created_at
FROM dblink('old', 'SELECT id, instrument_id, verification_id, label, size, path, type, created_at FROM public.documents')
AS d (id uuid, instrument_id uuid, verification_id uuid, label text, size integer, path text, type text, created_at timestamptz)
ON CONFLICT (id) DO NOTHING;

-- ======================================================================
-- 9. verification_docs
--     Связки для документов, у которых был verification_id
-- ======================================================================
INSERT INTO public.verification_docs (id, verification_id, doc_id, name, created_at, updated_at)
SELECT
    gen_random_uuid() AS id,
    d.verification_id,
    d.id AS doc_id,
    d.label AS name,
    d.created_at,
    d.created_at AS updated_at
FROM dblink('old', 'SELECT id, verification_id, label, created_at FROM public.documents WHERE verification_id IS NOT NULL')
AS d (id uuid, verification_id uuid, label text, created_at timestamptz)
WHERE NOT EXISTS (
    SELECT 1 FROM public.verification_docs x
    WHERE x.doc_id = d.id AND x.verification_id = d.verification_id
)
AND EXISTS (SELECT 1 FROM public.verifications v WHERE v.id = d.verification_id)
AND EXISTS (SELECT 1 FROM public.documents nd WHERE nd.id = d.id);

-- ======================================================================
-- 10. Filters (бывший default_filters)
--      filter_name → name; section_id через realm_id
-- ======================================================================
INSERT INTO public.filters (id, sso_id, name, field_type, compare_type, value, section_id, created_at, updated_at)
SELECT
    df.id,
    df.sso_id,
    df.filter_name AS name,
    df.filter_name AS field_type,
    df.compare_type,
    df.value,
    COALESCE(
        (SELECT s.id FROM public.sections s WHERE s.realm_id = df.realm_id LIMIT 1),
        (SELECT id FROM public.sections LIMIT 1)
    ) AS section_id,
    now() AS created_at,
    now() AS updated_at
FROM dblink('old', 'SELECT id, sso_id, filter_name, compare_type, value, realm_id FROM public.default_filters')
AS df (id uuid, sso_id uuid, filter_name text, compare_type text, value text, realm_id uuid)
ON CONFLICT (id) DO NOTHING;

-- ======================================================================
-- Отключение
-- ======================================================================
SELECT dblink_disconnect('old');
