-- =============================================================
-- 03_grants.sql
-- Assigns privileges to application roles.
-- Depends on: 01_roles.sql, 02_tables.sql
-- =============================================================

-- bff: full DML on all tables in public schema
GRANT SELECT, INSERT, UPDATE, DELETE
    ON ALL TABLES IN SCHEMA public
    TO bff;

-- bff: sequence usage so SERIAL columns work on INSERT
GRANT USAGE, SELECT
    ON ALL SEQUENCES IN SCHEMA public
    TO bff;

-- admin: full ownership-level access
GRANT ALL PRIVILEGES
    ON ALL TABLES IN SCHEMA public
    TO admin;

GRANT ALL PRIVILEGES
    ON ALL SEQUENCES IN SCHEMA public
    TO admin;

-- Ensure future tables/sequences created in this schema
-- automatically carry the same privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO bff;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT USAGE, SELECT ON SEQUENCES TO bff;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON TABLES TO admin;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON SEQUENCES TO admin;