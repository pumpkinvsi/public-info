-- =============================================================
-- 02_tables.sql
-- Creates all application tables.
-- Depends on: 01_roles.sql
-- =============================================================

-- ----------------------------
-- levels
-- Skill level lookup table.
-- `level` carries the numeric weight used by the frontend.
-- ----------------------------
CREATE TABLE IF NOT EXISTS levels (
    id    SERIAL      PRIMARY KEY,
    level INTEGER     NOT NULL,
    title VARCHAR(64) NOT NULL
);

-- ----------------------------
-- info
-- Singleton table: one row, no surrogate key.
-- Enforced by a CHECK constraint on row count via a partial unique index.
-- `bio` is a JSONB localized map, e.g. {"en": "...", "uk": "..."}.
-- ----------------------------
CREATE TABLE IF NOT EXISTS info (
    id  INTEGER     PRIMARY KEY DEFAULT 1,
    bio JSONB       NOT NULL,

    CONSTRAINT info_single_row CHECK (id = 1)
);

-- ----------------------------
-- skills
-- Individual skills.
-- `level` references the levels lookup.
-- ----------------------------
CREATE TABLE IF NOT EXISTS skills (
    id    SERIAL      PRIMARY KEY,
    name  VARCHAR(128) NOT NULL,
    level INTEGER      NOT NULL REFERENCES levels (id) ON DELETE RESTRICT
);

-- ----------------------------
-- projects
-- Personal / professional projects.
-- All text fields are JSONB localized maps.
-- `note` is optional (nullable).
-- ----------------------------
CREATE TABLE IF NOT EXISTS projects (
    id          SERIAL  PRIMARY KEY,
    name        JSONB   NOT NULL,
    description JSONB   NOT NULL,
    note        JSONB
);

-- ----------------------------
-- project_skills  (many-to-many: projects ↔ skills)
-- Composite primary key prevents duplicate pairs.
-- ----------------------------
CREATE TABLE IF NOT EXISTS project_skills (
    project_id INTEGER NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    skill_id   INTEGER NOT NULL REFERENCES skills   (id) ON DELETE CASCADE,

    PRIMARY KEY (project_id, skill_id)
);

-- ----------------------------
-- technologies
-- Technology / stack entries.
-- ----------------------------
CREATE TABLE IF NOT EXISTS technologies (
    id   SERIAL       PRIMARY KEY,
    name VARCHAR(128) NOT NULL
);

-- ----------------------------
-- project_groups  (many-to-many: projects ↔ technologies)
-- Groups projects under a technology umbrella.
-- Composite primary key prevents duplicate pairs.
-- ----------------------------
CREATE TABLE IF NOT EXISTS project_groups (
    technology_id INTEGER NOT NULL REFERENCES technologies (id) ON DELETE CASCADE,
    project_id    INTEGER NOT NULL REFERENCES projects     (id) ON DELETE CASCADE,

    PRIMARY KEY (technology_id, project_id)
);

-- ----------------------------
-- contacts
-- Contact information entries (e.g. email, LinkedIn).
-- ----------------------------
CREATE TABLE IF NOT EXISTS contacts (
    id    SERIAL       PRIMARY KEY,
    name  VARCHAR(128) NOT NULL,
    value VARCHAR(256) NOT NULL
);