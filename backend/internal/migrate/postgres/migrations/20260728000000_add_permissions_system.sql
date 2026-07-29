-- +goose Up

-- 1. Таблица permissions
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    object TEXT NOT NULL,
    action TEXT NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(object, action)
);

-- 2. Таблица role_permissions
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

-- 3. Таблица role_hierarchy
CREATE TABLE IF NOT EXISTS role_hierarchy (
    parent_role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY(parent_role_id, role_id),
    CHECK(parent_role_id <> role_id)
);

-- 4. Триггер предотвращения циклов в role_hierarchy
CREATE OR REPLACE FUNCTION check_role_hierarchy_cycle() RETURNS TRIGGER AS $$
DECLARE
    cycle_exists BOOLEAN;
BEGIN
    WITH RECURSIVE ancestors AS (
        SELECT parent_role_id AS id
        FROM role_hierarchy
        WHERE role_id = NEW.role_id

        UNION ALL

        SELECT rh.parent_role_id
        FROM role_hierarchy rh
        JOIN ancestors a ON rh.role_id = a.id
    )
    SELECT EXISTS(SELECT 1 FROM ancestors WHERE id = NEW.parent_role_id) INTO cycle_exists;

    IF cycle_exists THEN
        RAISE EXCEPTION 'circular inheritance detected: role % cannot inherit from %', NEW.role_id, NEW.parent_role_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_role_hierarchy_cycle
    BEFORE INSERT OR UPDATE ON role_hierarchy
    FOR EACH ROW
    EXECUTE FUNCTION check_role_hierarchy_cycle();

-- 5. Data migration: permissions из rule_item
INSERT INTO permissions (id, object, action, name, description)
SELECT gen_random_uuid(), name, method, name, COALESCE(description, '')
FROM rule_item
ON CONFLICT (object, action) DO UPDATE
SET name = EXCLUDED.name, description = EXCLUDED.description;

-- 6. Data migration: role_permissions из rule + rule_item
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), r.role_id, p.id
FROM rule r
INNER JOIN rule_item ri ON ri.id = r.rule_item_id
INNER JOIN permissions p ON p.object = ri.name AND p.action = ri.method
ON CONFLICT DO NOTHING;

-- 7. Data migration: role_hierarchy из roles.extends
INSERT INTO role_hierarchy (parent_role_id, role_id)
SELECT unnest(r.extends), r.id
FROM roles r
WHERE r.extends IS NOT NULL AND array_length(r.extends, 1) > 0
ON CONFLICT DO NOTHING;

-- +goose Down
DROP TRIGGER IF EXISTS trigger_check_role_hierarchy_cycle ON role_hierarchy;
DROP FUNCTION IF EXISTS check_role_hierarchy_cycle();
DROP TABLE IF EXISTS role_hierarchy;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
