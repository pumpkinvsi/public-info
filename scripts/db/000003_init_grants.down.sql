-- =============================================================
-- 000003_init_grants.down.sql
-- Revokes privileges granted to application roles.
-- =============================================================

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    REVOKE ALL PRIVILEGES ON SEQUENCES FROM admin;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    REVOKE ALL PRIVILEGES ON TABLES FROM admin;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    REVOKE USAGE, SELECT ON SEQUENCES FROM bff;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    REVOKE SELECT, INSERT, UPDATE, DELETE ON TABLES FROM bff;

REVOKE ALL PRIVILEGES
    ON ALL SEQUENCES IN SCHEMA public
    FROM admin;

REVOKE ALL PRIVILEGES
    ON ALL TABLES IN SCHEMA public
    FROM admin;

REVOKE USAGE, SELECT
    ON ALL SEQUENCES IN SCHEMA public
    FROM bff;

REVOKE SELECT, INSERT, UPDATE, DELETE
    ON ALL TABLES IN SCHEMA public
    FROM bff;