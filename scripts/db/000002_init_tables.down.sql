-- =============================================================
-- 000002_init_tables.down.sql
-- Drops all application tables.
-- Order is strict: junction tables first, then dependants,
-- then lookup tables to satisfy foreign key constraints.
-- =============================================================

DROP TABLE IF EXISTS project_groups;
DROP TABLE IF EXISTS project_skills;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS technologies;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS skills;
DROP TABLE IF EXISTS info;
DROP TABLE IF EXISTS levels;