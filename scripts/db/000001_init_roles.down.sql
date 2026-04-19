-- =============================================================
-- 000001_init_roles.down.sql
-- Drops application roles.
-- Roles can only be dropped when they own no objects
-- and hold no granted privileges — run after 000003 down
-- and 000002 down to ensure that precondition is met.
-- =============================================================

DROP ROLE IF EXISTS bff;
DROP ROLE IF EXISTS admin;