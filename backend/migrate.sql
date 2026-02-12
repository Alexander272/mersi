-- Проверить наличие расширения
SELECT * FROM pg_available_extensions WHERE name = 'dblink';

-- Установить расширение (если не установлено)
CREATE EXTENSION IF NOT EXISTS dblink;

-- копирование инструментов
INSERT INTO public.instruments (
  id, section_id, user_id, name, date_of_receipt, 
  type, factory_number, measurement_limits, 
  accuracy, state_register, country_of_produce, 
  manufacturer, responsible, inventory, 
  year_of_issue, inter_verification_interval, 
  act_of_entering, act_of_entering_id, 
  notes, status, created_at, updated_at, 
  "position", deleted
) 
SELECT 
  id, 
  '4102cfbd-4a31-4d05-a89c-d20ab6653582'::uuid AS section_id, -- Укажите нужный UUID или значение по умолчанию
  '328f6b3d-6211-4013-bedf-2272349b1f44'::uuid AS user_id, -- Значение по умолчанию для user_id
  name, 
  created_at AS date_of_receipt,
  type, 
  factory_number, 
  measurement_limits, 
  accuracy, 
  state_register, 
  ''::text AS country_of_produce, -- Значение по умолчанию
  manufacturer, 
  ''::text AS responsible, -- Значение по умолчанию
  ''::text AS inventory, -- Значение по умолчанию
  CASE WHEN year_of_issue ~ '^[0-9]+$' THEN year_of_issue::integer ELSE 0 END AS year_of_issue, -- Преобразование текста в число с проверкой
  CASE WHEN inter_verification_interval ~ '^[0-9]+$' THEN inter_verification_interval::integer ELSE 0 END AS inter_verification_interval, 
  -- Преобразование текста в число с проверкой
  ''::text AS act_of_entering, -- Значение по умолчанию
  '00000000-0000-0000-0000-000000000000'::uuid AS act_of_entering_id, -- Значение по умолчанию
  notes, 
  status, 
  created_at, 
  NOW() AS updated_at, 
  ROW_NUMBER() OVER (ORDER BY created_at) AS "position",
  CASE WHEN status='deleted' THEN NOW() ELSE NULL END AS deleted -- Значение по умолчанию
FROM dblink(
    'host= port=5432 dbname= user= password=',
    'SELECT 
        id,
        name,
        type,
        factory_number,
        measurement_limits,
        accuracy,
        state_register,
        manufacturer,
        year_of_issue,
        inter_verification_interval,
        notes,
        status,
		realm_id,
        created_at
     FROM instruments'
) AS remote_data(
    id uuid,
    name text,
    type text,
    factory_number text,
    measurement_limits text,
    accuracy text,
    state_register text,
    manufacturer text,
    year_of_issue text,
    inter_verification_interval text,
    notes text,
    status text,
	realm_id uuid,
    created_at timestamp with time zone
)
WHERE realm_id='2236af5f-3151-4ac0-b10b-1975ffabde5a';


-- копирование поверок
INSERT INTO public.verifications(
	id, instrument_id, register_link, status, date, next_date, notes, not_verified, created_at, updated_at)
SELECT 
  id, instrument_id, register_link, status,
  to_timestamp(date)::date, 
  to_timestamp(next_date)::date,
  notes, not_verified, 
  created_at, created_at AS updated_at
FROM dblink(
    'host= port=5432 dbname= user= password=',
    'SELECT 
        id, instrument_id, register_link, status, date, next_date, notes, created_at, not_verified
     FROM verification_history AS v
     WHERE EXISTS (
         SELECT 1 FROM instruments AS i
         WHERE i.id = v.instrument_id 
         AND i.realm_id = ''2236af5f-3151-4ac0-b10b-1975ffabde5a''
         AND i.inter_verification_interval != ''''
     )'
) AS remote_data(
    id uuid,
    instrument_id uuid,
    register_link text,
    status text,
    date integer,
    next_date integer,
    notes text,
    created_at timestamp with time zone,
    not_verified boolean
);


-- копирование перемещений
INSERT INTO public.locations(
	id, instrument_id, date_of_receiving, date_of_issue, status, has_confirmed, need_confirmed, 
	person, person_id, place, department_id, last_place, last_place_id, user_id, created_at)
SELECT 
  id, instrument_id, 
  to_timestamp(date_of_receiving)::date, to_timestamp(date_of_issue)::date, 
  status, has_confirmed, need_confirmed, 
  person, person_id, place, department_id, last_place, last_place_id,
  user_id, created_at
FROM dblink(
    'host= port=5432 dbname= user= password=',
    'SELECT 
        id, instrument_id, date_of_receiving, date_of_issue, place, status, person, last_place, 
		created_at, person_id, department_id, has_confirmed, need_confirmed, last_place_id,
		''328f6b3d-6211-4013-bedf-2272349b1f44''::uuid AS user_id
     FROM si_movement_history AS m
     WHERE EXISTS (
         SELECT 1 FROM instruments AS i
         WHERE i.id = m.instrument_id 
         AND i.realm_id = ''2236af5f-3151-4ac0-b10b-1975ffabde5a''
     )'
) AS remote_data(
    id uuid,
    instrument_id uuid,
    date_of_receiving integer,
    date_of_issue integer,
    place text,
    status text,
    person text,
    last_place text,
    created_at timestamp with time zone,
    person_id uuid,
    department_id uuid,
    has_confirmed boolean,
    need_confirmed boolean,
    last_place_id text,
	user_id uuid
);