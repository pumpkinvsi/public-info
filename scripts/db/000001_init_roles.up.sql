-- =============================================================
-- 01_roles.sql
-- Creates application roles.
-- Run as superuser (postgres).
-- Passwords must be set via ALTER ROLE after provisioning,
-- or injected by a secrets manager / compose env substitution.
-- =============================================================

DO $$
BEGIN
    -- admin: privileged role for manual DBA tasks
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'admin') THEN
        CREATE ROLE admin
            WITH LOGIN
                 CREATEDB
                 CREATEROLE
                 REPLICATION;
    END IF;

    -- bff: read-write role for the backend-for-frontend service
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'bff') THEN
        CREATE ROLE bff
            WITH LOGIN;
    END IF;
END
$$;